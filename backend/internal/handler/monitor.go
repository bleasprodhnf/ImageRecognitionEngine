package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/image-recognition-engine/internal/model"
	"github.com/image-recognition-engine/internal/repository"
)

// MonitorHandler 监控处理器
type MonitorHandler struct {
	repo *repository.MonitorRepository
}

// NewMonitorHandler 创建监控处理器实例
func NewMonitorHandler(repo *repository.MonitorRepository) *MonitorHandler {
	return &MonitorHandler{repo: repo}
}

// GetServerMetrics 获取服务器监控指标
func (h *MonitorHandler) GetServerMetrics(c *gin.Context) {
	// 解析时间范围参数
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	// 转换时间字符串为time.Time
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的开始时间格式"})
		return
	}

	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的结束时间格式"})
		return
	}

	// 获取监控数据
	metrics, err := h.repo.GetServerMetrics(c.Request.Context(), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取服务器监控数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data":    metrics,
	})
}

// GetAPIMetrics 获取API调用指标
func (h *MonitorHandler) GetAPIMetrics(c *gin.Context) {
	// 解析时间范围参数
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	// 转换时间字符串为time.Time
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的开始时间格式"})
		return
	}

	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的结束时间格式"})
		return
	}

	// 获取监控数据
	metrics, err := h.repo.GetAPIMetrics(c.Request.Context(), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取API监控数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data":    metrics,
	})
}

// GetStorageMetrics 获取存储空间使用指标
func (h *MonitorHandler) GetStorageMetrics(c *gin.Context) {
	// 解析时间范围参数
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	// 转换时间字符串为time.Time
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的开始时间格式"})
		return
	}

	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的结束时间格式"})
		return
	}

	// 获取监控数据
	metrics, err := h.repo.GetStorageMetrics(c.Request.Context(), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取存储空间监控数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data":    metrics,
	})
}