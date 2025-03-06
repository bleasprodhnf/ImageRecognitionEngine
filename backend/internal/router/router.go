package router

import (
	"github.com/gin-gonic/gin"
	"github.com/image-recognition-engine/internal/handler/auth"
	"github.com/image-recognition-engine/internal/handler/client"
	"github.com/image-recognition-engine/internal/handler/common"
	"github.com/image-recognition-engine/internal/middleware"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(app *gin.Engine) {
	// 添加CORS中间件
	app.Use(middleware.CORSMiddleware())

	// API版本前缀
	apiV1 := app.Group("/api/v1")

	// 健康检查
	app.GET("/health", common.HealthCheck)

	// 认证相关路由
	authRoutes := apiV1.Group("/auth")
	{
		authRoutes.POST("/login", auth.Login)
		authRoutes.POST("/register", auth.Register)
	}

	// 管理员端路由 - 暂时注释掉未实现的路由
	// _ = apiV1.Group("/admin")

	// 客户端路由
	clientRoutes := apiV1.Group("/client")
	{
		// 图像识别
		clientRoutes.POST("/recognize", client.RecognizeImage)
		// 暂时注释掉未实现的路由
		// clientRoutes.GET("/history", client.GetRecognitionHistory)
	}

	// 注册统计路由 - 暂时注释掉
	// 暂时不实现
}