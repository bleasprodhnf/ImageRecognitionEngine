package router

import (
	"github.com/gin-gonic/gin"
	"github.com/image-recognition-engine/internal/handler"
	"github.com/image-recognition-engine/internal/handler/admin"
	"github.com/image-recognition-engine/internal/handler/auth"
	"github.com/image-recognition-engine/internal/handler/client"
	"github.com/image-recognition-engine/internal/handler/common"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(app *gin.Engine) {
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

	// 管理员端路由
	adminRoutes := apiV1.Group("/admin")
	{
		// 系统管理
		systemRoutes := adminRoutes.Group("/system")
		{
			systemRoutes.GET("/params", admin.GetSystemParams)
			systemRoutes.POST("/params", admin.CreateSystemParam)
			systemRoutes.PUT("/params/:id", admin.UpdateSystemParam)
			systemRoutes.DELETE("/params/:id", admin.DeleteSystemParam)
		}

		// 用户权限管理
		userRoutes := adminRoutes.Group("/users")
		{
			userRoutes.GET("", admin.GetUsers)
			userRoutes.POST("", admin.CreateUser)
			userRoutes.PUT("/:id", admin.UpdateUser)
			userRoutes.DELETE("/:id", admin.DeleteUser)
		}

		// 角色管理
		roleRoutes := adminRoutes.Group("/roles")
		{
			roleRoutes.GET("", admin.GetRoles)
			roleRoutes.POST("", admin.CreateRole)
			roleRoutes.PUT("/:id", admin.UpdateRole)
			roleRoutes.DELETE("/:id", admin.DeleteRole)
		}

		// 模型管理
		modelRoutes := adminRoutes.Group("/models")
		{
			modelRoutes.GET("", admin.GetModels)
			modelRoutes.POST("", admin.CreateModel)
			modelRoutes.PUT("/:id", admin.UpdateModel)
			modelRoutes.DELETE("/:id", admin.DeleteModel)
		}

		// 客户管理
		customerRoutes := adminRoutes.Group("/customers")
		{
			customerRoutes.GET("", admin.GetCustomers)
			customerRoutes.POST("", admin.CreateCustomer)
			customerRoutes.PUT("/:id", admin.UpdateCustomer)
			customerRoutes.DELETE("/:id", admin.DeleteCustomer)
		}

		// 数据统计
		statsRoutes := adminRoutes.Group("/stats")
		{
			statsRoutes.GET("/system", admin.GetSystemStats)
			statsRoutes.GET("/recognition", admin.GetRecognitionStats)
			statsRoutes.GET("/customers", admin.GetCustomerStats)
		}

		// 注册日志路由
		logHandler := handler.NewLogHandler(nil) // 需要传入实际的日志仓库实例
		RegisterLogRoutes(adminRoutes, logHandler)
	}

	// 客户端路由
	clientRoutes := apiV1.Group("/client")
	{
		// 图像识别
		clientRoutes.POST("/recognize", client.RecognizeImage)
		clientRoutes.GET("/history", client.GetRecognitionHistory)

		// 账户管理
		clientRoutes.GET("/account", client.GetAccountInfo)
		clientRoutes.PUT("/account", client.UpdateAccountInfo)
		clientRoutes.PUT("/password", client.ChangePassword)

		// 使用统计
		clientRoutes.GET("/stats", client.GetClientStats)
	}

	// 注册统计路由
	statsHandler := handler.NewStatsHandler(nil) // 需要传入实际的统计仓库实例
	RegisterStatsRoutes(apiV1, statsHandler)
}