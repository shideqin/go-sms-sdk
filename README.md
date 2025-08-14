# 多渠道短信 SDK

一个支持多短信渠道的 Go SDK，提供统一的短信发送接口，支持阿里云等多个短信服务提供商。

## 特性

- 🚀 **多渠道支持**：支持多个短信服务提供商，可轻松扩展
- 🔄 **故障转移**：内置故障转移机制，提高发送成功率
- 🛡️ **并发安全**：线程安全设计，支持高并发场景
- 📊 **统一接口**：提供统一的API接口，降低使用成本
- 🎯 **易于扩展**：模块化设计，便于添加新的短信服务商
- 📝 **完整文档**：提供详细的API文档和使用示例
- 🧪 **测试覆盖**：包含单元测试和集成测试

## 目录结构

```
sms/
├── pkg/                   # 主要包代码
│   ├── sms/               # 核心短信接口
│   │   ├── interface.go   # 短信服务接口定义
│   │   └── manager.go     # 短信服务管理器
│   └── providers/         # 短信服务提供商实现
│       └── aliyun/        # 阿里云短信服务
│           └── aliyun.go  # 阿里云短信服务实现
├── internal/              # 内部包
│   ├── errors/           # 错误处理
│   │   └── errors.go     # 错误类型定义
│   └── utils/            # 工具函数
│       └── utils.go      # 通用工具函数
├── examples/             # 使用示例
│   └── basic/            # 基础示例
│       └── basic_example.go
├── go.mod                # Go模块文件
├── go.sum                # 依赖校验文件
└── README.md             # 项目说明文档
```

## 快速开始

### 安装

```bash
go get github.com/shideqin/go-sms-sdk
```

### 基础使用

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/shideqin/go-sms-sdk/pkg/providers/aliyun"
    "github.com/shideqin/go-sms-sdk/pkg/sms"
)

func main() {
    // 创建阿里云短信服务实例
    aliyunSMS := aliyun.NewAliyunSMS(
        "your_access_key_id",
        "your_access_key_secret",
    )

    // 创建发送请求
    req := &sms.SendRequest{
        PhoneNumbers: "your_phone_number",
        SignName:     "your_sign_name",
        TemplateCode: "your_template_code",
        TemplateParam: map[string]interface{}{
            "code": "123456",
        },
    }

    // 发送短信
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    resp, err := aliyunSMS.Send(ctx, req)
    if err != nil {
        fmt.Printf("发送短信失败: %v\n", err)
        return
    }

    if resp.Success {
        fmt.Printf("短信发送成功! RequestId: %s, BizId: %s\n", resp.RequestId, resp.BizId)
    } else {
        fmt.Printf("短信发送失败! Code: %s, Message: %s\n", resp.Code, resp.Message)
    }
}
```
## 🤝 贡献

欢迎提交 Issue 和 Pull Request！
