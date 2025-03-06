package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/image-recognition-engine/internal/model"
)

type ModelRepository struct {
	db *mongo.Database
}

func NewModelRepository(db *mongo.Database) *ModelRepository {
	return &ModelRepository{db: db}
}

// CreateModelVersion 创建新的模型版本
func (r *ModelRepository) CreateModelVersion(ctx context.Context, version *model.ModelVersion) error {
	version.CreateTime = time.Now()
	version.UpdateTime = time.Now()
	_, err := r.db.Collection("model_versions").InsertOne(ctx, version)
	return errors.Wrap(err, "create model version failed")
}

// GetModelVersions 获取模型版本列表
func (r *ModelRepository) GetModelVersions(ctx context.Context, status string, page, pageSize int64) ([]*model.ModelVersion, int64, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "createTime", Value: -1}}).
		SetSkip((page - 1) * pageSize).
		SetLimit(pageSize)

	cursor, err := r.db.Collection("model_versions").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, errors.Wrap(err, "find model versions failed")
	}
	defer cursor.Close(ctx)

	var versions []*model.ModelVersion
	if err = cursor.All(ctx, &versions); err != nil {
		return nil, 0, errors.Wrap(err, "decode model versions failed")
	}

	total, err := r.db.Collection("model_versions").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "count model versions failed")
	}

	return versions, total, nil
}

// SaveModelPerformance 保存模型性能数据
func (r *ModelRepository) SaveModelPerformance(ctx context.Context, perf *model.ModelPerformance) error {
	perf.CreateTime = time.Now()
	_, err := r.db.Collection("model_performance").InsertOne(ctx, perf)
	return errors.Wrap(err, "save model performance failed")
}

// GetModelPerformance 获取模型性能数据
func (r *ModelRepository) GetModelPerformance(ctx context.Context, modelID primitive.ObjectID, startTime, endTime time.Time) (*model.ModelPerformance, error) {
	filter := bson.M{"modelId": modelID}
	if !startTime.IsZero() {
		filter["createTime"] = bson.M{"$gte": startTime}
	}
	if !endTime.IsZero() {
		filter["createTime"] = bson.M{"$lte": endTime}
	}

	var perf model.ModelPerformance
	err := r.db.Collection("model_performance").FindOne(ctx, filter).Decode(&perf)
	if err != nil {
		return nil, errors.Wrap(err, "get model performance failed")
	}

	return &perf, nil
}