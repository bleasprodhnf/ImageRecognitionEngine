package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"

	"github.com/image-recognition-engine/internal/middleware/auth"
)

// RegisterMiddlewares 注册所有中间件
func RegisterMiddlewares(app *gin.Engine) {
	// 跨域中间件
	app.Use(CORSMiddleware())
	// 响应中间件
	app.Use(ResponseMiddleware())
	// 请求日志中间件
	app.Use(LoggerMiddleware())
	// 恢复中间件
	app.Use(gin.Recovery())
	// API认证中间件
	app.Use(auth.AuthMiddleware())
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-App-ID, X-API-Key")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 日志记录
		gin.DefaultWriter.Write([]byte(fmt.Sprintf("%s | %d | %s | %s | %v\n",
			c.Request.Method,
			c.Writer.Status(),
			c.Request.URL.Path,
			c.ClientIP(),
			latency,
		)))
	}
}

// skipAuth 判断是否跳过认证
func skipAuth(path string) bool {
	// 静态资源路径跳过认证
	if strings.HasPrefix(path, "/static") {
		return true
	}

	// 健康检查和监控端点跳过认证
	if path == "/health" || path == "/metrics" {
		return true
	}

	// 登录和注册接口跳过认证
	if path == "/api/v1/auth/login" || path == "/api/v1/auth/register" {
		return true
	}

	return false
}