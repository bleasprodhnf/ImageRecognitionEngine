package security

import (
	"os"
	"strconv"
	"strings"
)

// OverrideSecurityConfigFromEnv 从环境变量覆盖安全配置
func OverrideSecurityConfigFromEnv(cfg *SecurityConfig) {
	// 加密密钥
	if encKey := os.Getenv("SECURITY_ENCRYPTION_KEY"); encKey != "" {
		cfg.EncryptionKey = encKey
	}

	// JWT密钥
	if jwtSecret := os.Getenv("SECURITY_JWT_SECRET"); jwtSecret != "" {
		cfg.JWTSecret = jwtSecret
	}

	// API密钥过期时间
	if apiKeyExp := os.Getenv("SECURITY_API_KEY_EXPIRATION"); apiKeyExp != "" {
		if exp, err := strconv.Atoi(apiKeyExp); err == nil {
			cfg.APIKeyExpiration = exp
		}
	}

	// 敏感字段列表
	if sensitiveFields := os.Getenv("SECURITY_SENSITIVE_FIELDS"); sensitiveFields != "" {
		cfg.SensitiveFields = strings.Split(sensitiveFields, ",")
	}

	// 默认API请求限制
	if rateLimit := os.Getenv("SECURITY_RATE_LIMIT_DEFAULT"); rateLimit != "" {
		if limit, err := strconv.Atoi(rateLimit); err == nil {
			cfg.RateLimitDefault = limit
		}
	}

	// 是否启用审计日志
	if enableAuditLog := os.Getenv("SECURITY_ENABLE_AUDIT_LOG"); enableAuditLog != "" {
		if enableAuditLog == "true" || enableAuditLog == "1" {
			cfg.EnableAuditLog = true
		} else if enableAuditLog == "false" || enableAuditLog == "0" {
			cfg.EnableAuditLog = false
		}
	}
}