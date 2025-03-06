package repository

import (
	"github.com/image-recognition-engine/internal/model"
)

// AdminUserRepository 管理员用户仓储接口
type AdminUserRepository interface {
	GetByID(id int64) (*model.AdminUser, error)
	GetByUsername(username string) (*model.AdminUser, error)
	Create(user *model.AdminUser) error
	Update(user *model.AdminUser) error
	Delete(id int64) error
	List(page, pageSize int) ([]*model.AdminUser, int64, error)
}

// CustomerRepository 客户仓储接口
type CustomerRepository interface {
	GetByID(id int64) (*model.Customer, error)
	GetByUsername(username string) (*model.Customer, error)
	GetByAPIKey(apiKey string) (*model.Customer, error)
	Create(customer *model.Customer) error
	Update(customer *model.Customer) error
	Delete(id int64) error
	List(page, pageSize int) ([]*model.Customer, int64, error)
}

// RoleRepository 角色仓储接口
type RoleRepository interface {
	GetByID(id int64) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	Create(role *model.Role) error
	Update(role *model.Role) error
	Delete(id int64) error
	List(page, pageSize int) ([]*model.Role, int64, error)
	GetRolePermissions(roleID int64) ([]*model.Permission, error)
}

// PermissionRepository 权限仓储接口
type PermissionRepository interface {
	GetByID(id int64) (*model.Permission, error)
	GetByCode(code string) (*model.Permission, error)
	Create(permission *model.Permission) error
	Update(permission *model.Permission) error
	Delete(id int64) error
	List(page, pageSize int) ([]*model.Permission, int64, error)
	GetByModule(module string) ([]*model.Permission, error)
}