package jwt

// import (
// 	"fmt"
// 	"time"
// )

// // ExampleGlobal 展示使用全局JWT管理器的方法（推荐方式）
// func ExampleGlobal() {
// 	fmt.Println("=== 使用全局JWT管理器（推荐方式）===\n")

// 	// ========== 1. 使用默认配置（已自动初始化） ==========
// 	// 包导入时会自动调用 init() 函数初始化全局Manager
// 	fmt.Println("✓ 全局JWT管理器已自动初始化\n")

// 	// ========== 2. 或者使用自定义配置重新初始化 ==========
// 	customConfig := &Config{
// 		SecretKey:            "my-super-secret-key-123456",
// 		AccessTokenDuration:  30 * time.Minute,
// 		RefreshTokenDuration: 24 * time.Hour,
// 		Issuer:               "my-app",
// 	}
// 	Init(customConfig)
// 	fmt.Println("✓ 使用自定义配置重新初始化成功\n")

// 	// ========== 3. 直接使用全局函数生成令牌 ==========
// 	userID := "user123"
// 	username := "张三"
// 	roles := []string{"admin", "user"}
// 	extra := map[string]interface{}{
// 		"department": "技术部",
// 		"level":      5,
// 	}

// 	// 直接调用包级别的函数，无需 Manager.XXX
// 	accessToken, refreshToken, err := GenerateTokenPair(userID, username, roles, extra)
// 	if err != nil {
// 		fmt.Printf("生成令牌对失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("✓ 生成令牌对成功:\n访问令牌: %s...\n刷新令牌: %s...\n\n",
// 		accessToken[:50], refreshToken[:50])

// 	// ========== 4. 解析令牌 ==========
// 	claims, err := ParseToken(accessToken)
// 	if err != nil {
// 		fmt.Printf("解析令牌失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("✓ 解析令牌成功:\n")
// 	fmt.Printf("  用户ID: %s\n", claims.UserID)
// 	fmt.Printf("  用户名: %s\n", claims.Username)
// 	fmt.Printf("  角色: %v\n", claims.Roles)
// 	fmt.Printf("  部门: %v\n\n", claims.Extra["department"])

// 	// ========== 5. 验证令牌 ==========
// 	valid, err := ValidateToken(accessToken)
// 	if err != nil {
// 		fmt.Printf("验证令牌失败: %v\n", err)
// 	} else {
// 		fmt.Printf("✓ 令牌验证结果: %v\n\n", valid)
// 	}

// 	// ========== 6. 检查角色 ==========
// 	hasAdmin, _ := HasRole(accessToken, "admin")
// 	fmt.Printf("✓ 是否包含admin角色: %v\n\n", hasAdmin)

// 	// ========== 7. 刷新令牌 ==========
// 	newAccessToken, err := RefreshAccessToken(refreshToken, username, roles, extra)
// 	if err != nil {
// 		fmt.Printf("刷新令牌失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("✓ 刷新访问令牌成功:\n%s...\n\n", newAccessToken[:50])

// 	// ========== 8. 获取全局管理器实例（如需要） ==========
// 	manager := GetManager()
// 	fmt.Printf("✓ 全局管理器配置:\n")
// 	fmt.Printf("  签发者: %s\n", manager.config.Issuer)
// 	fmt.Printf("  访问令牌过期时间: %v\n", manager.config.AccessTokenDuration)
// 	fmt.Printf("  刷新令牌过期时间: %v\n", manager.config.RefreshTokenDuration)

// 	fmt.Println("\n=== 全局方式示例结束 ===")
// }

// // Example 展示JWT工具的基本使用方法（创建独立实例）
// func Example() {
// 	// ========== 1. 创建JWT管理器 ==========
// 	// 使用默认配置
// 	jwtManager := NewJWTManager(nil)

// 	// 或者使用自定义配置
// 	customConfig := &Config{
// 		SecretKey:            "my-super-secret-key-123456",
// 		AccessTokenDuration:  30 * time.Minute,
// 		RefreshTokenDuration: 24 * time.Hour,
// 		Issuer:               "my-app",
// 	}
// 	jwtManager = NewJWTManager(customConfig)

// 	fmt.Println("=== JWT工具使用示例（独立实例方式）===\n")

// 	// ========== 2. 生成访问令牌 ==========
// 	userID := "user123"
// 	username := "张三"
// 	roles := []string{"admin", "user"}
// 	extra := map[string]interface{}{
// 		"department": "技术部",
// 		"level":      5,
// 	}

// 	accessToken, err := jwtManager.GenerateToken(userID, username, roles, extra)
// 	if err != nil {
// 		fmt.Printf("生成访问令牌失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("✓ 生成访问令牌成功:\n%s\n\n", accessToken)

// 	// ========== 3. 生成刷新令牌 ==========
// 	refreshToken, err := jwtManager.GenerateRefreshToken(userID)
// 	if err != nil {
// 		fmt.Printf("生成刷新令牌失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("✓ 生成刷新令牌成功:\n%s\n\n", refreshToken)

// 	// ========== 4. 同时生成访问令牌和刷新令牌 ==========
// 	accessToken2, refreshToken2, err := jwtManager.GenerateTokenPair(userID, username, roles, extra)
// 	if err != nil {
// 		fmt.Printf("生成令牌对失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("✓ 生成令牌对成功:\n访问令牌: %s...\n刷新令牌: %s...\n\n",
// 		accessToken2[:50], refreshToken2[:50])

// 	// ========== 5. 解析和验证令牌 ==========
// 	claims, err := jwtManager.ParseToken(accessToken)
// 	if err != nil {
// 		fmt.Printf("解析令牌失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("✓ 解析令牌成功:\n")
// 	fmt.Printf("  用户ID: %s\n", claims.UserID)
// 	fmt.Printf("  用户名: %s\n", claims.Username)
// 	fmt.Printf("  角色: %v\n", claims.Roles)
// 	fmt.Printf("  部门: %v\n", claims.Extra["department"])
// 	fmt.Printf("  级别: %v\n", claims.Extra["level"])
// 	fmt.Printf("  签发时间: %s\n", claims.IssuedAt.Time.Format("2006-01-02 15:04:05"))
// 	fmt.Printf("  过期时间: %s\n\n", claims.ExpiresAt.Time.Format("2006-01-02 15:04:05"))

// 	// ========== 6. 验证令牌是否有效 ==========
// 	valid, err := jwtManager.ValidateToken(accessToken)
// 	if err != nil {
// 		fmt.Printf("验证令牌失败: %v\n", err)
// 	} else {
// 		fmt.Printf("✓ 令牌验证结果: %v\n\n", valid)
// 	}

// 	// ========== 7. 从令牌中提取用户ID ==========
// 	extractedUserID, err := jwtManager.GetUserIDFromToken(accessToken)
// 	if err != nil {
// 		fmt.Printf("提取用户ID失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("✓ 提取用户ID: %s\n\n", extractedUserID)

// 	// ========== 8. 检查令牌是否过期 ==========
// 	isExpired := jwtManager.IsTokenExpired(accessToken)
// 	fmt.Printf("✓ 令牌是否过期: %v\n\n", isExpired)

// 	// ========== 9. 检查是否包含指定角色 ==========
// 	hasAdmin, err := jwtManager.HasRole(accessToken, "admin")
// 	if err != nil {
// 		fmt.Printf("检查角色失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("✓ 是否包含admin角色: %v\n\n", hasAdmin)

// 	hasGuest, err := jwtManager.HasRole(accessToken, "guest")
// 	if err != nil {
// 		fmt.Printf("检查角色失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("✓ 是否包含guest角色: %v\n\n", hasGuest)

// 	// ========== 10. 刷新访问令牌 ==========
// 	newAccessToken, err := jwtManager.RefreshAccessToken(refreshToken, username, roles, extra)
// 	if err != nil {
// 		fmt.Printf("刷新访问令牌失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("✓ 刷新访问令牌成功:\n%s...\n\n", newAccessToken[:50])

// 	// ========== 11. 演示过期令牌 ==========
// 	// 创建一个已过期的令牌（过期时间设为1纳秒）
// 	expiredJWT := NewJWTManager(&Config{
// 		SecretKey:           "test-key",
// 		AccessTokenDuration: 1 * time.Nanosecond,
// 		Issuer:              "test",
// 	})
// 	expiredToken, _ := expiredJWT.GenerateToken(userID, username, roles, nil)

// 	// 等待令牌过期
// 	time.Sleep(10 * time.Millisecond)

// 	_, err = expiredJWT.ParseToken(expiredToken)
// 	if err != nil {
// 		fmt.Printf("✓ 过期令牌验证失败（预期行为）: %v\n\n", err)
// 	}

// 	// ========== 12. 修改配置 ==========
// 	jwtManager.SetSecretKey("new-secret-key")
// 	jwtManager.SetAccessTokenDuration(1 * time.Hour)
// 	jwtManager.SetRefreshTokenDuration(30 * 24 * time.Hour)
// 	fmt.Println("✓ 更新JWT配置成功")

// 	fmt.Println("\n=== 示例结束 ===")
// }

// // ExampleGinMiddleware 在Gin框架中使用JWT的示例
// func ExampleGinMiddleware() {
// 	fmt.Println(`
// === Gin中间件使用示例 ===

// // 创建JWT管理器
// var jwtManager = jwt.NewJWTManager(&jwt.Config{
//     SecretKey:            "your-secret-key",
//     AccessTokenDuration:  15 * time.Minute,
//     RefreshTokenDuration: 7 * 24 * time.Hour,
//     Issuer:               "your-app",
// })

// // JWT中间件
// func JWTAuthMiddleware() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         // 从请求头获取token
//         tokenString := c.GetHeader("Authorization")
//         if tokenString == "" {
//             c.JSON(401, gin.H{"error": "未提供token"})
//             c.Abort()
//             return
//         }

//         // 移除 "Bearer " 前缀
//         if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
//             tokenString = tokenString[7:]
//         }

//         // 解析token
//         claims, err := jwtManager.ParseToken(tokenString)
//         if err != nil {
//             c.JSON(401, gin.H{"error": "无效的token", "detail": err.Error()})
//             c.Abort()
//             return
//         }

//         // 将用户信息存入上下文
//         c.Set("userID", claims.UserID)
//         c.Set("username", claims.Username)
//         c.Set("roles", claims.Roles)
//         c.Next()
//     }
// }

// // 角色检查中间件
// func RequireRole(role string) gin.HandlerFunc {
//     return func(c *gin.Context) {
//         roles, exists := c.Get("roles")
//         if !exists {
//             c.JSON(403, gin.H{"error": "未找到角色信息"})
//             c.Abort()
//             return
//         }

//         roleList, ok := roles.([]string)
//         if !ok {
//             c.JSON(500, gin.H{"error": "角色信息格式错误"})
//             c.Abort()
//             return
//         }

//         hasRole := false
//         for _, r := range roleList {
//             if r == role {
//                 hasRole = true
//                 break
//             }
//         }

//         if !hasRole {
//             c.JSON(403, gin.H{"error": "权限不足"})
//             c.Abort()
//             return
//         }

//         c.Next()
//     }
// }

// // 登录处理
// func Login(c *gin.Context) {
//     var loginReq struct {
//         Username string ` + "`json:\"username\"`" + `
//         Password string ` + "`json:\"password\"`" + `
//     }

//     if err := c.ShouldBindJSON(&loginReq); err != nil {
//         c.JSON(400, gin.H{"error": "参数错误"})
//         return
//     }

//     // TODO: 验证用户名和密码
//     userID := "user123"
//     username := loginReq.Username
//     roles := []string{"user"}

//     // 生成令牌对
//     accessToken, refreshToken, err := jwtManager.GenerateTokenPair(
//         userID, username, roles, nil,
//     )
//     if err != nil {
//         c.JSON(500, gin.H{"error": "生成token失败"})
//         return
//     }

//     c.JSON(200, gin.H{
//         "access_token":  accessToken,
//         "refresh_token": refreshToken,
//         "expires_in":    900, // 15分钟
//     })
// }

// // 刷新token
// func RefreshToken(c *gin.Context) {
//     var req struct {
//         RefreshToken string ` + "`json:\"refresh_token\"`" + `
//     }

//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(400, gin.H{"error": "参数错误"})
//         return
//     }

//     // TODO: 从数据库获取用户信息
//     username := "张三"
//     roles := []string{"user"}
//     extra := map[string]interface{}{}

//     // 刷新访问令牌
//     newAccessToken, err := jwtManager.RefreshAccessToken(
//         req.RefreshToken, username, roles, extra,
//     )
//     if err != nil {
//         c.JSON(401, gin.H{"error": "刷新token失败", "detail": err.Error()})
//         return
//     }

//     c.JSON(200, gin.H{
//         "access_token": newAccessToken,
//         "expires_in":   900,
//     })
// }

// // 路由设置
// r := gin.Default()

// // 公开路由
// r.POST("/login", Login)
// r.POST("/refresh", RefreshToken)

// // 需要认证的路由
// auth := r.Group("/api")
// auth.Use(JWTAuthMiddleware())
// {
//     auth.GET("/profile", func(c *gin.Context) {
//         userID, _ := c.Get("userID")
//         username, _ := c.Get("username")
//         c.JSON(200, gin.H{
//             "user_id":  userID,
//             "username": username,
//         })
//     })

//     // 需要admin角色的路由
//     admin := auth.Group("/admin")
//     admin.Use(RequireRole("admin"))
//     {
//         admin.GET("/users", func(c *gin.Context) {
//             c.JSON(200, gin.H{"users": []string{"user1", "user2"}})
//         })
//     }
// }
// `)
// }
