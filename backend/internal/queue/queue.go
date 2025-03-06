package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// TaskType 定义任务类型
type TaskType string

const (
	TaskTypeImageRecognition TaskType = "image_recognition"
	TaskTypeModelTraining    TaskType = "model_training"
	TaskTypeDataAnalysis     TaskType = "data_analysis"
)

// Task 定义任务结构
type Task struct {
	ID        string          `json:"id"`
	Type      TaskType        `json:"type"`
	Data      json.RawMessage `json:"data"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Queue 消息队列管理器
type Queue struct {
	redis  *redis.Client
	prefix string
}

// NewQueue 创建新的队列管理器
func NewQueue(client *redis.Client, prefix string) *Queue {
	return &Queue{
		redis:  client,
		prefix: prefix,
	}
}

// Push 将任务推送到队列
func (q *Queue) Push(ctx context.Context, taskType TaskType, data interface{}) error {
	taskData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal task data error: %v", err)
	}

	task := &Task{
		ID:        fmt.Sprintf("%s_%d", taskType, time.Now().UnixNano()),
		Type:      taskType,
		Data:      taskData,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	taskBytes, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal task error: %v", err)
	}

	key := fmt.Sprintf("%s:%s", q.prefix, taskType)
	err = q.redis.LPush(ctx, key, taskBytes).Err()
	if err != nil {
		return fmt.Errorf("push task error: %v", err)
	}

	return nil
}

// Pop 从队列中获取任务
func (q *Queue) Pop(ctx context.Context, taskType TaskType) (*Task, error) {
	key := fmt.Sprintf("%s:%s", q.prefix, taskType)
	result, err := q.redis.RPop(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("pop task error: %v", err)
	}

	var task Task
	if err := json.Unmarshal([]byte(result), &task); err != nil {
		return nil, fmt.Errorf("unmarshal task error: %v", err)
	}

	return &task, nil
}

// UpdateTaskStatus 更新任务状态
func (q *Queue) UpdateTaskStatus(ctx context.Context, taskID string, status string) error {
	key := fmt.Sprintf("%s:status:%s", q.prefix, taskID)
	err := q.redis.Set(ctx, key, status, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("update task status error: %v", err)
	}

	return nil
}

// GetTaskStatus 获取任务状态
func (q *Queue) GetTaskStatus(ctx context.Context, taskID string) (string, error) {
	key := fmt.Sprintf("%s:status:%s", q.prefix, taskID)
	status, err := q.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get task status error: %v", err)
	}

	return status, nil
}