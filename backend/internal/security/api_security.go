package security

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

// APISecurityService 提供API安全相关功能
type APISecurityService struct {
	redisClient *redis.Client
	limiters    map[string]*rate.Limiter
	encryption  *EncryptionService
}

// APIKeyInfo API密钥信息
type APIKeyInfo struct {
	AppID       string    `json:"appId"`
	APIKey      string    `json:"apiKey"`
	OwnerID     int64     `json:"ownerId"`
	Permissions []string  `json:"permissions"`
	RateLimit   int       `json:"rateLimit"` // 每分钟请求数限制
	ExpireAt    time.Time `json:"expireAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NewAPISecurityService 创建一个新的API安全服务实例
func NewAPISecurityService(redisClient *redis.Client, encryption *EncryptionService) *APISecurityService {
	return &APISecurityService{
		redisClient: redisClient,
		limiters:    make(map[string]*rate.Limiter),
		encryption:  encryption,
	}
}

// ValidateAPIKey 验证API密钥
func (s *APISecurityService) ValidateAPIKey(appID, apiKey string) (*APIKeyInfo, error) {
	if appID == "" || apiKey == "" {
		return nil, errors.New("AppID和APIKey不能为空")
	}

	// 这里应该从数据库或缓存中获取API密钥信息
	// 为了演示，这里使用硬编码的方式
	// TODO: 实现从数据库或Redis缓存中获取API密钥信息
	if appID == "test_app" && apiKey == "test_key" {
		return &APIKeyInfo{
			AppID:       appID,
			APIKey:      apiKey,
			OwnerID:     1,
			Permissions: []string{"recognition:read", "recognition:write"},
			RateLimit:   100,
			ExpireAt:    time.Now().AddDate(1, 0, 0),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}

	return nil, errors.New("无效的API密钥")
}

// CheckRateLimit 检查请求频率限制
func (s *APISecurityService) CheckRateLimit(ctx context.Context, apiKey string, limit int) error {
	// 使用Redis记录和检查请求频率
	key := fmt.Sprintf("rate_limit:%s", apiKey)
	count, err := s.redisClient.Get(ctx, key).Int()
	if err == redis.Nil {
		// 第一次请求
		err = s.redisClient.Set(ctx, key, 1, time.Hour).Err()
		if err != nil {
			return fmt.Errorf("设置请求计数失败: %v", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("获取请求计数失败: %v", err)
	}

	if count >= limit {
		return fmt.Errorf("超出API请求限制")
	}

	// 增加请求计数
	_, err = s.redisClient.Incr(ctx, key).Result()
	return err
}

// GenerateAPIKey 生成新的API密钥
func (s *APISecurityService) GenerateAPIKey(appID string, salt string) string {
	// 使用应用ID、时间戳和盐值生成API密钥
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	data := appID + timestamp + salt

	// 使用SHA-256哈希
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// LogAPIAccess 记录API访问日志
func (s *APISecurityService) LogAPIAccess(appID, endpoint, method, clientIP string, statusCode int, responseTime time.Duration) error {
	// TODO: 将日志保存到数据库或日志系统
	return nil
}

// APISecurityMiddleware API安全中间件
func (s *APISecurityService) APISecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取API密钥
		appID := c.GetHeader("X-App-ID")
		apiKey := c.GetHeader("X-API-Key")

		// 验证API密钥
		apiInfo, err := s.ValidateAPIKey(appID, apiKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的API密钥",
				"data":    nil,
			})
			return
		}

		// 检查API密钥是否过期
		if apiInfo.ExpireAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "API密钥已过期",
				"data":    nil,
			})
			return
		}

		// 检查权限
		endpoint := c.Request.URL.Path
		method := c.Request.Method
		if !s.checkPermission(apiInfo.Permissions, endpoint, method) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "没有访问权限",
				"data":    nil,
			})
			return
		}

		// 检查频率限制
		ctx := c.Request.Context()
		err = s.CheckRateLimit(ctx, appID, apiInfo.RateLimit)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求频率超过限制",
				"data":    nil,
			})
			return
		}

		// 记录请求开始时间
		startTime := time.Now()

		// 将API信息存储到上下文中
		c.Set("appID", appID)
		c.Set("ownerID", apiInfo.OwnerID)
		c.Set("permissions", apiInfo.Permissions)

		// 处理请求
		c.Next()

		// 计算响应时间
		responseTime := time.Since(startTime)

		// 记录访问日志
		s.LogAPIAccess(
			appID,
			endpoint,
			method,
			c.ClientIP(),
			c.Writer.Status(),
			responseTime,
		)
	}
}

// contains 检查字符串切片中是否包含指定字符串
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// checkPermission 检查是否有权限访问指定的端点
func (s *APISecurityService) checkPermission(permissions []string, endpoint, method string) bool {
	// 简化的权限检查逻辑
	// 实际应用中应该有更复杂的权限管理系统
	if contains(permissions, "admin") {
		return true
	}

	// 检查端点特定权限
	if strings.HasPrefix(endpoint, "/api/v1/recognition") {
		if method == http.MethodGet && contains(permissions, "recognition:read") {
			return true
		}
		if (method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete) && 
		   contains(permissions, "recognition:write") {
			return true
		}
	}

	return false
}