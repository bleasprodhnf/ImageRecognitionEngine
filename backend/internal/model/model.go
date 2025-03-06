package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ModelVersion 模型版本信息
type ModelVersion struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Version     string            `bson:"version" json:"version"`
	Description string            `bson:"description" json:"description"`
	ReleaseDate string            `bson:"releaseDate" json:"releaseDate"`
	Status      string            `bson:"status" json:"status"` // development/testing/production
	Accuracy    float64           `bson:"accuracy" json:"accuracy"`
	Parameters  ModelParameters   `bson:"parameters" json:"parameters"`
	CreateTime  time.Time         `bson:"createTime" json:"createTime"`
	UpdateTime  time.Time         `bson:"updateTime" json:"updateTime"`
}

// ModelParameters 模型训练参数
type ModelParameters struct {
	BatchSize    int     `bson:"batchSize" json:"batchSize"`
	LearningRate float64 `bson:"learningRate" json:"learningRate"`
	Epochs       int     `bson:"epochs" json:"epochs"`
}

// ModelPerformance 模型性能数据
type ModelPerformance struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ModelID     primitive.ObjectID  `bson:"modelId" json:"modelId"`
	ModelVersion string             `bson:"modelVersion" json:"modelVersion"`
	Metrics      ModelMetrics       `bson:"metrics" json:"metrics"`
	CreateTime   time.Time          `bson:"createTime" json:"createTime"`
}

// ModelMetrics 模型性能指标
type ModelMetrics struct {
	Accuracy   []MetricPoint `bson:"accuracy" json:"accuracy"`
	Latency    []MetricPoint `bson:"latency" json:"latency"`
	Throughput []MetricPoint `bson:"throughput" json:"throughput"`
}

// MetricPoint 性能指标数据点
type MetricPoint struct {
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Value     float64   `bson:"value" json:"value"`
}