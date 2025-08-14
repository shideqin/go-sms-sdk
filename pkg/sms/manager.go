package sms

import (
	"context"
	"fmt"
	"sync"
)

// SMSManager 短信服务管理器
type SMSManager struct {
	providers  map[string]SMS
	defaultSMS SMS
	mu         sync.RWMutex
}

// NewSMSManager 创建短信服务管理器
func NewSMSManager() *SMSManager {
	return &SMSManager{
		providers: make(map[string]SMS),
	}
}

// Register 注册短信服务商
func (m *SMSManager) Register(name string, sms SMS) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.providers[name] = sms
}

// SetDefault 设置默认短信服务商
func (m *SMSManager) SetDefault(name string) error {
	m.mu.RLock()
	sms, exists := m.providers[name]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("provider %s not found", name)
	}

	m.mu.Lock()
	m.defaultSMS = sms
	m.mu.Unlock()

	return nil
}

// Send 发送短信
func (m *SMSManager) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	m.mu.RLock()
	if m.defaultSMS == nil {
		m.mu.RUnlock()
		return nil, fmt.Errorf("no default SMS provider configured")
	}
	sms := m.defaultSMS
	m.mu.RUnlock()

	return sms.Send(ctx, req)
}

// SendWithProvider 使用指定服务商发送短信
func (m *SMSManager) SendWithProvider(ctx context.Context, provider string, req *SendRequest) (*SendResponse, error) {
	m.mu.RLock()
	sms, exists := m.providers[provider]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("provider %s not found", provider)
	}

	return sms.Send(ctx, req)
}

// GetProvider 获取短信服务商
func (m *SMSManager) GetProvider(name string) (SMS, error) {
	m.mu.RLock()
	sms, exists := m.providers[name]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	return sms, nil
}

// GetProviders 获取所有已注册的服务商
func (m *SMSManager) GetProviders() map[string]SMS {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 创建副本
	providers := make(map[string]SMS)
	for k, v := range m.providers {
		providers[k] = v
	}

	return providers
}

// RemoveProvider 移除短信服务商
func (m *SMSManager) RemoveProvider(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.providers[name]; !exists {
		return fmt.Errorf("provider %s not found", name)
	}

	// 如果要移除的是默认服务商，清空默认设置
	if m.providers[name] == m.defaultSMS {
		m.defaultSMS = nil
	}

	delete(m.providers, name)
	return nil
}
