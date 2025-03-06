package router

import (
	"github.com/gin-gonic/gin"
	"github.com/image-recognition-engine/internal/handler"
	"github.com/image-recognition-engine/internal/middleware"
)

// RegisterLogRoutes 注册日志相关路由
func RegisterLogRoutes(r *gin.RouterGroup, logHandler *handler.LogHandler) {
	// 日志管理路由组
	logGroup := r.Group("/logs")
	logGroup.Use(middleware.RequireAdmin()) // 需要管理员权限

	// 系统日志
	logGroup.GET("/system", logHandler.GetSystemLogs)

	// API使用日志
	logGroup.GET("/api-usage", logHandler.GetAPIUsageLogs)
}