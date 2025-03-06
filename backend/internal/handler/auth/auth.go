package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/image-recognition-engine/config"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
	UserInfo  struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		RealName string `json:"realName"`
		Email    string `json:"email"`
	} `json:"userInfo"`
}

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	RealName string `json:"realName"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// Login 处理登录请求
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// TODO: 从数据库验证用户名和密码
	// 这里简化处理，实际应用中需要查询数据库并验证密码
	if req.Username != "admin" || req.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
			"data":    nil,
		})
		return
	}

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务器内部错误",
			"data":    nil,
		})
		return
	}

	// 创建JWT令牌
	expiration := time.Duration(cfg.JWT.Expiration) * time.Hour
	expiresAt := time.Now().Add(expiration).Unix()

	claims := JWTClaims{
		UserID:   1, // 示例ID，实际应从数据库获取
		Username: req.Username,
		RealName: "管理员", // 示例，实际应从数据库获取
		Email:    "admin@example.com", // 示例，实际应从数据库获取
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresAt, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "image-recognition-engine",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成令牌失败",
			"data":    nil,
		})
		return
	}

	// 构建响应
	response := LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}
	response.UserInfo.ID = 1
	response.UserInfo.Username = req.Username
	response.UserInfo.RealName = "管理员"
	response.UserInfo.Email = "admin@example.com"

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    response,
	})
}

// Register 处理注册请求
func Register(c *gin.Context) {
	// TODO: 实现注册逻辑
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data":    nil,
	})
}