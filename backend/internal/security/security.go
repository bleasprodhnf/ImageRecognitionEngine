package security

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

// SecurityService 提供统一的安全服务接口
type SecurityService struct {
	Encryption  *EncryptionService
	APISecurity *APISecurityService
	DataMasking *DataMaskingService
	AuditLog    *AuditLogService
}

// NewSecurityService 创建一个新的安全服务实例
func NewSecurityService(mongoDB *mongo.Database, redisClient *redis.Client, encryptionKey string) (*SecurityService, error) {
	// 创建加密服务
	encryption, err := NewEncryptionService(encryptionKey)
	if err != nil {
		return nil, err
	}

	// 创建数据脱敏服务
	dataMasking := NewDataMaskingService(encryption)

	// 创建API安全服务
	apiSecurity := NewAPISecurityService(redisClient, encryption)

	// 创建审计日志服务
	auditLog := NewAuditLogService(mongoDB, redisClient, dataMasking, true)

	return &SecurityService{
		Encryption:  encryption,
		APISecurity: apiSecurity,
		DataMasking: dataMasking,
		AuditLog:    auditLog,
	}, nil
}

// RegisterMiddlewares 注册安全中间件
func (s *SecurityService) RegisterMiddlewares(r *gin.Engine) {
	// 注册API安全中间件
	r.Use(s.APISecurity.APISecurityMiddleware())

	// 可以在这里注册其他安全相关的中间件
}

// EncryptData 加密数据
func (s *SecurityService) EncryptData(plaintext string) (string, error) {
	return s.Encryption.Encrypt(plaintext)
}

// DecryptData 解密数据
func (s *SecurityService) DecryptData(ciphertext string) (string, error) {
	return s.Encryption.Decrypt(ciphertext)
}

// MaskSensitiveData 对敏感数据进行脱敏处理
func (s *SecurityService) MaskSensitiveData(data map[string]string) map[string]string {
	return s.DataMasking.MaskSensitiveData(data)
}

// LogSecurityEvent 记录安全事件
func (s *SecurityService) LogSecurityEvent(userID int64, username, action, resource, resourceID, clientIP, userAgent string, details map[string]interface{}) error {
	return s.AuditLog.LogUserAction(nil, userID, username, action, resource, resourceID, clientIP, userAgent, "success", details)
}

// ValidateAPIAccess 验证API访问权限
func (s *SecurityService) ValidateAPIAccess(appID, apiKey string) (*APIKeyInfo, error) {
	return s.APISecurity.ValidateAPIKey(appID, apiKey)
}

// GenerateNewAPIKey 生成新的API密钥
func (s *SecurityService) GenerateNewAPIKey(appID, salt string) string {
	return s.APISecurity.GenerateAPIKey(appID, salt)
}