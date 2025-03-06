package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// User 用户结构
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	RealName     string    `json:"realName"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Status       int       `json:"status"`
	Roles        []Role    `json:"roles"`
	LastLoginTime time.Time `json:"lastLoginTime"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Role 角色结构
type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UserCreateRequest 创建用户请求结构
type UserCreateRequest struct {
	Username string   `json:"username" binding:"required"`
	Password string   `json:"password" binding:"required"`
	RealName string   `json:"realName"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Status   int      `json:"status"`
	RoleIDs  []int64  `json:"roleIds"`
}

// UserUpdateRequest 更新用户请求结构
type UserUpdateRequest struct {
	RealName string   `json:"realName"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Status   int      `json:"status"`
	RoleIDs  []int64  `json:"roleIds"`
	Password string   `json:"password,omitempty"`
}

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// TODO: 从数据库查询用户列表
	// 这里返回示例数据
	users := []User{
		{
			ID:           1,
			Username:     "admin",
			RealName:     "系统管理员",
			Email:        "admin@example.com",
			Phone:        "13800138000",
			Status:       1,
			Roles:        []Role{{ID: 1, Name: "超级管理员", Description: "拥有所有权限"}},
			LastLoginTime: time.Now().Add(-24 * time.Hour),
			CreatedAt:    time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:    time.Now().Add(-7 * 24 * time.Hour),
		},
	}

	// 构建分页信息
	pagination := map[string]interface{}{
		"total":     1,
		"page":      page,
		"pageSize":  pageSize,
		"totalPage": 1,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    users,
		"pagination": pagination,
	})
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var req UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// TODO: 验证用户名是否已存在
	// TODO: 密码加密
	// TODO: 保存用户信息到数据库
	// TODO: 关联用户角色

	// 构建响应
	user := User{
		ID:           1,
		Username:     req.Username,
		RealName:     req.RealName,
		Email:        req.Email,
		Phone:        req.Phone,
		Status:       req.Status,
		Roles:        []Role{{ID: 1, Name: "管理员", Description: "管理员权限"}},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    user,
	})
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID无效",
			"data":    nil,
		})
		return
	}

	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// TODO: 验证用户是否存在
	// TODO: 如果更新密码，进行密码加密
	// TODO: 更新用户信息到数据库
	// TODO: 更新用户角色关联

	// 构建响应
	user := User{
		ID:           id,
		Username:     "admin", // 用户名不允许修改，这里使用示例值
		RealName:     req.RealName,
		Email:        req.Email,
		Phone:        req.Phone,
		Status:       req.Status,
		Roles:        []Role{{ID: 1, Name: "管理员", Description: "管理员权限"}},
		LastLoginTime: time.Now().Add(-24 * time.Hour),
		CreatedAt:    time.Now().Add(-30 * 24 * time.Hour),
		UpdatedAt:    time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    user,
	})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	_, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID无效",
			"data":    nil,
		})
		return
	}

	// TODO: 验证用户是否存在
	// TODO: 从数据库删除用户及关联信息

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data":    nil,
	})
}