package mysql

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/image-recognition-engine/internal/model"
)

type customerRepository struct {
	db *sqlx.DB
}

// NewCustomerRepository 创建客户数据访问实例
func NewCustomerRepository(db *sqlx.DB) model.CustomerRepository {
	return &customerRepository{db: db}
}

// FindByUsername 根据用户名查找客户
func (r *customerRepository) FindByUsername(username string) (*model.Customer, error) {
	var customer model.Customer
	query := `SELECT * FROM customers WHERE username = ?`
	err := r.db.Get(&customer, query, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "查询客户信息失败")
	}
	return &customer, nil
}

// Create 创建客户
func (r *customerRepository) Create(customer *model.Customer) (int64, error) {
	customer.CreateTime = time.Now()
	customer.UpdateTime = time.Now()

	query := `INSERT INTO customers (username, password, company, email, phone, plan_id, status, expire_time, last_login, create_time, update_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		customer.Username,
		customer.Password,
		customer.Company,
		customer.Email,
		customer.Phone,
		customer.PlanID,
		customer.Status,
		customer.ExpireTime,
		customer.LastLogin,
		customer.CreateTime,
		customer.UpdateTime,
	)

	if err != nil {
		return 0, errors.Wrap(err, "创建客户失败")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "获取客户ID失败")
	}

	return id, nil
}

// Update 更新客户信息
func (r *customerRepository) Update(customer *model.Customer) error {
	customer.UpdateTime = time.Now()

	query := `UPDATE customers SET 
		company = ?, 
		email = ?, 
		phone = ?, 
		plan_id = ?, 
		status = ?, 
		expire_time = ?, 
		update_time = ? 
		WHERE id = ?`

	_, err := r.db.Exec(query,
		customer.Company,
		customer.Email,
		customer.Phone,
		customer.PlanID,
		customer.Status,
		customer.ExpireTime,
		customer.UpdateTime,
		customer.ID,
	)

	if err != nil {
		return errors.Wrap(err, "更新客户信息失败")
	}

	return nil
}

// Delete 删除客户
func (r *customerRepository) Delete(id int64) error {
	query := `DELETE FROM customers WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return errors.Wrap(err, "删除客户失败")
	}
	return nil
}

// List 获取客户列表
func (r *customerRepository) List(page, size int) ([]*model.Customer, int64, error) {
	var total int64
	if err := r.db.Get(&total, "SELECT COUNT(*) FROM customers"); err != nil {
		return nil, 0, errors.Wrap(err, "获取客户总数失败")
	}

	offset := (page - 1) * size
	query := `SELECT * FROM customers ORDER BY create_time DESC LIMIT ? OFFSET ?`
	customers := make([]*model.Customer, 0)
	if err := r.db.Select(&customers, query, size, offset); err != nil {
		return nil, 0, errors.Wrap(err, "获取客户列表失败")
	}

	return customers, total, nil
}

// FindByID 根据ID查找客户
func (r *customerRepository) FindByID(id int64) (*model.Customer, error) {
	var customer model.Customer
	query := `SELECT * FROM customers WHERE id = ?`
	err := r.db.Get(&customer, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "查询客户信息失败")
	}
	return &customer, nil
}

// UpdateLastLogin 更新客户登录时间
func (r *customerRepository) UpdateLastLogin(id int64) error {
	query := `UPDATE customers SET last_login = ? WHERE id = ?`
	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return errors.Wrap(err, "更新登录时间失败")
	}
	return nil
}

// UpdateUsage 更新客户使用统计
func (r *customerRepository) UpdateUsage(usage *model.CustomerUsage) error {
	usage.UpdateTime = time.Now()

	query := `UPDATE customer_usage SET 
		request_count = ?, 
		success_count = ?, 
		fail_count = ?, 
		avg_latency = ?, 
		update_time = ? 
		WHERE customer_id = ? AND month = ?`

	result, err := r.db.Exec(query,
		usage.RequestCount,
		usage.SuccessCount,
		usage.FailCount,
		usage.AvgLatency,
		usage.UpdateTime,
		usage.CustomerID,
		usage.Month,
	)

	if err != nil {
		return errors.Wrap(err, "更新使用统计失败")
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "获取影响行数失败")
	}

	// 如果没有更新任何记录，说明需要插入新记录
	if affected == 0 {
		usage.CreateTime = time.Now()
		query = `INSERT INTO customer_usage (
			customer_id, request_count, success_count, fail_count, avg_latency, month, create_time, update_time
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

		_, err = r.db.Exec(query,
			usage.CustomerID,
			usage.RequestCount,
			usage.SuccessCount,
			usage.FailCount,
			usage.AvgLatency,
			usage.Month,
			usage.CreateTime,
			usage.UpdateTime,
		)

		if err != nil {
			return errors.Wrap(err, "创建使用统计失败")
		}
	}

	return nil
}