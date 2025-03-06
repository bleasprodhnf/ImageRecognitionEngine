package model

import (
	"time"
)

// SystemParam 系统参数模型
type SystemParam struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Key         string    `json:"key" bson:"key"`
	Value       string    `json:"value" bson:"value"`
	Description string    `json:"description" bson:"description"`
	Type        string    `json:"type" bson:"type"` // string, number, boolean, json
	CreateTime  time.Time `json:"createTime" bson:"create_time"`
	UpdateTime  time.Time `json:"updateTime" bson:"update_time"`
}

// TrainingParam 训练参数模型
type TrainingParam struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	DisplayName string    `json:"displayName" bson:"display_name"`
	Value       string    `json:"value" bson:"value"`
	Type        string    `json:"type" bson:"type"` // string, number, boolean, json
	Range       string    `json:"range" bson:"range"` // 可选值范围，如"32,64,128,256"
	Description string    `json:"description" bson:"description"`
	ModelID     string    `json:"modelId" bson:"model_id"`
	UpdateTime  time.Time `json:"updateTime" bson:"update_time"`
}

// SystemParamRepository 系统参数数据访问接口
type SystemParamRepository interface {
	// 创建系统参数
	Create(param *SystemParam) (string, error)
	// 更新系统参数
	Update(param *SystemParam) error
	// 删除系统参数
	Delete(id string) error
	// 获取系统参数列表
	List() ([]*SystemParam, error)
	// 根据ID查找系统参数
	FindByID(id string) (*SystemParam, error)
	// 根据Key查找系统参数
	FindByKey(key string) (*SystemParam, error)
}

// TrainingParamRepository 训练参数数据访问接口
type TrainingParamRepository interface {
	// 创建训练参数
	Create(param *TrainingParam) (string, error)
	// 更新训练参数
	Update(param *TrainingParam) error
	// 删除训练参数
	Delete(id string) error
	// 获取训练参数列表
	List(modelID string) ([]*TrainingParam, error)
	// 根据ID查找训练参数
	FindByID(id string) (*TrainingParam, error)
	// 根据模型ID和参数名查找训练参数
	FindByModelIDAndName(modelID, name string) (*TrainingParam, error)
}