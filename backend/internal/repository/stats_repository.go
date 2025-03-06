package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/image-recognition-engine/internal/model"
)

type StatsRepository struct {
	db *mongo.Database
}

func NewStatsRepository(db *mongo.Database) *StatsRepository {
	return &StatsRepository{db: db}
}

// GetSystemStats 获取系统使用统计数据
func (r *StatsRepository) GetSystemStats(ctx context.Context, startTime, endTime time.Time) ([]*model.SystemStats, error) {
	coll := r.db.Collection("system_stats")

	filter := bson.M{
		"date": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	opts := options.Find().SetSort(bson.M{"date": 1})

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.Wrap(err, "查询系统统计数据失败")
	}
	defer cursor.Close(ctx)

	var stats []*model.SystemStats
	if err = cursor.All(ctx, &stats); err != nil {
		return nil, errors.Wrap(err, "解析系统统计数据失败")
	}

	return stats, nil
}

// GetAccuracyStats 获取识别准确率统计数据
func (r *StatsRepository) GetAccuracyStats(ctx context.Context, modelID string, category string, startTime, endTime time.Time) ([]*model.AccuracyStats, error) {
	coll := r.db.Collection("accuracy_stats")

	filter := bson.M{
		"date": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	if modelID != "" {
		filter["model_id"] = modelID
	}
	if category != "" {
		filter["category"] = category
	}

	opts := options.Find().SetSort(bson.M{"date": 1})

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.Wrap(err, "查询准确率统计数据失败")
	}
	defer cursor.Close(ctx)

	var stats []*model.AccuracyStats
	if err = cursor.All(ctx, &stats); err != nil {
		return nil, errors.Wrap(err, "解析准确率统计数据失败")
	}

	return stats, nil
}

// GetResourceUsage 获取资源使用统计数据
func (r *StatsRepository) GetResourceUsage(ctx context.Context, startTime, endTime time.Time) ([]*model.ResourceUsage, error) {
	coll := r.db.Collection("resource_usage")

	filter := bson.M{
		"date": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	opts := options.Find().SetSort(bson.M{"date": 1})

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.Wrap(err, "查询资源使用统计数据失败")
	}
	defer cursor.Close(ctx)

	var stats []*model.ResourceUsage
	if err = cursor.All(ctx, &stats); err != nil {
		return nil, errors.Wrap(err, "解析资源使用统计数据失败")
	}

	return stats, nil
}

// SaveSystemStats 保存系统使用统计数据
func (r *StatsRepository) SaveSystemStats(ctx context.Context, stats *model.SystemStats) error {
	coll := r.db.Collection("system_stats")

	stats.CreatedAt = time.Now()
	_, err := coll.InsertOne(ctx, stats)
	if err != nil {
		return errors.Wrap(err, "保存系统统计数据失败")
	}

	return nil
}

// SaveAccuracyStats 保存识别准确率统计数据
func (r *StatsRepository) SaveAccuracyStats(ctx context.Context, stats *model.AccuracyStats) error {
	coll := r.db.Collection("accuracy_stats")

	stats.CreatedAt = time.Now()
	_, err := coll.InsertOne(ctx, stats)
	if err != nil {
		return errors.Wrap(err, "保存准确率统计数据失败")
	}

	return nil
}

// SaveResourceUsage 保存资源使用统计数据
func (r *StatsRepository) SaveResourceUsage(ctx context.Context, stats *model.ResourceUsage) error {
	coll := r.db.Collection("resource_usage")

	stats.CreatedAt = time.Now()
	_, err := coll.InsertOne(ctx, stats)
	if err != nil {
		return errors.Wrap(err, "保存资源使用统计数据失败")
	}

	return nil
}