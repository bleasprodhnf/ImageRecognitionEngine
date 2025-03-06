package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/image-recognition-engine/internal/model"
)

type ServicePlanHandler struct {
	planRepo model.ServicePlanRepository
}

// NewServicePlanHandler 创建服务套餐处理器实例
func NewServicePlanHandler(planRepo model.ServicePlanRepository) *ServicePlanHandler {
	return &ServicePlanHandler{
		planRepo: planRepo,
	}
}

// Create 创建服务套餐
func (h *ServicePlanHandler) Create(c *gin.Context) {
	var plan model.ServicePlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	id, err := h.planRepo.Create(&plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建服务套餐失败"})
		return
	}

	plan.ID = id
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    plan,
	})
}

// Update 更新服务套餐
func (h *ServicePlanHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的套餐ID"})
		return
	}

	plan, err := h.planRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "系统错误"})
		return
	}
	if plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "服务套餐不存在"})
		return
	}

	if err := c.ShouldBindJSON(plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	if err := h.planRepo.Update(plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新服务套餐失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    plan,
	})
}

// Delete 删除服务套餐
func (h *ServicePlanHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的套餐ID"})
		return
	}

	plan, err := h.planRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "系统错误"})
		return
	}
	if plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "服务套餐不存在"})
		return
	}

	if err := h.planRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除服务套餐失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// List 获取服务套餐列表
func (h *ServicePlanHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	plans, total, err := h.planRepo.List(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取服务套餐列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": size,
			"list":     plans,
		},
	})
}

// GetByID 根据ID获取服务套餐信息
func (h *ServicePlanHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的套餐ID"})
		return
	}

	plan, err := h.planRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "系统错误"})
		return
	}
	if plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "服务套餐不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data":    plan,
	})
}