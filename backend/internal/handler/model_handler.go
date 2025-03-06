package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/image-recognition-engine/internal/model"
	"github.com/image-recognition-engine/internal/repository"
)

type ModelHandler struct {
	modelRepo *repository.ModelRepository
}

func NewModelHandler(modelRepo *repository.ModelRepository) *ModelHandler {
	return &ModelHandler{modelRepo: modelRepo}
}

// CreateModelVersion 创建新的模型版本
func (h *ModelHandler) CreateModelVersion(c *gin.Context) {
	var version model.ModelVersion
	if err := c.ShouldBindJSON(&version); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := h.modelRepo.CreateModelVersion(c.Request.Context(), &version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建模型版本失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": version})
}

// GetModelVersions 获取模型版本列表
func (h *ModelHandler) GetModelVersions(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("pageSize", "20"), 10, 64)
	status := c.Query("status")

	versions, total, err := h.modelRepo.GetModelVersions(c.Request.Context(), status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取模型版本列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"total": total,
			"items": versions,
		},
	})
}

// SaveModelPerformance 保存模型性能数据
func (h *ModelHandler) SaveModelPerformance(c *gin.Context) {
	var perf model.ModelPerformance
	if err := c.ShouldBindJSON(&perf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := h.modelRepo.SaveModelPerformance(c.Request.Context(), &perf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存模型性能数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "保存成功", "data": perf})
}

// GetModelPerformance 获取模型性能数据
func (h *ModelHandler) GetModelPerformance(c *gin.Context) {
	modelIDStr := c.Query("modelId")
	modelID, err := primitive.ObjectIDFromHex(modelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "模型ID格式错误"})
		return
	}

	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")

	var startTime, endTime time.Time
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "开始时间格式错误"})
			return
		}
	}
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "结束时间格式错误"})
			return
		}
	}

	perf, err := h.modelRepo.GetModelPerformance(c.Request.Context(), modelID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取模型性能数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "成功", "data": perf})
}