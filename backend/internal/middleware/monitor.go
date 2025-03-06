package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/image-recognition-engine/internal/model"
	"github.com/image-recognition-engine/internal/repository"
)

// APIMonitor 创建API监控中间件
func APIMonitor(repo *repository.MonitorRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 计算响应时间
		duration := time.Since(start).Seconds()

		// 创建API指标记录
		metrics := &model.APIMetrics{
			Timestamp:    time.Now(),
			Endpoint:     c.Request.URL.Path,
			Method:       c.Request.Method,
			ResponseTime: duration,
			StatusCode:   c.Writer.Status(),
			ClientIP:     c.ClientIP(),
		}

		// 异步保存API指标
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := repo.SaveAPIMetrics(ctx, metrics); err != nil {
				// 记录错误但不影响请求处理
				// TODO: 添加日志记录
			}
		}()
	}
}