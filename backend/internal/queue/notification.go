package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// NotificationType 定义通知类型
type NotificationType string

const (
	NotificationTypeTaskComplete NotificationType = "task_complete"
	NotificationTypeTaskFailed   NotificationType = "task_failed"
)

// Notification 定义通知结构
type Notification struct {
	ID        string           `json:"id"`
	Type      NotificationType `json:"type"`
	TaskID    string           `json:"task_id"`
	Message   string           `json:"message"`
	CreatedAt time.Time        `json:"created_at"`
}

// NotificationService 通知服务
type NotificationService struct {
	redis  *redis.Client
	prefix string
}

// NewNotificationService 创建新的通知服务
func NewNotificationService(client *redis.Client, prefix string) *NotificationService {
	return &NotificationService{
		redis:  client,
		prefix: prefix,
	}
}

// SendNotification 发送通知
func (ns *NotificationService) SendNotification(ctx context.Context, taskID string, nType NotificationType, message string) error {
	notification := &Notification{
		ID:        fmt.Sprintf("%s_%d", nType, time.Now().UnixNano()),
		Type:      nType,
		TaskID:    taskID,
		Message:   message,
		CreatedAt: time.Now(),
	}

	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("marshal notification error: %v", err)
	}

	// 将通知保存到Redis
	key := fmt.Sprintf("%s:notification:%s", ns.prefix, taskID)
	err = ns.redis.Set(ctx, key, notificationBytes, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("save notification error: %v", err)
	}

	// 发布通知事件
	pubKey := fmt.Sprintf("%s:notification:channel", ns.prefix)
	err = ns.redis.Publish(ctx, pubKey, notificationBytes).Err()
	if err != nil {
		log.Printf("publish notification error: %v", err)
	}

	return nil
}

// GetNotification 获取通知
func (ns *NotificationService) GetNotification(ctx context.Context, taskID string) (*Notification, error) {
	key := fmt.Sprintf("%s:notification:%s", ns.prefix, taskID)
	result, err := ns.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get notification error: %v", err)
	}

	var notification Notification
	if err := json.Unmarshal([]byte(result), &notification); err != nil {
		return nil, fmt.Errorf("unmarshal notification error: %v", err)
	}

	return &notification, nil
}