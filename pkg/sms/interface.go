package sms

import (
	"context"
)

// SMS 短信服务接口
type SMS interface {
	// Send 发送短信
	Send(ctx context.Context, req *SendRequest) (*SendResponse, error)
	// GetName 获取服务商名称
	GetName() string
}

// SendRequest 发送短信请求
type SendRequest struct {
	PhoneNumbers  string         `json:"phone_numbers"`  // 手机号码
	SignName      string         `json:"sign_name"`      // 签名名称
	TemplateCode  string         `json:"template_code"`  // 模板代码
	TemplateParam map[string]any `json:"template_param"` // 模板参数
}

// SendResponse 发送短信响应
type SendResponse struct {
	Success   bool           `json:"success"`    // 是否成功
	Message   string         `json:"message"`    // 响应消息
	Code      string         `json:"code"`       // 响应码
	RequestId string         `json:"request_id"` // 请求ID
	BizId     string         `json:"biz_id"`     // 业务ID
	Data      map[string]any `json:"data"`       // 额外数据
}
