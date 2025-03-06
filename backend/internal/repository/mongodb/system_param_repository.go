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

// SystemParamRepositoryImpl 系统参数数据访问实现
type SystemParamRepositoryImpl struct {
	collection *mongo.Collection
}

// NewSystemParamRepository 创建系统参数数据访问实例
func NewSystemParamRepository() model.SystemParamRepository {
	return &SystemParamRepositoryImpl{
		collection: database.MongoDB.Collection("system_params"),
	}
}

// Create 创建系统参数
func (r *SystemParamRepositoryImpl) Create(param *model.SystemParam) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 设置创建和更新时间
	now := time.Now()
	param.CreateTime = now
	param.UpdateTime = now

	// 插入文档
	result, err := r.collection.InsertOne(ctx, param)
	if err != nil {
		return "", fmt.Errorf("创建系统参数失败: %w", err)
	}

	// 获取插入的ID
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// Update 更新系统参数
func (r *SystemParamRepositoryImpl) Update(param *model.SystemParam) error {
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
		"key":         param.Key,
		"value":       param.Value,
		"description": param.Description,
		"type":        param.Type,
		"update_time": param.UpdateTime,
	}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新系统参数失败: %w", err)
	}

	return nil
}

// Delete 删除系统参数
func (r *SystemParamRepositoryImpl) Delete(id string) error {
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
		return fmt.Errorf("删除系统参数失败: %w", err)
	}

	return nil
}

// List 获取系统参数列表
func (r *SystemParamRepositoryImpl) List() ([]*model.SystemParam, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 查询所有文档
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("查询系统参数列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	// 解析结果
	params := make([]*model.SystemParam, 0)
	for cursor.Next(ctx) {
		var param model.SystemParam
		if err := cursor.Decode(&param); err != nil {
			return nil, fmt.Errorf("解析系统参数失败: %w", err)
		}
		params = append(params, &param)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("遍历系统参数失败: %w", err)
	}

	return params, nil
}

// FindByID 根据ID查找系统参数
func (r *SystemParamRepositoryImpl) FindByID(id string) (*model.SystemParam, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 转换ID为ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID格式: %w", err)
	}

	// 查询文档
	filter := bson.M{"_id": objectID}
	var param model.SystemParam
	err = r.collection.FindOne(ctx, filter).Decode(&param)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 参数不存在
		}
		return nil, fmt.Errorf("查询系统参数失败: %w", err)
	}

	return &param, nil
}

// FindByKey 根据Key查找系统参数
func (r *SystemParamRepositoryImpl) FindByKey(key string) (*model.SystemParam, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 查询文档
	filter := bson.M{"key": key}
	var param model.SystemParam
	err := r.collection.FindOne(ctx, filter).Decode(&param)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 参数不存在
		}
		return nil, fmt.Errorf("查询系统参数失败: %w", err)
	}

	return &param, nil
}