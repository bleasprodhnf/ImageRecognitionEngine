package security

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuditLogService 提供安全审计日志功能
type AuditLogService struct {
	mongoDB     *mongo.Database
	redisClient *redis.Client
	dataMasking *DataMaskingService
	enabled     bool
}

// AuditLogEntry 审计日志条目
type AuditLogEntry struct {
	ID          string                 `json:"id" bson:"_id,omitempty"`
	UserID      int64                  `json:"userId" bson:"userId"`
	Username    string                 `json:"username" bson:"username"`
	Action      string                 `json:"action" bson:"action"`
	Resource    string                 `json:"resource" bson:"resource"`
	ResourceID  string                 `json:"resourceId" bson:"resourceId"`
	ClientIP    string                 `json:"clientIP" bson:"clientIP"`
	UserAgent   string                 `json:"userAgent" bson:"userAgent"`
	Status      string                 `json:"status" bson:"status"` // success, failure
	Details     map[string]interface{} `json:"details" bson:"details"`
	Timestamp   time.Time              `json:"timestamp" bson:"timestamp"`
	Environment string                 `json:"environment" bson:"environment"`
}

// NewAuditLogService 创建一个新的审计日志服务实例
func NewAuditLogService(mongoDB *mongo.Database, redisClient *redis.Client, dataMasking *DataMaskingService, enabled bool) *AuditLogService {
	return &AuditLogService{
		mongoDB:     mongoDB,
		redisClient: redisClient,
		dataMasking: dataMasking,
		enabled:     enabled,
	}
}

// LogUserAction 记录用户操作
func (s *AuditLogService) LogUserAction(ctx context.Context, userID int64, username, action, resource, resourceID, clientIP, userAgent, status string, details map[string]interface{}) error {
	if !s.enabled {
		return nil
	}

	// 创建审计日志条目
	entry := AuditLogEntry{
		UserID:      userID,
		Username:    username,
		Action:      action,
		Resource:    resource,
		ResourceID:  resourceID,
		ClientIP:    clientIP,
		UserAgent:   userAgent,
		Status:      status,
		Details:     s.maskSensitiveDetails(details),
		Timestamp:   time.Now(),
		Environment: "production", // 可以从配置中获取
	}

	// 将日志保存到MongoDB
	_, err := s.mongoDB.Collection("audit_logs").InsertOne(ctx, entry)
	if err != nil {
		// 如果MongoDB写入失败，尝试写入Redis作为备份
		s.backupToRedis(entry)
		return fmt.Errorf("保存审计日志到MongoDB失败: %w", err)
	}

	return nil
}

// maskSensitiveDetails 对敏感详情进行脱敏处理
func (s *AuditLogService) maskSensitiveDetails(details map[string]interface{}) map[string]interface{} {
	if s.dataMasking == nil {
		return details
	}

	// 将详情转换为JSON字符串
	jsonData, err := json.Marshal(details)
	if err != nil {
		return details
	}

	// 对JSON字符串进行脱敏处理
	maskedJSON := s.dataMasking.MaskJSON(string(jsonData))

	// 将脱敏后的JSON字符串转换回map
	var maskedDetails map[string]interface{}
	if err := json.Unmarshal([]byte(maskedJSON), &maskedDetails); err != nil {
		return details
	}

	return maskedDetails
}

// backupToRedis 将审计日志备份到Redis
func (s *AuditLogService) backupToRedis(entry AuditLogEntry) error {
	// 将日志条目转换为JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	// 生成唯一键
	key := fmt.Sprintf("audit_log_backup:%d:%s", entry.UserID, entry.Timestamp.Format(time.RFC3339))

	// 存储到Redis，设置过期时间为24小时
	ctx := context.Background()
	err = s.redisClient.Set(ctx, key, string(jsonData), 24*time.Hour).Err()
	return err
}

// GetUserLogs 获取用户操作日志
func (s *AuditLogService) GetUserLogs(ctx context.Context, userID int64, startTime, endTime time.Time, page, pageSize int) ([]AuditLogEntry, int64, error) {
	// 构建查询条件
	filter := bson.M{"userId": userID}

	// 添加时间范围条件
	if !startTime.IsZero() || !endTime.IsZero() {
		timeFilter := bson.M{}
		if !startTime.IsZero() {
			timeFilter["$gte"] = startTime
		}
		if !endTime.IsZero() {
			timeFilter["$lte"] = endTime
		}
		filter["timestamp"] = timeFilter
	}

	// 计算总数
	total, err := s.mongoDB.Collection("audit_logs").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 查询日志
	options := bson.D{
		{"$sort", bson.D{{"timestamp", -1}}},
		{"$skip", int64((page - 1) * pageSize)},
		{"$limit", int64(pageSize)},
	}

	cursor, err := s.mongoDB.Collection("audit_logs").Aggregate(ctx, []bson.D{{
		{"$match", filter},
	}, options})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// 解析结果
	var logs []AuditLogEntry
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetResourceLogs 获取资源操作日志
func (s *AuditLogService) GetResourceLogs(ctx context.Context, resource, resourceID string, startTime, endTime time.Time, page, pageSize int) ([]AuditLogEntry, int64, error) {
	// 构建查询条件
	filter := bson.M{}
	if resource != "" {
		filter["resource"] = resource
	}
	if resourceID != "" {
		filter["resourceId"] = resourceID
	}

	// 添加时间范围条件
	if !startTime.IsZero() || !endTime.IsZero() {
		timeFilter := bson.M{}
		if !startTime.IsZero() {
			timeFilter["$gte"] = startTime
		}
		if !endTime.IsZero() {
			timeFilter["$lte"] = endTime
		}
		filter["timestamp"] = timeFilter
	}

	// 计算总数
	total, err := s.mongoDB.Collection("audit_logs").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 查询日志
	options := bson.D{
		{"$sort", bson.D{{"timestamp", -1}}},
		{"$skip", int64((page - 1) * pageSize)},
		{"$limit", int64(pageSize)},
	}

	cursor, err := s.mongoDB.Collection("audit_logs").Aggregate(ctx, []bson.D{{
		{"$match", filter},
	}, options})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// 解析结果
	var logs []AuditLogEntry
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}