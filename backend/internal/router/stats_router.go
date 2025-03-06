package router

import (
	"github.com/gin-gonic/gin"
	"github.com/image-recognition-engine/internal/handler"
	"github.com/image-recognition-engine/internal/middleware"
)

// RegisterStatsRoutes 注册统计相关的路由
func RegisterStatsRoutes(r *gin.RouterGroup, statsHandler *handler.StatsHandler) {
	// 统计相关路由组
	stats := r.Group("/stats")
	stats.Use(middleware.RequireAuth()) // 需要认证

	// 系统使用统计
	stats.GET("/system", middleware.RequirePermission("stats:view"), statsHandler.GetSystemStats)

	// 识别准确率统计
	stats.GET("/accuracy", middleware.RequirePermission("stats:view"), statsHandler.GetAccuracyStats)

	// 资源使用统计
	stats.GET("/resources", middleware.RequirePermission("stats:view"), statsHandler.GetResourceStats)
}