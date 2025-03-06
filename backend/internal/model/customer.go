package model

import (
	"time"
)

// Role 角色模型
type Role struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Permissions []string  `json:"permissions" db:"permissions"`
	CreateTime  time.Time `json:"createTime" db:"create_time"`
	UpdateTime  time.Time `json:"updateTime" db:"update_time"`
}

// ServicePlan 服务套餐模型
type ServicePlan struct {
	ID           int64     `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	Price        float64   `json:"price" db:"price"`
	RequestLimit int64     `json:"requestLimit" db:"request_limit"` // 每月请求次数限制
	Concurrent   int       `json:"concurrent" db:"concurrent"`      // 并发请求数限制
	Features     []string  `json:"features" db:"features"`        // 支持的特性
	CreateTime   time.Time `json:"createTime" db:"create_time"`
	UpdateTime   time.Time `json:"updateTime" db:"update_time"`
}

// CustomerUsage 客户使用情况统计
type CustomerUsage struct {
	ID           int64     `json:"id" db:"id"`
	CustomerID   int64     `json:"customerId" db:"customer_id"`
	RequestCount int64     `json:"requestCount" db:"request_count"` // 当月请求次数
	SuccessCount int64     `json:"successCount" db:"success_count"` // 成功次数
	FailCount    int64     `json:"failCount" db:"fail_count"`       // 失败次数
	AvgLatency   float64   `json:"avgLatency" db:"avg_latency"`    // 平均响应时间
	Month        time.Time `json:"month" db:"month"`              // 统计月份
	CreateTime   time.Time `json:"createTime" db:"create_time"`
	UpdateTime   time.Time `json:"updateTime" db:"update_time"`
}

// CustomerRepository 客户数据访问接口
type CustomerRepository interface {
	// 根据用户名查找客户
	FindByUsername(username string) (*Customer, error)
	// 创建客户
	Create(customer *Customer) (int64, error)
	// 更新客户
	Update(customer *Customer) error
	// 删除客户
	Delete(id int64) error
	// 获取客户列表
	List(page, size int) ([]*Customer, int64, error)
	// 根据ID查找客户
	FindByID(id int64) (*Customer, error)
	// 更新客户登录时间
	UpdateLastLogin(id int64) error
	// 更新客户使用统计
	UpdateUsage(usage *CustomerUsage) error
}

// ServicePlanRepository 服务套餐数据访问接口
type ServicePlanRepository interface {
	// 创建套餐
	Create(plan *ServicePlan) (int64, error)
	// 更新套餐
	Update(plan *ServicePlan) error
	// 删除套餐
	Delete(id int64) error
	// 获取套餐列表
	List(page, size int) ([]*ServicePlan, int64, error)
	// 根据ID查找套餐
	FindByID(id int64) (*ServicePlan, error)
}