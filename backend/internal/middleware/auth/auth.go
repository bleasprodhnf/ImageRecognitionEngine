package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要认证的路径
		if isSkippedPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 检查认证方式
		if strings.HasPrefix(c.Request.URL.Path, "/api/admin") {
			// 管理员JWT认证
			HandleJWTAuth(c)
		} else {
			// 客户端API密钥认证
			HandleAPIKeyAuth(c)
		}
	}
}

// isSkippedPath 判断是否跳过认证
func isSkippedPath(path string) bool {
	// 静态资源路径跳过认证
	if strings.HasPrefix(path, "/static") {
		return true
	}

	// 健康检查和监控端点跳过认证
	if path == "/health" || path == "/metrics" {
		return true
	}

	// 登录和注册接口跳过认证
	if path == "/api/v1/login" || path == "/api/v1/register" {
		return true
	}

	return false
}

// HandleJWTAuth 处理JWT认证
func HandleJWTAuth(c *gin.Context) {
	// 获取Authorization头
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未提供认证信息",
		})
		return
	}

	// 验证JWT token
	// TODO: 实现JWT token验证逻辑

	c.Next()
}

// HandleAPIKeyAuth 处理API密钥认证
func HandleAPIKeyAuth(c *gin.Context) {
	// 获取API认证信息
	appID := c.GetHeader("X-App-ID")
	apiKey := c.GetHeader("X-API-Key")

	if appID == "" || apiKey == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未提供API认证信息",
		})
		return
	}

	// 验证API密钥
	// TODO: 实现API密钥验证逻辑

	c.Next()
}
