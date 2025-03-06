package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/image-recognition-engine/internal/repository"
)

type AuthHandler struct {
	adminRepo  repository.AdminUserRepository
	customerRepo repository.CustomerRepository
	jwtSecret  []byte
}

func NewAuthHandler(adminRepo repository.AdminUserRepository, customerRepo repository.CustomerRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		adminRepo:  adminRepo,
		customerRepo: customerRepo,
		jwtSecret:  []byte(jwtSecret),
	}
}

// AdminLogin 管理员登录
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var loginReq struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	// 查找用户
	user, err := h.adminRepo.GetByUsername(loginReq.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	// 生成JWT令牌
	token, err := h.generateJWT(user.ID, user.Username, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成令牌失败"})
		return
	}

	// 更新最后登录时间
	user.LastLoginTime = time.Now()
	if err := h.adminRepo.Update(user); err != nil {
		// 仅记录错误，不影响登录
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "登录成功",
		"data": gin.H{
			"token": token,
			"user": user,
		},
	})
}

// CustomerLogin 客户登录
func (h *AuthHandler) CustomerLogin(c *gin.Context) {
	var loginReq struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的请求参数"})
		return
	}

	// 查找客户
	customer, err := h.customerRepo.GetByUsername(loginReq.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(loginReq.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	// 生成JWT令牌
	token, err := h.generateJWT(customer.ID, customer.Username, "customer")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成令牌失败"})
		return
	}

	// 更新最后登录时间
	customer.LastLoginTime = time.Now()
	if err := h.customerRepo.Update(customer); err != nil {
		// 仅记录错误，不影响登录
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "登录成功",
		"data": gin.H{
			"token": token,
			"customer": customer,
		},
	})
}

// generateJWT 生成JWT令牌
func (h *AuthHandler) generateJWT(userID int64, username, userType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"username":  username,
		"user_type": userType,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		return "", errors.New("生成令牌失败")
	}

	return tokenString, nil
}