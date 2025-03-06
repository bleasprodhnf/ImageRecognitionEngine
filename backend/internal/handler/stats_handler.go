package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/image-recognition-engine/internal/model"
	"github.com/image-recognition-engine/internal/repository"
)

type StatsHandler struct {
	statsRepo *repository.StatsRepository
}

func NewStatsHandler(statsRepo *repository.StatsRepository) *StatsHandler {
	return &StatsHandler{statsRepo: statsRepo}
}

// GetSystemStats 获取系统使用统计数据
func (h *StatsHandler) GetSystemStats(c *gin.Context) {
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	_ = c.Query("interval") // 保留获取参数，但标记为已使用

	// 解析时间参数
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "开始时间格式错误"})
		return
	}

	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "结束时间格式错误"})
		return
	}

	// 获取统计数据
	stats, err := h.statsRepo.GetSystemStats(c.Request.Context(), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": errors.Wrap(err, "获取系统统计数据失败").Error()})
		return
	}

	// 处理统计数据
	var totalRequests, totalUsers, activeUsers int64
	requestsChart := make([]map[string]interface{}, 0)
	usersChart := make([]map[string]interface{}, 0)

	for _, stat := range stats {
		totalRequests += stat.TotalRequests
		totalUsers = stat.TotalUsers // 使用最新的总用户数
		activeUsers = stat.ActiveUsers // 使用最新的活跃用户数

		requestsChart = append(requestsChart, map[string]interface{}{
			"date":  stat.Date.Format("2006-01-02"),
			"value": stat.TotalRequests,
		})

		usersChart = append(usersChart, map[string]interface{}{
			"date":  stat.Date.Format("2006-01-02"),
			"value": stat.ActiveUsers,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"totalRequests":  totalRequests,
			"totalUsers":     totalUsers,
			"activeUsers":    activeUsers,
			"requestsChart": requestsChart,
			"usersChart":    usersChart,
		},
	})
}

// GetAccuracyStats 获取识别准确率统计数据
func (h *StatsHandler) GetAccuracyStats(c *gin.Context) {
	modelID := c.Query("modelId")
	category := c.Query("category")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	// 解析时间参数
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "开始时间格式错误"})
		return
	}

	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "结束时间格式错误"})
		return
	}

	// 获取统计数据
	stats, err := h.statsRepo.GetAccuracyStats(c.Request.Context(), modelID, category, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": errors.Wrap(err, "获取准确率统计数据失败").Error()})
		return
	}

	// 处理统计数据
	var totalAccuracy float64
	categoryStats := make(map[string]float64)
	categoryCount := make(map[string]int)

	for _, stat := range stats {
		totalAccuracy += stat.Accuracy
		categoryStats[stat.Category] += stat.Accuracy
		categoryCount[stat.Category]++
	}

	// 计算平均准确率
	overallAccuracy := totalAccuracy / float64(len(stats))
	categoryAccuracy := make([]map[string]interface{}, 0)

	for category, total := range categoryStats {
		avg := total / float64(categoryCount[category])
		categoryAccuracy = append(categoryAccuracy, map[string]interface{}{
			"category": category,
			"accuracy": avg,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"overallAccuracy":   overallAccuracy,
			"categoryAccuracy": categoryAccuracy,
		},
	})
}

// GetResourceUsage 获取资源使用统计数据
func (h *StatsHandler) GetResourceUsage(c *gin.Context) {
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	// 解析时间参数
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "开始时间格式错误"})
		return
	}

	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "结束时间格式错误"})
		return
	}

	// 获取统计数据
	stats, err := h.statsRepo.GetResourceUsage(c.Request.Context(), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": errors.Wrap(err, "获取资源使用统计数据失败").Error()})
		return
	}

	// 处理统计数据
	var lastStat *model.ResourceUsage
	if len(stats) > 0 {
		lastStat = stats[len(stats)-1]
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"cpu": gin.H{
				"usage": lastStat.CPUUsage,
			},
			"memory": gin.H{
				"usage": lastStat.MemUsage,
			},
			"disk": gin.H{
				"usage": lastStat.DiskUsage,
			},
		},
	})
}

// GetResourceStats 获取资源使用统计
func (h *StatsHandler) GetResourceStats(c *gin.Context) {
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	// 解析时间参数
	start, _ := time.Parse(time.RFC3339, startTime)
	end, _ := time.Parse(time.RFC3339, endTime)

	// 获取资源使用统计数据
	_, err := h.statsRepo.GetResourceUsage(c.Request.Context(), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取资源使用统计数据失败"})
		return
	}

	// 处理统计数据 - 这里使用模拟数据
	cpuChart := make([]map[string]interface{}, 0)
	memoryChart := make([]map[string]interface{}, 0)
	diskChart := make([]map[string]interface{}, 0)

	// 模拟数据
	for i := 0; i < 24; i++ {
		time := time.Now().Add(time.Duration(-i) * time.Hour).Format("15:04")
		cpuChart = append(cpuChart, map[string]interface{}{
			"time":  time,
			"usage": 40.0 + float64(i%10)*5, // 示例数据
		})
		memoryChart = append(memoryChart, map[string]interface{}{
			"time":  time,
			"usage": 60.0 + float64(i%8)*3, // 示例数据
		})
		diskChart = append(diskChart, map[string]interface{}{
			"time":  time,
			"usage": 50.0 + float64(i%5), // 示例数据
		})
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"currentCPU":    45.2,
			"currentMemory": 68.7,
			"currentDisk":   52.3,
			"cpuChart":      cpuChart,
			"memoryChart":   memoryChart,
			"diskChart":     diskChart,
		},
	})
}