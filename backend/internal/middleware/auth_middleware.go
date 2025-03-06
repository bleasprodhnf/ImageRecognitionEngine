package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RequireAdmin 需要管理员权限的中间件
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息，检查是否为管理员
		// 这里应该从JWT或会话中获取用户信息
		// 为了示例，我们简单实现
		
		// TODO: 从JWT或会话中获取用户角色信息
		// 假设我们从Authorization头中获取的JWT已经在AuthMiddleware中验证
		// 并且用户信息已经存储在上下文中
		
		// 模拟检查用户角色
		isAdmin := false
		
		// 实际实现应该从JWT中提取用户信息并验证角色
		// userInfo, exists := c.Get("userInfo")
		// if !exists {
		//     isAdmin = false
		// } else {
		//     user := userInfo.(UserInfo)
		//     isAdmin = user.Role == "admin"
		// }
		
		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "需要管理员权限",
			})
			return
		}
		
		c.Next()
	}
}

// RequireAuth 需要认证的中间件
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查用户是否已认证
		// 这里应该从JWT或会话中获取用户信息
		
		// TODO: 实现实际的认证检查逻辑
		// 假设认证信息已经在AuthMiddleware中验证
		// 这里只是再次确认用户信息存在
		
		// 模拟检查用户是否已认证
		isAuthenticated := false
		
		// 实际实现应该检查JWT或会话
		// userInfo, exists := c.Get("userInfo")
		// isAuthenticated = exists
		
		if !isAuthenticated {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "需要登录",
			})
			return
		}
		
		c.Next()
	}
}

// RequirePermission 需要特定权限的中间件
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查用户是否有特定权限
		// 这里应该从JWT或会话中获取用户权限信息
		
		// TODO: 实现实际的权限检查逻辑
		// 假设用户信息和权限已经在AuthMiddleware中验证并存储在上下文中
		
		// 模拟检查用户权限
		hasPermission := false
		
		// 实际实现应该检查用户的权限列表
		// userInfo, exists := c.Get("userInfo")
		// if !exists {
		//     hasPermission = false
		// } else {
		//     user := userInfo.(UserInfo)
		//     hasPermission = user.HasPermission(permission)
		// }
		
		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "没有操作权限",
			})
			return
		}
		
		c.Next()
	}
}