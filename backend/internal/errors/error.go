package errors

import (
	"fmt"
	"net/http"
	"time"
)

// ErrorCode 定义错误码类型
type ErrorCode int

// 系统级错误码定义 (1000-1999)
const (
	InternalServerError ErrorCode = 1000 + iota
	DatabaseError
	CacheError
	ValidationError
	AuthenticationError
	AuthorizationError
	RateLimitError
	NotFoundError
	BadRequestError
	TimeoutError
)

// 业务级错误码定义 (2000-2999)
const (
	ModelNotFoundError ErrorCode = 2000 + iota
	InvalidModelParameterError
	ModelTrainingError
	ImageProcessingError
	RecognitionError
)

// AppError 定义应用错误结构
type AppError struct {
	Code      ErrorCode  `json:"code"`
	HTTPCode  int       `json:"httpCode"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details,omitempty"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
}

// NewAppError 创建新的应用错误
func NewAppError(code ErrorCode, message string, details string) *AppError {
	httpCode := getHTTPStatusCode(code)
	return &AppError{
		Code:      code,
		HTTPCode:  httpCode,
		Message:   message,
		Timestamp: time.Now(),
		Details:   details,
	}
}

// getHTTPStatusCode 根据错误码获取对应的HTTP状态码
func getHTTPStatusCode(code ErrorCode) int {
	switch code {
	case InternalServerError:
		return http.StatusInternalServerError
	case DatabaseError, CacheError:
		return http.StatusServiceUnavailable
	case ValidationError, BadRequestError:
		return http.StatusBadRequest
	case AuthenticationError:
		return http.StatusUnauthorized
	case AuthorizationError:
		return http.StatusForbidden
	case RateLimitError:
		return http.StatusTooManyRequests
	case NotFoundError, ModelNotFoundError:
		return http.StatusNotFound
	case TimeoutError:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}

// IsNotFound 检查是否为未找到错误
func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == NotFoundError || appErr.Code == ModelNotFoundError
	}
	return false
}

// IsBadRequest 检查是否为请求参数错误
func IsBadRequest(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ValidationError || appErr.Code == BadRequestError
	}
	return false
}

// IsAuthError 检查是否为认证/授权错误
func IsAuthError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == AuthenticationError || appErr.Code == AuthorizationError
	}
	return false
}