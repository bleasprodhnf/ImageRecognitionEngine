package repository

import (
	"context"
	"time"

	"github.com/image-recognition-engine/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MonitorRepository 监控数据仓库
type MonitorRepository struct {
	db *mongo.Database
}

// NewMonitorRepository 创建监控数据仓库实例
func NewMonitorRepository(db *mongo.Database) *MonitorRepository {
	return &MonitorRepository{db: db}
}

// SaveServerMetrics 保存服务器监控指标
func (r *MonitorRepository) SaveServerMetrics(ctx context.Context, metrics *model.ServerMetrics) error {
	_, err := r.db.Collection("server_metrics").InsertOne(ctx, metrics)
	return err
}

// SaveAPIMetrics 保存API调用指标
func (r *MonitorRepository) SaveAPIMetrics(ctx context.Context, metrics *model.APIMetrics) error {
	_, err := r.db.Collection("api_metrics").InsertOne(ctx, metrics)
	return err
}

// SaveStorageMetrics 保存存储空间使用指标
func (r *MonitorRepository) SaveStorageMetrics(ctx context.Context, metrics *model.StorageMetrics) error {
	_, err := r.db.Collection("storage_metrics").InsertOne(ctx, metrics)
	return err
}

// GetServerMetrics 获取服务器监控指标
func (r *MonitorRepository) GetServerMetrics(ctx context.Context, start, end time.Time) ([]*model.ServerMetrics, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	filter := bson.M{
		"timestamp": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	cursor, err := r.db.Collection("server_metrics").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var metrics []*model.ServerMetrics
	if err = cursor.All(ctx, &metrics); err != nil {
		return nil, err
	}
	return metrics, nil
}

// GetAPIMetrics 获取API调用指标
func (r *MonitorRepository) GetAPIMetrics(ctx context.Context, start, end time.Time) ([]*model.APIMetrics, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	filter := bson.M{
		"timestamp": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	cursor, err := r.db.Collection("api_metrics").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var metrics []*model.APIMetrics
	if err = cursor.All(ctx, &metrics); err != nil {
		return nil, err
	}
	return metrics, nil
}

// GetStorageMetrics 获取存储空间使用指标
func (r *MonitorRepository) GetStorageMetrics(ctx context.Context, start, end time.Time) ([]*model.StorageMetrics, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	filter := bson.M{
		"timestamp": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	cursor, err := r.db.Collection("storage_metrics").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var metrics []*model.StorageMetrics
	if err = cursor.All(ctx, &metrics); err != nil {
		return nil, err
	}
	return metrics, nil
}