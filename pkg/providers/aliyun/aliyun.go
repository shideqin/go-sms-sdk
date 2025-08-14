package aliyun

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"go-sms-sdk/internal/utils"
	"go-sms-sdk/pkg/sms"
)

// AliyunSMS 阿里云短信服务
type AliyunSMS struct {
	accessKeyId     string
	accessKeySecret string
	endpoint        string
	version         string
	action          string
}

// NewAliyunSMS 创建阿里云短信服务实例
func NewAliyunSMS(accessKeyId, accessKeySecret string) *AliyunSMS {
	return &AliyunSMS{
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
		endpoint:        "dysmsapi.aliyuncs.com",
		version:         "2017-05-25",
		action:          "SendSms",
	}
}

// GetName 获取服务商名称
func (a *AliyunSMS) GetName() string {
	return "aliyun"
}

// Send 发送短信
func (a *AliyunSMS) Send(ctx context.Context, req *sms.SendRequest) (*sms.SendResponse, error) {
	// 1. 准备模板参数
	templateParam := a.buildTemplateParam(req.TemplateParam)

	// 2. 构建 body
	form := map[string]string{
		"PhoneNumbers":  req.PhoneNumbers,
		"SignName":      req.SignName,
		"TemplateCode":  req.TemplateCode,
		"TemplateParam": templateParam,
	}
	body := a.encodeForm(form)

	// 3. 生成时间戳和随机 nonce
	date := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())

	// 4. 计算 body SHA256
	bodyHash := utils.SHA256Hash(body)

	// 5. 构建规范化请求头（按字典序）
	headers := map[string]string{
		"host":                  a.endpoint,
		"x-acs-action":          a.action,
		"x-acs-content-sha256":  bodyHash,
		"x-acs-date":            date,
		"x-acs-signature-nonce": nonce,
		"x-acs-version":         a.version,
	}
	canonicalHeaders, signedHeaders := a.buildCanonicalHeaders(headers)

	// 6. 构建规范化请求串
	canonicalRequest := fmt.Sprintf("POST\n/\n\n%s\n%s\n%s", canonicalHeaders, signedHeaders, bodyHash)

	// 7. 构建待签名字符串
	stringToSign := "ACS3-HMAC-SHA256\n" + utils.SHA256Hash(canonicalRequest)

	// 8. 计算 HMAC-SHA256 签名
	signature := a.hmacSha256Hex(a.accessKeySecret, stringToSign)

	// 9. 构建 Authorization 头
	authorization := fmt.Sprintf(
		"ACS3-HMAC-SHA256 Credential=%s,SignedHeaders=%s,Signature=%s",
		a.accessKeyId,
		signedHeaders,
		signature,
	)

	// 10. 发送 HTTP POST 请求
	return a.sendRequestWithNewAPI(ctx, "https://"+a.endpoint+"/", body, authorization, headers)
}

// buildTemplateParam 构建模板参数
func (a *AliyunSMS) buildTemplateParam(params map[string]any) string {
	if params == nil {
		return ""
	}

	data, err := json.Marshal(params)
	if err != nil {
		return ""
	}

	return string(data)
}

// encodeForm 编码 form 为 x-www-form-urlencoded
func (a *AliyunSMS) encodeForm(m map[string]string) string {
	pairs := []string{}
	for k, v := range m {
		pairs = append(pairs, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
	}
	return strings.Join(pairs, "&")
}

// hmacSha256Hex 计算 HMAC-SHA256 哈希的十六进制表示
func (a *AliyunSMS) hmacSha256Hex(secret, msg string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))
}

// buildCanonicalHeaders 构建规范化请求头和已签名的请求头列表
func (a *AliyunSMS) buildCanonicalHeaders(headers map[string]string) (string, string) {
	// 1. 获取所有 header 名并排序
	keys := make([]string, 0, len(headers))
	for k := range headers {
		keys = append(keys, strings.ToLower(k))
	}
	sort.Strings(keys)

	// 2. 构建规范化请求头
	lines := make([]string, 0, len(headers))
	for _, k := range keys {
		// 规范化格式：key:value，去除值两端的空格
		lines = append(lines, fmt.Sprintf("%s:%s", k, strings.TrimSpace(headers[k])))
	}

	// 3. 返回规范化请求头（用换行符连接，最后要有换行符）和已签名的请求头列表
	return strings.Join(lines, "\n") + "\n", strings.Join(keys, ";")
}

// sendRequestWithNewAPI 使用新API发送HTTP请求
func (a *AliyunSMS) sendRequestWithNewAPI(ctx context.Context, url, body, authorization string, headers map[string]string) (*sms.SendResponse, error) {
	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", authorization)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	return a.parseResponse(respBody)
}

// parseResponse 解析响应
func (a *AliyunSMS) parseResponse(body []byte) (*sms.SendResponse, error) {
	var response struct {
		Code      string `json:"Code"`
		Message   string `json:"Message"`
		RequestId string `json:"RequestId"`
		BizId     string `json:"BizId"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &sms.SendResponse{
		Success:   response.Code == "OK",
		Message:   response.Message,
		Code:      response.Code,
		RequestId: response.RequestId,
		BizId:     response.BizId,
		Data:      make(map[string]any),
	}, nil
}
