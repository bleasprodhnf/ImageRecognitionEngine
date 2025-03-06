package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/image-recognition-engine/internal/database"
	"github.com/image-recognition-engine/internal/model"
)

// UserRepositoryImpl 用户数据访问实现
type UserRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository 创建用户数据访问实例
func NewUserRepository() model.UserRepository {
	return &UserRepositoryImpl{
		db: database.MySQLDB,
	}
}

// FindByUsername 根据用户名查找用户
func (r *UserRepositoryImpl) FindByUsername(username string) (*model.User, error) {
	query := `SELECT id, username, password, real_name, email, phone, avatar, role_id, status, 
		last_login, create_time, update_time FROM users WHERE username = ?`

	var user model.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.RealName, &user.Email, 
		&user.Phone, &user.Avatar, &user.RoleID, &user.Status, &user.LastLogin, 
		&user.CreateTime, &user.UpdateTime,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 用户不存在
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

// Create 创建用户
func (r *UserRepositoryImpl) Create(user *model.User) (int64, error) {
	query := `INSERT INTO users (username, password, real_name, email, phone, avatar, role_id, status, 
		last_login, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	user.CreateTime = now
	user.UpdateTime = now

	result, err := r.db.Exec(query, 
		user.Username, user.Password, user.RealName, user.Email, user.Phone, 
		user.Avatar, user.RoleID, user.Status, user.LastLogin, user.CreateTime, user.UpdateTime,
	)

	if err != nil {
		return 0, fmt.Errorf("创建用户失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取用户ID失败: %w", err)
	}

	return id, nil
}

// Update 更新用户
func (r *UserRepositoryImpl) Update(user *model.User) error {
	query := `UPDATE users SET real_name = ?, email = ?, phone = ?, avatar = ?, role_id = ?, 
		status = ?, update_time = ? WHERE id = ?`

	user.UpdateTime = time.Now()

	_, err := r.db.Exec(query, 
		user.RealName, user.Email, user.Phone, user.Avatar, 
		user.RoleID, user.Status, user.UpdateTime, user.ID,
	)

	if err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}

	return nil
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	return nil
}

// List 获取用户列表
func (r *UserRepositoryImpl) List(page, size int) ([]*model.User, int64, error) {
	// 计算总数
	countQuery := `SELECT COUNT(*) FROM users`
	var total int64
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * size
	query := `SELECT id, username, real_name, email, phone, avatar, role_id, status, 
		last_login, create_time, update_time FROM users ORDER BY id DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, size, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %w", err)
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.RealName, &user.Email, 
			&user.Phone, &user.Avatar, &user.RoleID, &user.Status, 
			&user.LastLogin, &user.CreateTime, &user.UpdateTime,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("扫描用户数据失败: %w", err)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("遍历用户数据失败: %w", err)
	}

	return users, total, nil
}

// FindByID 根据ID查找用户
func (r *UserRepositoryImpl) FindByID(id int64) (*model.User, error) {
	query := `SELECT id, username, real_name, email, phone, avatar, role_id, status, 
		last_login, create_time, update_time FROM users WHERE id = ?`

	var user model.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.RealName, &user.Email, 
		&user.Phone, &user.Avatar, &user.RoleID, &user.Status, 
		&user.LastLogin, &user.CreateTime, &user.UpdateTime,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 用户不存在
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

// UpdateLastLogin 更新用户登录时间
func (r *UserRepositoryImpl) UpdateLastLogin(id int64) error {
	query := `UPDATE users SET last_login = ? WHERE id = ?`

	now := time.Now()
	_, err := r.db.Exec(query, now, id)
	if err != nil {
		return fmt.Errorf("更新用户登录时间失败: %w", err)
	}

	return nil
}