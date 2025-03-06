package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/image-recognition-engine/internal/database"
	"github.com/image-recognition-engine/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ModelRepositoryImpl 模型数据访问实现
type ModelRepositoryImpl struct {
	collection *mongo.Collection
}

// NewModelRepository 创建模型数据访问实例
func NewModelRepository() model.ModelRepository {
	return &ModelRepositoryImpl{
		collection: database.MongoDB.Collection("models"),
	}
}

// Create 创建模型
func (r *ModelRepositoryImpl) Create(model *model.Model) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 设置创建和更新时间
	now := time.Now()
	model.CreateTime = now
	model.UpdateTime = now

	// 插入文档
	result, err := r.collection.InsertOne(ctx, model)
	if err != nil {
		return "", fmt.Errorf("创建模型失败: %w", err)
	}

	// 获取插入的ID
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// Update 更新模型
func (r *ModelRepositoryImpl) Update(model *model.Model) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 设置更新时间
	model.UpdateTime = time.Now()

	// 转换ID为ObjectID
	objectID, err := primitive.ObjectIDFromHex(model.ID)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	// 更新文档
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"name":        model.Name,
		"description": model.Description,
		"version":     model.Version,
		"type":        model.Type,
		"status":      model.Status,
		"accuracy":    model.Accuracy,
		"file_path":   model.FilePath,
		"update_time": model.UpdateTime,
	}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新模型失败: %w", err)
	}

	return nil
}

// Delete 删除模型
func (r *ModelRepositoryImpl) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 转换ID为ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	// 删除文档
	filter := bson.M{"_id": objectID}
	_, err = r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除模型失败: %w", err)
	}

	return nil
}

// List 获取模型列表
func (r *ModelRepositoryImpl) List(userID int64, page, size int) ([]*model.Model, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建查询条件
	filter := bson.M{}
	if userID > 0 {
		filter["user_id"] = userID
	}

	// 计算总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("获取模型总数失败: %w", err)
	}

	// 分页查询
	opts := options.Find()
	opts.SetSort(bson.M{"create_time": -1}) // 按创建时间降序
	opts.SetSkip(int64((page - 1) * size))
	opts.SetLimit(int64(size))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询模型列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	// 解析结果
	models := make([]*model.Model, 0)
	for cursor.Next(ctx) {
		var m model.Model
		if err := cursor.Decode(&m); err != nil {
			return nil, 0, fmt.Errorf("解析模型数据失败: %w", err)
		}
		models = append(models, &m)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("遍历模型数据失败: %w", err)
	}

	return models, total, nil
}

// FindByID 根据ID查找模型
func (r *ModelRepositoryImpl) FindByID(id string) (*model.Model, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 转换ID为ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID格式: %w", err)
	}

	// 查询文档
	filter := bson.M{"_id": objectID}
	var m model.Model
	err = r.collection.FindOne(ctx, filter).Decode(&m)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 模型不存在
		}
		return nil, fmt.Errorf("查询模型失败: %w", err)
	}

	return &m, nil
}

// UpdateStatus 更新模型状态
func (r *ModelRepositoryImpl) UpdateStatus(id string, status int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 转换ID为ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	// 更新文档
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"status":      status,
		"update_time": time.Now(),
	}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新模型状态失败: %w", err)
	}

	return nil
}

// UpdateAccuracy 更新模型准确率
func (r *ModelRepositoryImpl) UpdateAccuracy(id string, accuracy float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 转换ID为ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	// 更新文档
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"accuracy":    accuracy,
		"update_time": time.Now(),
	}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新模型准确率失败: %w", err)
	}

	return nil
}