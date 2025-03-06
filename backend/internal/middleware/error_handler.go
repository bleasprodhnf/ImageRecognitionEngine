package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/image-recognition-engine/internal/errors"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic堆栈信息
				log.Printf("Panic: %v\n%s", err, debug.Stack())

				// 创建内部服务器错误
				appErr := errors.NewAppError(
					errors.InternalServerError,
					"服务器内部错误",
					"系统发生未预期的错误",
				)

				// 返回错误响应
				c.JSON(appErr.HTTPCode, gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				})
				c.Abort()
			}
		}()

		c.Next()

		// 处理请求过程中的错误
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				if appErr, ok := e.Err.(*errors.AppError); ok {
					// 记录错误日志
					log.Printf("Error: %s", appErr.Error())

					// 返回错误响应
					c.JSON(appErr.HTTPCode, gin.H{
						"code":    appErr.Code,
						"message": appErr.Message,
						"details": appErr.Details,
					})
					return
				}

				// 处理非AppError类型的错误
				log.Printf("Unexpected error: %s", e.Error())
				appErr := errors.NewAppError(
					errors.InternalServerError,
					"服务器内部错误",
					e.Error(),
				)
				c.JSON(appErr.HTTPCode, gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				})
				return
			}
		}
	}
}

// RecoveryHandler 恢复处理中间件
func RecoveryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic信息
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())

				// 返回500错误
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    errors.InternalServerError,
					"message": "服务器内部错误",
				})
			}
		}()
		c.Next()
	}
}