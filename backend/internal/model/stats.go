package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SystemStats 系统使用统计
type SystemStats struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TotalRequests int64             `bson:"total_requests" json:"totalRequests"`
	TotalUsers    int64             `bson:"total_users" json:"totalUsers"`
	ActiveUsers   int64             `bson:"active_users" json:"activeUsers"`
	Timestamp     time.Time         `bson:"timestamp" json:"timestamp"`
	Date          time.Time         `bson:"date" json:"date"`
	CreatedAt     time.Time         `bson:"created_at" json:"createdAt"`
}

// AccuracyStats 识别准确率统计
type AccuracyStats struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ModelID         string            `bson:"model_id" json:"modelId"`
	Category        string            `bson:"category" json:"category"`
	Accuracy        float64           `bson:"accuracy" json:"accuracy"`
	TotalSamples    int64             `bson:"total_samples" json:"totalSamples"`
	CorrectSamples  int64             `bson:"correct_samples" json:"correctSamples"`
	Date            time.Time         `bson:"date" json:"date"`
	CreatedAt       time.Time         `bson:"created_at" json:"createdAt"`
}

// ResourceUsage 资源使用统计
type ResourceUsage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CPUUsage  float64           `bson:"cpu_usage" json:"cpuUsage"`
	MemUsage  float64           `bson:"mem_usage" json:"memUsage"`
	DiskUsage float64           `bson:"disk_usage" json:"diskUsage"`
	Date      time.Time         `bson:"date" json:"date"`
	CreatedAt time.Time         `bson:"created_at" json:"createdAt"`
}