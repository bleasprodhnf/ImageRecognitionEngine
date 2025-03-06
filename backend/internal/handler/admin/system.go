package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SystemParam 系统参数结构
type SystemParam struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	Description string `json:"description"`
	UpdateTime  string `json:"updateTime"`
}

// GetSystemParams 获取系统参数列表
func GetSystemParams(c *gin.Context) {
	// TODO: 从数据库获取系统参数列表
	// 这里返回示例数据
	params := []SystemParam{
		{
			ID:          1,
			Name:        "maxUploadSize",
			DisplayName: "最大上传大小",
			Value:       "10",
			Type:        "number",
			Description: "单个文件最大上传大小(MB)",
			UpdateTime:  "2024-01-25T00:00:00Z",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    params,
	})
}

// CreateSystemParam 创建系统参数
func CreateSystemParam(c *gin.Context) {
	var param SystemParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// TODO: 保存系统参数到数据库
	// 这里返回示例数据
	param.ID = 1
	param.UpdateTime = "2024-01-25T00:00:00Z"

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    param,
	})
}

// UpdateSystemParam 更新系统参数
func UpdateSystemParam(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数ID无效",
			"data":    nil,
		})
		return
	}

	var param SystemParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// TODO: 更新数据库中的系统参数
	// 这里返回示例数据
	param.ID = id
	param.UpdateTime = "2024-01-25T00:00:00Z"

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    param,
	})
}

// DeleteSystemParam 删除系统参数
func DeleteSystemParam(c *gin.Context) {
	_, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数ID无效",
			"data":    nil,
		})
		return
	}

	// TODO: 从数据库删除系统参数

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data":    nil,
	})
}