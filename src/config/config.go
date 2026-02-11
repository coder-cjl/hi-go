package config

import "time"

// JWT 配置常量
const (
	//  访问令牌过期时间（秒）
	JWTAccessTokenDuration = 7200 // 2小时

	//  刷新令牌过期时间（秒）
	JWTRefreshTokenDuration = 604800 // 7天

	// 签名密钥（生产环境应使用环境变量）
	JWTSecretKey = "your-secret-key-change-it"

	//  JWT 签发者
	JWTIssuer = "hi-go-app"
)

// Redis 配置常量
const (
	//  Token 在 Redis 中的过期时间（秒）
	RedisTokenTTL = JWTAccessTokenDuration // 与 JWT 过期时间一致

	//  会话在 Redis 中的过期时间（秒）
	RedisSessionTTL = 86400 // 24小时
)

// 数据库配置常量
const (
	//  最大打开连接数
	DBMaxOpenConns = 100

	//  最大空闲连接数
	DBMaxIdleConns = 10

	//  连接最大存活时间
	DBConnMaxLifetime = time.Hour

	// 空闲连接最大存活时间
	DBConnMaxIdleTime = 10 * time.Minute
)

// 业务配置常量
const (
	//  默认分页大小
	DefaultPageSize = 20

	//  最大分页大小
	MaxPageSize = 100

	//  密码最小长度
	PasswordMinLength = 6

	//  用户名最小长度
	UsernameMinLength = 3
)

// 获取 JWT Access Token 过期时间（time.Duration）
func GetJWTAccessTokenDuration() time.Duration {
	return time.Duration(JWTAccessTokenDuration) * time.Second
}

// 获取 JWT Refresh Token 过期时间（time.Duration）
func GetJWTRefreshTokenDuration() time.Duration {
	return time.Duration(JWTRefreshTokenDuration) * time.Second
}

// 获取 Redis Token TTL（time.Duration）
func GetRedisTokenTTL() time.Duration {
	return time.Duration(RedisTokenTTL) * time.Second
}

// 获取 Redis Session TTL（time.Duration）
func GetRedisSessionTTL() time.Duration {
	return time.Duration(RedisSessionTTL) * time.Second
}
