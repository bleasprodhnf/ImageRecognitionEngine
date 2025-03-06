package mysql

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/image-recognition-engine/internal/model"
)

type servicePlanRepository struct {
	db *sqlx.DB
}

// NewServicePlanRepository 创建服务套餐数据访问实例
func NewServicePlanRepository(db *sqlx.DB) model.ServicePlanRepository {
	return &servicePlanRepository{db: db}
}

// Create 创建服务套餐
func (r *servicePlanRepository) Create(plan *model.ServicePlan) (int64, error) {
	plan.CreateTime = time.Now()
	plan.UpdateTime = time.Now()

	query := `INSERT INTO service_plans (name, description, price, request_limit, concurrent, features, create_time, update_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		plan.Name,
		plan.Description,
		plan.Price,
		plan.RequestLimit,
		plan.Concurrent,
		plan.Features,
		plan.CreateTime,
		plan.UpdateTime,
	)

	if err != nil {
		return 0, errors.Wrap(err, "创建服务套餐失败")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "获取套餐ID失败")
	}

	return id, nil
}

// Update 更新服务套餐
func (r *servicePlanRepository) Update(plan *model.ServicePlan) error {
	plan.UpdateTime = time.Now()

	query := `UPDATE service_plans SET 
		name = ?, 
		description = ?, 
		price = ?, 
		request_limit = ?, 
		concurrent = ?, 
		features = ?, 
		update_time = ? 
		WHERE id = ?`

	_, err := r.db.Exec(query,
		plan.Name,
		plan.Description,
		plan.Price,
		plan.RequestLimit,
		plan.Concurrent,
		plan.Features,
		plan.UpdateTime,
		plan.ID,
	)

	if err != nil {
		return errors.Wrap(err, "更新服务套餐失败")
	}

	return nil
}

// Delete 删除服务套餐
func (r *servicePlanRepository) Delete(id int64) error {
	query := `DELETE FROM service_plans WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return errors.Wrap(err, "删除服务套餐失败")
	}
	return nil
}

// List 获取服务套餐列表
func (r *servicePlanRepository) List(page, size int) ([]*model.ServicePlan, int64, error) {
	var total int64
	if err := r.db.Get(&total, "SELECT COUNT(*) FROM service_plans"); err != nil {
		return nil, 0, errors.Wrap(err, "获取服务套餐总数失败")
	}

	offset := (page - 1) * size
	query := `SELECT * FROM service_plans ORDER BY create_time DESC LIMIT ? OFFSET ?`
	plans := make([]*model.ServicePlan, 0)
	if err := r.db.Select(&plans, query, size, offset); err != nil {
		return nil, 0, errors.Wrap(err, "获取服务套餐列表失败")
	}

	return plans, total, nil
}

// FindByID 根据ID查找服务套餐
func (r *servicePlanRepository) FindByID(id int64) (*model.ServicePlan, error) {
	var plan model.ServicePlan
	query := `SELECT * FROM service_plans WHERE id = ?`
	err := r.db.Get(&plan, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "查询服务套餐失败")
	}
	return &plan, nil
}