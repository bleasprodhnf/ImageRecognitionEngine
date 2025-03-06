package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/image-recognition-engine/internal/errors"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ResponseMiddleware 响应处理中间件
func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理错误响应
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if appErr, ok := err.(*errors.AppError); ok {
				c.JSON(appErr.HTTPCode, Response{
					Code:    int(appErr.Code),
					Message: appErr.Message,
				})
				return
			}

			// 处理其他类型的错误
			c.JSON(http.StatusInternalServerError, Response{
				Code:    http.StatusInternalServerError,
				Message: "服务器内部错误",
			})
			return
		}

		// 处理成功响应
		if c.Writer.Written() {
			return
		}

		if data, exists := c.Get("response_data"); exists {
			c.JSON(http.StatusOK, Response{
				Code:    http.StatusOK,
				Message: "成功",
				Data:    data,
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Code:    http.StatusOK,
			Message: "成功",
		})
	}
}

// Success 设置成功响应数据
func Success(c *gin.Context, data interface{}) {
	c.Set("response_data", data)
}

// Error 设置错误响应
func Error(c *gin.Context, err error) {
	c.Error(err)
}