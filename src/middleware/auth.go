package middleware

import (
	"hi-go/src/model"
	"hi-go/src/utils/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			model.Unauthorized(c, "请求头中缺少 Authorization")
			c.Abort()
			return
		}

		// 2. 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			model.Unauthorized(c, "Authorization 格式错误，应为: Bearer {token}")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. 验证 token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			model.Unauthorized(c, "Token 验证失败: "+err.Error())
			c.Abort()
			return
		}

		// 4. 将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)

		// 5. 继续处理请求
		c.Next()
	}
}

// RoleAuth 角色权限中间件
func RoleAuth(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		roles, exists := c.Get("roles")
		if !exists {
			model.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		userRoles := roles.([]string)

		// 检查是否有权限
		hasPermission := false
		for _, userRole := range userRoles {
			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			model.Forbidden(c, "权限不足，需要角色: "+strings.Join(allowedRoles, ", "))
			c.Abort()
			return
		}

		c.Next()
	}
}
