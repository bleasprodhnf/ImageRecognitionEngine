package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/image-recognition-engine/internal/model"
)

type CustomerHandler struct {
	customerRepo model.CustomerRepository
	planRepo     model.ServicePlanRepository
}

// NewCustomerHandler 创建客户处理器实例
func NewCustomerHandler(customerRepo model.CustomerRepository, planRepo model.ServicePlanRepository) *CustomerHandler {
	return &CustomerHandler{
		customerRepo: customerRepo,
		planRepo:     planRepo,
	}
}

// Create 创建客户
func (h *CustomerHandler) Create(c *gin.Context) {
	var customer model.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	// 检查服务套餐是否存在
	if customer.PlanID > 0 {
		plan, err := h.planRepo.FindByID(customer.PlanID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "系统错误"})
			return
		}
		if plan == nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "服务套餐不存在"})
			return
		}
	}

	// 检查用户名是否已存在
	exist, err := h.customerRepo.FindByUsername(customer.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "系统错误"})
		return
	}
	if exist != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名已存在"})
		return
	}

	id, err := h.customerRepo.Create(&customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建客户失败"})
		return
	}

	customer.ID = id
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    customer,
	})
}

// Update 更新客户信息
func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的客户ID"})
		return
	}

	customer, err := h.customerRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "系统错误"})
		return
	}
	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "客户不存在"})
		return
	}

	if err := c.ShouldBindJSON(customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	// 检查服务套餐是否存在
	if customer.PlanID > 0 {
		plan, err := h.planRepo.FindByID(customer.PlanID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "系统错误"})
			return
		}
		if plan == nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "服务套餐不存在"})
			return
		}
	}

	if err := h.customerRepo.Update(customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新客户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    customer,
	})
}

// Delete 删除客户
func (h *CustomerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的客户ID"})
		return
	}

	customer, err := h.customerRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "系统错误"})
		return
	}
	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "客户不存在"})
		return
	}

	if err := h.customerRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除客户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// List 获取客户列表
func (h *CustomerHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	customers, total, err := h.customerRepo.List(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取客户列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": size,
			"list":     customers,
		},
	})
}

// GetByID 根据ID获取客户信息
func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的客户ID"})
		return
	}

	customer, err := h.customerRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "系统错误"})
		return
	}
	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "客户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data":    customer,
	})
}

// UpdateUsage 更新客户使用统计
func (h *CustomerHandler) UpdateUsage(c *gin.Context) {
	var usage model.CustomerUsage
	if err := c.ShouldBindJSON(&usage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	// 检查客户是否存在
	customer, err := h.customerRepo.FindByID(usage.CustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "系统错误"})
		return
	}
	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "客户不存在"})
		return
	}

	if err := h.customerRepo.UpdateUsage(&usage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新使用统计失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    usage,
	})
}