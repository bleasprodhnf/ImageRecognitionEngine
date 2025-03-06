package security

import (
	"github.com/image-recognition-engine/config"
)

// SecurityConfig 安全配置结构
type SecurityConfig struct {
	EncryptionKey    string   `json:"encryptionKey"`    // 数据加密密钥
	JWTSecret        string   `json:"jwtSecret"`        // JWT密钥
	APIKeyExpiration int      `json:"apiKeyExpiration"` // API密钥过期时间(天)
	SensitiveFields  []string `json:"sensitiveFields"`  // 敏感字段列表
	RateLimitDefault int      `json:"rateLimitDefault"` // 默认API请求限制(每分钟)
	EnableAuditLog   bool     `json:"enableAuditLog"`   // 是否启用审计日志
}

// LoadSecurityConfig 从配置文件加载安全配置
func LoadSecurityConfig() (*SecurityConfig, error) {
	// 加载主配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// 这里应该从配置文件中读取安全配置
	// 为了演示，这里使用硬编码的方式
	return &SecurityConfig{
		EncryptionKey:    cfg.JWT.Secret, // 使用JWT密钥作为加密密钥(实际应用中应该分开)
		JWTSecret:        cfg.JWT.Secret,
		APIKeyExpiration: 365, // 默认一年
		SensitiveFields:  []string{"password", "phone", "email", "idCard", "bankCard"},
		RateLimitDefault: 100, // 每分钟100次请求
		EnableAuditLog:   true,
	}, nil
}