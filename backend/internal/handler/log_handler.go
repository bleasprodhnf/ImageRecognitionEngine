package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/image-recognition-engine/internal/model"
	"github.com/image-recognition-engine/internal/repository"
)

// LogHandler 日志处理器
type LogHandler struct {
	logRepo repository.LogRepository
}

// NewLogHandler 创建日志处理器实例
func NewLogHandler(logRepo repository.LogRepository) *LogHandler {
	return &LogHandler{
		logRepo: logRepo,
	}
}

// GetSystemLogs 获取系统日志
func (h *LogHandler) GetSystemLogs(c *gin.Context) {
	var params model.LogQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效"})
		return
	}

	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}

	// 解析时间范围
	if startStr := c.Query("startTime"); startStr != "" {
		start, err := time.Parse(time.RFC3339, startStr)
		if err == nil {
			params.StartTime = start
		}
	}
	if endStr := c.Query("endTime"); endStr != "" {
		end, err := time.Parse(time.RFC3339, endStr)
		if err == nil {
			params.EndTime = end
		}
	}

	// 查询日志
	response, err := h.logRepo.QuerySystemLogs(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询日志失败"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAPIUsageLogs 获取API使用日志
func (h *LogHandler) GetAPIUsageLogs(c *gin.Context) {
	var params model.LogQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效"})
		return
	}

	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}

	// 解析时间范围
	if startStr := c.Query("startTime"); startStr != "" {
		start, err := time.Parse(time.RFC3339, startStr)
		if err == nil {
			params.StartTime = start
		}
	}
	if endStr := c.Query("endTime"); endStr != "" {
		end, err := time.Parse(time.RFC3339, endStr)
		if err == nil {
			params.EndTime = end
		}
	}

	// 查询日志
	response, err := h.logRepo.QueryAPIUsageLogs(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询日志失败"})
		return
	}

	c.JSON(http.StatusOK, response)
}