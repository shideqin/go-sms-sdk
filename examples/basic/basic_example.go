package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shideqin/go-sms-sdk/pkg/providers/aliyun"
	"github.com/shideqin/go-sms-sdk/pkg/sms"
)

func main() {
	// 基础使用示例
	basicUsage()
}

// basicUsage 基础使用示例
func basicUsage() {
	fmt.Println("=== 基础使用示例 ===")

	// 1. 创建阿里云短信服务实例
	aliyunSMS := aliyun.NewAliyunSMS(
		"your_access_id",
		"your_access_key",
	)

	// 2. 创建发送请求
	req := &sms.SendRequest{
		PhoneNumbers: "your_phone_number",
		SignName:     "your_sign_name",
		TemplateCode: "your_template_code",
		TemplateParam: map[string]interface{}{
			"code": "123456",
		},
	}

	// 3. 发送短信
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := aliyunSMS.Send(ctx, req)
	if err != nil {
		log.Fatalf("发送短信失败: %v", err)
	}

	// 4. 处理结果
	if resp.Success {
		fmt.Printf("短信发送成功! RequestId: %s, BizId: %s\n", resp.RequestId, resp.BizId)
	} else {
		fmt.Printf("短信发送失败! Code: %s, Message: %s\n", resp.Code, resp.Message)
	}
}
