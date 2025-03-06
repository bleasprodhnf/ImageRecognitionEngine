package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID         int64     `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	Password   string    `json:"password,omitempty" db:"password"`
	RealName   string    `json:"realName" db:"real_name"`
	Email      string    `json:"email" db:"email"`
	Phone      string    `json:"phone" db:"phone"`
	Avatar     string    `json:"avatar" db:"avatar"`
	RoleID     int64     `json:"roleId" db:"role_id"`
	Status     int       `json:"status" db:"status"` // 0-禁用 1-启用
	LastLogin  time.Time `json:"lastLogin" db:"last_login"`
	CreateTime time.Time `json:"createTime" db:"create_time"`
	UpdateTime time.Time `json:"updateTime" db:"update_time"`
}

// UserRepository 用户数据访问接口
type UserRepository interface {
	// 根据用户名查找用户
	FindByUsername(username string) (*User, error)
	// 创建用户
	Create(user *User) (int64, error)
	// 更新用户
	Update(user *User) error
	// 删除用户
	Delete(id int64) error
	// 获取用户列表
	List(page, size int) ([]*User, int64, error)
	// 根据ID查找用户
	FindByID(id int64) (*User, error)
	// 更新用户登录时间
	UpdateLastLogin(id int64) error
}

// RoleRepository 角色数据访问接口
type RoleRepository interface {
	// 创建角色
	Create(role *Role) (int64, error)
	// 更新角色
	Update(role *Role) error
	// 删除角色
	Delete(id int64) error
	// 获取角色列表
	List() ([]*Role, error)
	// 根据ID查找角色
	FindByID(id int64) (*Role, error)
}