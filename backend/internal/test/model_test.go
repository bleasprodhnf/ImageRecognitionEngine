package test

import (
	"testing"
	"time"

	"github.com/image-recognition-engine/internal/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestModelVersion(t *testing.T) {
	// 创建测试数据
	modelVersion := model.ModelVersion{
		ID:          primitive.NewObjectID(),
		Version:     "1.0.0",
		Description: "测试版本",
		ReleaseDate: time.Now().Format("2006-01-02"),
		Status:      "development",
		Accuracy:    0.95,
		Parameters: model.ModelParameters{
			BatchSize:    32,
			LearningRate: 0.001,
			Epochs:       100,
		},
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	// 验证字段
	assert.NotEmpty(t, modelVersion.ID)
	assert.Equal(t, "1.0.0", modelVersion.Version)
	assert.Equal(t, "测试版本", modelVersion.Description)
	assert.Equal(t, "development", modelVersion.Status)
	assert.Equal(t, 0.95, modelVersion.Accuracy)
	assert.Equal(t, 32, modelVersion.Parameters.BatchSize)
	assert.Equal(t, 0.001, modelVersion.Parameters.LearningRate)
	assert.Equal(t, 100, modelVersion.Parameters.Epochs)
}

func TestRecognitionTask(t *testing.T) {
	// 创建测试数据
	task := model.RecognitionTask{
		ID:          "task_123",
		ModelID:     "model_456",
		UserID:      1001,
		Status:      0,
		ImageURL:    "http://example.com/test.jpg",
		ResultURL:   "",
		ResultData:  "",
		ErrorMsg:    "",
		ProcessTime: 0,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	// 验证字段
	assert.Equal(t, "task_123", task.ID)
	assert.Equal(t, "model_456", task.ModelID)
	assert.Equal(t, int64(1001), task.UserID)
	assert.Equal(t, 0, task.Status)
	assert.Equal(t, "http://example.com/test.jpg", task.ImageURL)
	assert.Empty(t, task.ResultURL)
	assert.Empty(t, task.ResultData)
	assert.Empty(t, task.ErrorMsg)
	assert.Equal(t, int64(0), task.ProcessTime)
}