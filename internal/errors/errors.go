package errors

import "fmt"

// SMS 短信相关错误类型
type SMS struct {
	Code    string
	Message string
	Err     error
}

func (e *SMS) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("SMS Error [Code: %s, Message: %s]: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("SMS Error [Code: %s, Message: %s]", e.Code, e.Message)
}

func (e *SMS) Unwrap() error {
	return e.Err
}

// NewSMS 创建短信错误
func NewSMS(code, message string) *SMS {
	return &SMS{
		Code:    code,
		Message: message,
	}
}

// NewSMSWithErr 创建带内部错误的短信错误
func NewSMSWithErr(code, message string, err error) *SMS {
	return &SMS{
		Code:    code,
		Message: message,
		Err:     err,
	}
}