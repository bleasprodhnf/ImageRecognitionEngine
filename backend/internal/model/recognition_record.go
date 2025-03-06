package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RecognitionRecord 图像识别记录模型
type RecognitionRecord struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CustomerID    int64             `json:"customerId" bson:"customer_id"`
	ModelVersion  string            `json:"modelVersion" bson:"model_version"`
	ImageURL      string            `json:"imageUrl" bson:"image_url"`
	ResultURL     string            `json:"resultUrl" bson:"result_url"`
	Category      string            `json:"category" bson:"category"`
	Confidence    float64           `json:"confidence" bson:"confidence"`
	ProcessTime   int64             `json:"processTime" bson:"process_time"` // 处理时间(毫秒)
	Status        int               `json:"status" bson:"status"`           // 0-处理中 1-成功 2-失败
	ErrorMessage  string            `json:"errorMessage" bson:"error_message"`
	CreateTime    time.Time         `json:"createTime" bson:"create_time"`
	UpdateTime    time.Time         `json:"updateTime" bson:"update_time"`
}