package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

// GenerateTimestamp 生成时间戳
func GenerateTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

// GenerateSignatureNonce 生成签名随机数
func GenerateSignatureNonce() string {
	// 生成16字节的随机数
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// 如果随机数生成失败，使用时间戳作为备选方案
		return fmt.Sprintf("%d%d", time.Now().UnixNano(), time.Now().Unix())
	}
	// 返回hex编码的随机字符串
	return hex.EncodeToString(b)
}

// PercentEncode 百分比编码
func PercentEncode(str string) string {
	// URL编码
	encoded := url.QueryEscape(str)
	// 阿里云API的特殊要求
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "*", "%2A")
	encoded = strings.ReplaceAll(encoded, "%7E", "~")
	return encoded
}

// BuildCanonicalQueryString 构建规范化查询字符串
func BuildCanonicalQueryString(params map[string]string) string {
	// 获取排序后的键
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建查询字符串
	var queryStrings []string
	for _, k := range keys {
		v := params[k]
		queryStrings = append(queryStrings, fmt.Sprintf("%s=%s", k, PercentEncode(v)))
	}

	return strings.Join(queryStrings, "&")
}

// HMACSHA256 HMAC-SHA256加密
func HMACSHA256(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// HMACSHA1 HMAC-SHA1加密
func HMACSHA1(data, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256Hash SHA256哈希
func SHA256Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// ValidatePhoneNumber 验证手机号码
func ValidatePhoneNumber(phone string) bool {
	// 简单的手机号码验证
	if len(phone) != 11 {
		return false
	}
	
	// 检查是否以1开头
	if phone[0] != '1' {
		return false
	}
	
	// 检查是否都是数字
	for _, c := range phone {
		if c < '0' || c > '9' {
			return false
		}
	}
	
	return true
}