package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/image-recognition-engine/internal/model"
)

// LogRepository 日志仓储接口
type LogRepository interface {
	CreateSystemLog(ctx context.Context, log *model.SystemLog) error
	CreateAPIUsageLog(ctx context.Context, log *model.APIUsageLog) error
	QuerySystemLogs(ctx context.Context, params model.LogQueryParams) (*model.LogResponse, error)
	QueryAPIUsageLogs(ctx context.Context, params model.LogQueryParams) (*model.LogResponse, error)
}

// mongoLogRepository MongoDB日志仓储实现
type mongoLogRepository struct {
	db             *mongo.Database
	systemLogColl  *mongo.Collection
	apiUsageLogColl *mongo.Collection
}

// NewMongoLogRepository 创建MongoDB日志仓储实例
func NewMongoLogRepository(db *mongo.Database) LogRepository {
	return &mongoLogRepository{
		db:             db,
		systemLogColl:  db.Collection("system_logs"),
		apiUsageLogColl: db.Collection("api_usage_logs"),
	}
}

// CreateSystemLog 创建系统日志
func (r *mongoLogRepository) CreateSystemLog(ctx context.Context, log *model.SystemLog) error {
	log.Timestamp = time.Now()
	_, err := r.systemLogColl.InsertOne(ctx, log)
	return err
}

// CreateAPIUsageLog 创建API使用日志
func (r *mongoLogRepository) CreateAPIUsageLog(ctx context.Context, log *model.APIUsageLog) error {
	log.Timestamp = time.Now()
	_, err := r.apiUsageLogColl.InsertOne(ctx, log)
	return err
}

// QuerySystemLogs 查询系统日志
func (r *mongoLogRepository) QuerySystemLogs(ctx context.Context, params model.LogQueryParams) (*model.LogResponse, error) {
	filter := bson.M{}

	// 构建查询条件
	if params.Keyword != "" {
		filter["$or"] = []bson.M{
			{"message": bson.M{"$regex": params.Keyword, "$options": "i"}},
			{"details": bson.M{"$regex": params.Keyword, "$options": "i"}},
		}
	}
	if params.Level != "" {
		filter["level"] = params.Level
	}
	if params.Module != "" {
		filter["module"] = params.Module
	}
	if !params.StartTime.IsZero() && !params.EndTime.IsZero() {
		filter["timestamp"] = bson.M{
			"$gte": params.StartTime,
			"$lte": params.EndTime,
		}
	}

	// 设置分页
	skip := int64((params.Page - 1) * params.PageSize)
	limit := int64(params.PageSize)

	// 设置排序
	opts := options.Find().SetSort(bson.M{"timestamp": -1}).SetSkip(skip).SetLimit(limit)

	// 查询总数
	total, err := r.systemLogColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 查询数据
	cursor, err := r.systemLogColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []model.SystemLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return &model.LogResponse{
		Total: int(total),
		Items: logs,
	}, nil
}

// QueryAPIUsageLogs 查询API使用日志
func (r *mongoLogRepository) QueryAPIUsageLogs(ctx context.Context, params model.LogQueryParams) (*model.LogResponse, error) {
	filter := bson.M{}

	// 构建查询条件
	if params.Keyword != "" {
		filter["$or"] = []bson.M{
			{"endpoint": bson.M{"$regex": params.Keyword, "$options": "i"}},
			{"customer_id": bson.M{"$regex": params.Keyword, "$options": "i"}},
		}
	}
	if !params.StartTime.IsZero() && !params.EndTime.IsZero() {
		filter["timestamp"] = bson.M{
			"$gte": params.StartTime,
			"$lte": params.EndTime,
		}
	}

	// 设置分页
	skip := int64((params.Page - 1) * params.PageSize)
	limit := int64(params.PageSize)

	// 设置排序
	opts := options.Find().SetSort(bson.M{"timestamp": -1}).SetSkip(skip).SetLimit(limit)

	// 查询总数
	total, err := r.apiUsageLogColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 查询数据
	cursor, err := r.apiUsageLogColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []model.APIUsageLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return &model.LogResponse{
		Total: int(total),
		Items: logs,
	}, nil
}