package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/image-recognition-engine/config"
	"github.com/image-recognition-engine/internal/handler/auth"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要认证的路径
		if isSkippedPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 根据路径前缀判断使用哪种认证方式
		if strings.HasPrefix(c.Request.URL.Path, "/api/v1/admin") {
			// 管理员端使用JWT认证
			handleJWTAuth(c)
		} else if strings.HasPrefix(c.Request.URL.Path, "/api/v1/client") {
			// 客户端使用API密钥认证
			handleAPIKeyAuth(c)
		} else {
			// 默认放行
			c.Next()
		}
	}
}

// isSkippedPath 判断是否为跳过认证的路径
func isSkippedPath(path string) bool {
	// 不需要认证的路径列表
	skippedPaths := []string{
		"/health",
		"/api/v1/auth/login",
		"/api/v1/auth/register",
	}

	for _, p := range skippedPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}

	return false
}

// handleJWTAuth 处理JWT认证
func handleJWTAuth(c *gin.Context) {
	// 从请求头获取令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未提供认证令牌",
			"data":    nil,
		})
		return
	}

	// 提取令牌
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "认证令牌格式错误",
			"data":    nil,
		})
		return
	}

	tokenString := parts[1]

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务器内部错误",
			"data":    nil,
		})
		return
	}

	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "无效的认证令牌",
			"data":    nil,
		})
		return
	}

	// 验证令牌
	if claims, ok := token.Claims.(*auth.JWTClaims); ok && token.Valid {
		// 将用户信息存储到上下文
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("realName", claims.RealName)
		c.Set("email", claims.Email)
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "无效的认证令牌",
			"data":    nil,
		})
		return
	}
}

// handleAPIKeyAuth 处理API密钥认证
func handleAPIKeyAuth(c *gin.Context) {
	// 从请求头获取AppID和API密钥
	appID := c.GetHeader("X-App-ID")
	apiKey := c.GetHeader("X-API-Key")

	if appID == "" || apiKey == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未提供API认证信息",
			"data":    nil,
		})
		return
	}

	// TODO: 从数据库验证AppID和API密钥
	// 这里简化处理，实际应用中需要查询数据库验证
	if appID != "test-app" || apiKey != "test-key" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "无效的API认证信息",
			"data":    nil,
		})
		return
	}

	// 将客户信息存储到上下文
	c.Set("customerId", 1) // 示例ID，实际应从数据库获取
	c.Set("appId", appID)
	c.Next()
}