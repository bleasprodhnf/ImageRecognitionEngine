package common

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Services  struct {
		Database bool `json:"database"`
		Redis    bool `json:"redis"`
		Storage  bool `json:"storage"`
	} `json:"services"`
}

// HealthCheck 处理健康检查请求
func HealthCheck(c *gin.Context) {
	// 构建响应
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}

	// TODO: 实际检查各服务状态
	response.Services.Database = true
	response.Services.Redis = true
	response.Services.Storage = true

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "服务正常",
		"data":    response,
	})
}