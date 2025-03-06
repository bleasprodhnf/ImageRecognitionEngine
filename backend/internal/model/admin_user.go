package model

import (
	"time"
)

// AdminUser u7ba1u7406u5458u7528u6237u6a21u578b
type AdminUser struct {
	ID            int64     `json:"id" db:"id"`
	Username      string    `json:"username" db:"username"`
	Password      string    `json:"password,omitempty" db:"password"`
	RealName      string    `json:"realName" db:"real_name"`
	Email         string    `json:"email" db:"email"`
	Phone         string    `json:"phone" db:"phone"`
	Avatar        string    `json:"avatar" db:"avatar"`
	RoleID        int64     `json:"roleId" db:"role_id"`
	Status        int       `json:"status" db:"status"` // 0-u7981u7528 1-u542fu7528
	LastLoginTime time.Time `json:"lastLoginTime" db:"last_login_time"`
	CreateTime    time.Time `json:"createTime" db:"create_time"`
	UpdateTime    time.Time `json:"updateTime" db:"update_time"`
}

// Customer u5ba2u6237u6a21u578b
type Customer struct {
	ID            int64     `json:"id" db:"id"`
	Username      string    `json:"username" db:"username"`
	Password      string    `json:"password,omitempty" db:"password"`
	APIKey        string    `json:"apiKey" db:"api_key"`
	Email         string    `json:"email" db:"email"`
	Phone         string    `json:"phone" db:"phone"`
	Company       string    `json:"company" db:"company"`
	Status        int       `json:"status" db:"status"` // 0-u7981u7528 1-u542fu7528
	PlanID        int64     `json:"planId" db:"plan_id"`
	ExpireTime    time.Time `json:"expireTime" db:"expire_time"`
	LastLogin     time.Time `json:"lastLogin" db:"last_login"`
	LastLoginTime time.Time `json:"lastLoginTime" db:"last_login_time"`
	CreateTime    time.Time `json:"createTime" db:"create_time"`
	UpdateTime    time.Time `json:"updateTime" db:"update_time"`
}

// Permission u6743u9650u6a21u578b
type Permission struct {
	ID          int64     `json:"id" db:"id"`
	Code        string    `json:"code" db:"code"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Module      string    `json:"module" db:"module"`
	CreateTime  time.Time `json:"createTime" db:"create_time"`
	UpdateTime  time.Time `json:"updateTime" db:"update_time"`
}