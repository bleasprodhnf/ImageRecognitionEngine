package model

import (
	"time"
)

// SystemLog 表示系统日志信息
type SystemLog struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Level     string    `json:"level" bson:"level"` // 日志级别：info, warning, error, fatal
	Module    string    `json:"module" bson:"module"` // 所属模块
	Message   string    `json:"message" bson:"message"` // 日志内容
	Details   string    `json:"details" bson:"details"` // 详细信息
	IP        string    `json:"ip" bson:"ip"` // 来源IP
	UserAgent string    `json:"userAgent" bson:"user_agent"` // 用户代理
	UserID    string    `json:"userId" bson:"user_id"` // 用户ID
	Timestamp time.Time `json:"timestamp" bson:"timestamp"` // 创建时间
}

// LogQueryParams 表示日志查询参数
type LogQueryParams struct {
	Page     int       `json:"page" form:"page"`
	PageSize int       `json:"pageSize" form:"pageSize"`
	Keyword  string    `json:"keyword" form:"keyword"`
	Level    string    `json:"level" form:"level"`
	Module   string    `json:"module" form:"module"`
	StartTime time.Time `json:"startTime" form:"startTime"`
	EndTime   time.Time `json:"endTime" form:"endTime"`
}

// LogResponse 表示日志查询响应
type LogResponse struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

// APIUsageLog 表示API使用日志
type APIUsageLog struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	CustomerID   string    `json:"customerId" bson:"customer_id"`
	Endpoint     string    `json:"endpoint" bson:"endpoint"`
	Method       string    `json:"method" bson:"method"`
	StatusCode   int       `json:"statusCode" bson:"status_code"`
	ResponseTime int64     `json:"responseTime" bson:"response_time"` // 响应时间(ms)
	RequestSize  int64     `json:"requestSize" bson:"request_size"`   // 请求大小(bytes)
	ResponseSize int64     `json:"responseSize" bson:"response_size"` // 响应大小(bytes)
	IP           string    `json:"ip" bson:"ip"`
	UserAgent    string    `json:"userAgent" bson:"user_agent"`
	Timestamp    time.Time `json:"timestamp" bson:"timestamp"`
}

// APIUsageStats 表示API使用统计
type APIUsageStats struct {
	Endpoint     string `json:"endpoint" bson:"_id"`
	Count        int64  `json:"count" bson:"count"`
	AvgResTime   int64  `json:"avgResTime" bson:"avg_res_time"`
	MinResTime   int64  `json:"minResTime" bson:"min_res_time"`
	MaxResTime   int64  `json:"maxResTime" bson:"max_res_time"`
	TotalReqSize int64  `json:"totalReqSize" bson:"total_req_size"`
	TotalResSize int64  `json:"totalResSize" bson:"total_res_size"`
}