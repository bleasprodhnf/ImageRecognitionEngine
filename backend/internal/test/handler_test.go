package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/image-recognition-engine/internal/handler/common"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建路由
	r := gin.New()
	r.GET("/health", common.HealthCheck)

	// 创建请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	// 执行请求
	r.ServeHTTP(w, req)

	// 检查状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "服务正常", response["message"])

	// 验证data字段
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "ok", data["status"])
	assert.NotEmpty(t, data["timestamp"])
	assert.Equal(t, "1.0.0", data["version"])

	// 验证services字段
	services, ok := data["services"].(map[string]interface{})
	assert.True(t, ok)
	assert.True(t, services["database"].(bool))
	assert.True(t, services["redis"].(bool))
	assert.True(t, services["storage"].(bool))
}