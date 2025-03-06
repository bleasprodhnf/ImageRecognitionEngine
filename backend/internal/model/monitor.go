package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ServerMetrics 服务器性能指标
type ServerMetrics struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp time.Time         `bson:"timestamp" json:"timestamp"`
	CPU       CPUMetrics        `bson:"cpu" json:"cpu"`
	Memory    MemoryMetrics     `bson:"memory" json:"memory"`
	Disk      DiskMetrics       `bson:"disk" json:"disk"`
	Network   NetworkMetrics    `bson:"network" json:"network"`
	Uptime    int64            `bson:"uptime" json:"uptime"`
}

// CPUMetrics CPU相关指标
type CPUMetrics struct {
	Usage float64 `bson:"usage" json:"usage"`
	Cores int     `bson:"cores" json:"cores"`
}

// MemoryMetrics 内存相关指标
type MemoryMetrics struct {
	Total int64 `bson:"total" json:"total"`
	Used  int64 `bson:"used" json:"used"`
	Free  int64 `bson:"free" json:"free"`
}

// DiskMetrics 磁盘相关指标
type DiskMetrics struct {
	Total int64 `bson:"total" json:"total"`
	Used  int64 `bson:"used" json:"used"`
	Free  int64 `bson:"free" json:"free"`
}

// NetworkMetrics 网络相关指标
type NetworkMetrics struct {
	In  int64 `bson:"in" json:"in"`
	Out int64 `bson:"out" json:"out"`
}

// APIMetrics API调用指标
type APIMetrics struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp    time.Time         `bson:"timestamp" json:"timestamp"`
	Endpoint     string            `bson:"endpoint" json:"endpoint"`
	Method       string            `bson:"method" json:"method"`
	ResponseTime float64           `bson:"response_time" json:"response_time"`
	StatusCode   int              `bson:"status_code" json:"status_code"`
	ClientIP     string            `bson:"client_ip" json:"client_ip"`
}

// StorageMetrics 存储空间使用指标
type StorageMetrics struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp       time.Time         `bson:"timestamp" json:"timestamp"`
	TotalSpace      int64             `bson:"total_space" json:"total_space"`
	UsedSpace       int64             `bson:"used_space" json:"used_space"`
	FreeSpace       int64             `bson:"free_space" json:"free_space"`
	ImageCount      int64             `bson:"image_count" json:"image_count"`
	AverageFileSize float64           `bson:"average_file_size" json:"average_file_size"`
}