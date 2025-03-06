package handler

import (
	"fmt"

	"github.com/image-recognition-engine/internal/model"
)

// calculateTotalRequests 计算总请求数
func calculateTotalRequests(stats []*model.SystemStats) int64 {
	var total int64
	for _, stat := range stats {
		total += stat.TotalRequests
	}
	return total
}

// calculateTotalUsers 计算总用户数
func calculateTotalUsers(stats []*model.SystemStats) int64 {
	if len(stats) == 0 {
		return 0
	}
	return stats[len(stats)-1].TotalUsers
}

// calculateActiveUsers 计算活跃用户数
func calculateActiveUsers(stats []*model.SystemStats) int64 {
	if len(stats) == 0 {
		return 0
	}
	return stats[len(stats)-1].ActiveUsers
}

// generateRequestsChart 生成请求数量图表数据
func generateRequestsChart(stats []*model.SystemStats, interval string) []map[string]interface{} {
	var result []map[string]interface{}
	if len(stats) == 0 {
		return result
	}

	// 根据interval参数聚合数据
	groupedStats := groupStatsByInterval(stats, interval)

	// 生成图表数据
	for date, requests := range groupedStats {
		result = append(result, map[string]interface{}{
			"date":  date,
			"value": requests,
		})
	}

	return result
}

// generateUsersChart 生成用户数量图表数据
func generateUsersChart(stats []*model.SystemStats, interval string) []map[string]interface{} {
	var result []map[string]interface{}
	if len(stats) == 0 {
		return result
	}

	// 根据interval参数聚合数据
	groupedStats := groupUsersByInterval(stats, interval)

	// 生成图表数据
	for date, users := range groupedStats {
		result = append(result, map[string]interface{}{
			"date":  date,
			"value": users,
		})
	}

	return result
}

// groupStatsByInterval 按时间间隔分组统计数据
func groupStatsByInterval(stats []*model.SystemStats, interval string) map[string]int64 {
	result := make(map[string]int64)

	for _, stat := range stats {
		var dateKey string
		switch interval {
		case "hour":
			dateKey = stat.Timestamp.Format("2006-01-02 15:00")
		case "day":
			dateKey = stat.Timestamp.Format("2006-01-02")
		case "week":
			year, week := stat.Timestamp.ISOWeek()
			dateKey = fmt.Sprintf("%d-W%02d", year, week)
		case "month":
			dateKey = stat.Timestamp.Format("2006-01")
		default:
			dateKey = stat.Timestamp.Format("2006-01-02")
		}

		result[dateKey] += stat.TotalRequests
	}

	return result
}

// groupUsersByInterval 按时间间隔分组用户数据
func groupUsersByInterval(stats []*model.SystemStats, interval string) map[string]int64 {
	result := make(map[string]int64)

	for _, stat := range stats {
		var dateKey string
		switch interval {
		case "hour":
			dateKey = stat.Timestamp.Format("2006-01-02 15:00")
		case "day":
			dateKey = stat.Timestamp.Format("2006-01-02")
		case "week":
			year, week := stat.Timestamp.ISOWeek()
			dateKey = fmt.Sprintf("%d-W%02d", year, week)
		case "month":
			dateKey = stat.Timestamp.Format("2006-01")
		default:
			dateKey = stat.Timestamp.Format("2006-01-02")
		}

		// 使用活跃用户数作为用户图表数据
		result[dateKey] = stat.ActiveUsers
	}

	return result
}