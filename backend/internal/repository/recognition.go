package repository

import (
	"context"
	"time"

	"github.com/image-recognition-engine/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RecognitionRepository 图像识别记录仓储接口
type RecognitionRepository interface {
	Create(ctx context.Context, record *model.RecognitionRecord) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.RecognitionRecord, error)
	GetByCustomerID(ctx context.Context, customerID int64, page, pageSize int) ([]*model.RecognitionRecord, int64, error)
	GetByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*model.RecognitionRecord, error)
	GetByModelVersion(ctx context.Context, modelVersion string) ([]*model.RecognitionRecord, error)
	GetStatsByModelVersion(ctx context.Context, modelVersion string) (float64, error)
	GetStatsByCategory(ctx context.Context, category string) (float64, error)
}

// ModelVersionRepository 模型版本仓储接口
type ModelVersionRepository interface {
	Create(ctx context.Context, version *model.ModelVersion) error
	Update(ctx context.Context, version *model.ModelVersion) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.ModelVersion, error)
	GetByVersion(ctx context.Context, version string) (*model.ModelVersion, error)
	List(ctx context.Context, page, pageSize int) ([]*model.ModelVersion, int64, error)
	GetLatestVersion(ctx context.Context) (*model.ModelVersion, error)
	GetByStatus(ctx context.Context, status string) ([]*model.ModelVersion, error)
}