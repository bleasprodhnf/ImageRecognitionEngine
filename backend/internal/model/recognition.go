package model

import (
	"time"
)

// Model 模型定义
type Model struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Version     string    `json:"version" bson:"version"`
	Type        string    `json:"type" bson:"type"` // 分类、检测、分割等
	Status      int       `json:"status" bson:"status"` // 0-未训练 1-训练中 2-已训练 3-已发布
	Accuracy    float64   `json:"accuracy" bson:"accuracy"` // 模型准确率
	FilePath    string    `json:"filePath" bson:"file_path"` // 模型文件路径
	UserID      int64     `json:"userId" bson:"user_id"` // 创建者ID
	CreateTime  time.Time `json:"createTime" bson:"create_time"`
	UpdateTime  time.Time `json:"updateTime" bson:"update_time"`
}

// RecognitionTask 识别任务
type RecognitionTask struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	ModelID     string    `json:"modelId" bson:"model_id"`
	UserID      int64     `json:"userId" bson:"user_id"`
	Status      int       `json:"status" bson:"status"` // 0-等待中 1-处理中 2-已完成 3-失败
	ImageURL    string    `json:"imageUrl" bson:"image_url"` // 原始图像URL
	ResultURL   string    `json:"resultUrl" bson:"result_url"` // 结果图像URL
	ResultData  string    `json:"resultData" bson:"result_data"` // JSON格式的识别结果数据
	ErrorMsg    string    `json:"errorMsg" bson:"error_msg"` // 错误信息
	ProcessTime int64     `json:"processTime" bson:"process_time"` // 处理时间(毫秒)
	CreateTime  time.Time `json:"createTime" bson:"create_time"`
	UpdateTime  time.Time `json:"updateTime" bson:"update_time"`
}

// Dataset 数据集
type Dataset struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Type        string    `json:"type" bson:"type"` // 训练集、验证集、测试集
	TotalCount  int       `json:"totalCount" bson:"total_count"` // 数据总数
	StoragePath string    `json:"storagePath" bson:"storage_path"` // 存储路径
	UserID      int64     `json:"userId" bson:"user_id"` // 创建者ID
	CreateTime  time.Time `json:"createTime" bson:"create_time"`
	UpdateTime  time.Time `json:"updateTime" bson:"update_time"`
}

// ModelRepository 模型数据访问接口
type ModelRepository interface {
	// 创建模型
	Create(model *Model) (string, error)
	// 更新模型
	Update(model *Model) error
	// 删除模型
	Delete(id string) error
	// 获取模型列表
	List(userID int64, page, size int) ([]*Model, int64, error)
	// 根据ID查找模型
	FindByID(id string) (*Model, error)
	// 更新模型状态
	UpdateStatus(id string, status int) error
	// 更新模型准确率
	UpdateAccuracy(id string, accuracy float64) error
}

// RecognitionTaskRepository 识别任务数据访问接口
type RecognitionTaskRepository interface {
	// 创建识别任务
	Create(task *RecognitionTask) (string, error)
	// 更新识别任务
	Update(task *RecognitionTask) error
	// 获取识别任务列表
	List(userID int64, modelID string, page, size int) ([]*RecognitionTask, int64, error)
	// 根据ID查找识别任务
	FindByID(id string) (*RecognitionTask, error)
	// 更新识别任务状态
	UpdateStatus(id string, status int, resultURL, resultData, errorMsg string, processTime int64) error
	// 获取用户的识别统计数据
	GetUserStats(userID int64, startTime, endTime time.Time) (map[string]interface{}, error)
}

// DatasetRepository 数据集数据访问接口
type DatasetRepository interface {
	// 创建数据集
	Create(dataset *Dataset) (string, error)
	// 更新数据集
	Update(dataset *Dataset) error
	// 删除数据集
	Delete(id string) error
	// 获取数据集列表
	List(userID int64, page, size int) ([]*Dataset, int64, error)
	// 根据ID查找数据集
	FindByID(id string) (*Dataset, error)
	// 更新数据集总数
	UpdateTotalCount(id string, totalCount int) error
}