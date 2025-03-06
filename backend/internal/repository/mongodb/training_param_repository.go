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
)

// TrainingParamRepositoryImpl 训练参数数据访问实现
type TrainingParamRepositoryImpl struct {
	collection *mongo.Collection
}

// NewTrainingParamRepository 创建训练参数数据访问实例
func NewTrainingParamRepository() model.TrainingParamRepository {
	return &TrainingParamRepositoryImpl{
		collection: database.MongoDB.Collection("training_params"),
	}
}

// Create 创建训练参数
func (r *TrainingParamRepositoryImpl) Create(param *model.TrainingParam) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 设置更新时间
	param.UpdateTime = time.Now()

	// 插入文档
	result, err := r.collection.InsertOne(ctx, param)
	if err != nil {
		return "", fmt.Errorf("创建训练参数失败: %w", err)
	}

	// 获取插入的ID
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// Update 更新训练参数
func (r *TrainingParamRepositoryImpl) Update(param *model.TrainingParam) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 设置更新时间
	param.UpdateTime = time.Now()

	// 转换ID为ObjectID
	objectID, err := primitive.ObjectIDFromHex(param.ID)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	// 更新文档
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"name":         param.Name,
		"display_name": param.DisplayName,
		"value":        param.Value,
		"type":         param.Type,
		"range":        param.Range,
		"description":  param.Description,
		"model_id":     param.ModelID,
		"update_time":  param.UpdateTime,
	}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新训练参数失败: %w", err)
	}

	return nil
}

// Delete 删除训练参数
func (r *TrainingParamRepositoryImpl) Delete(id string) error {
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
		return fmt.Errorf("删除训练参数失败: %w", err)
	}

	return nil
}

// List 获取训练参数列表
func (r *TrainingParamRepositoryImpl) List(modelID string) ([]*model.TrainingParam, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 查询指定模型的所有参数
	filter := bson.M{"model_id": modelID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查询训练参数列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	// 解析结果
	params := make([]*model.TrainingParam, 0)
	for cursor.Next(ctx) {
		var param model.TrainingParam
		if err := cursor.Decode(&param); err != nil {
			return nil, fmt.Errorf("解析训练参数失败: %w", err)
		}
		params = append(params, &param)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("遍历训练参数失败: %w", err)
	}

	return params, nil
}

// FindByID 根据ID查找训练参数
func (r *TrainingParamRepositoryImpl) FindByID(id string) (*model.TrainingParam, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 转换ID为ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID格式: %w", err)
	}

	// 查询文档
	filter := bson.M{"_id": objectID}
	var param model.TrainingParam
	err = r.collection.FindOne(ctx, filter).Decode(&param)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 参数不存在
		}
		return nil, fmt.Errorf("查询训练参数失败: %w", err)
	}

	return &param, nil
}

// FindByModelIDAndName 根据模型ID和参数名查找训练参数
func (r *TrainingParamRepositoryImpl) FindByModelIDAndName(modelID, name string) (*model.TrainingParam, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 查询文档
	filter := bson.M{"model_id": modelID, "name": name}
	var param model.TrainingParam
	err := r.collection.FindOne(ctx, filter).Decode(&param)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 参数不存在
		}
		return nil, fmt.Errorf("查询训练参数失败: %w", err)
	}

	return &param, nil
}