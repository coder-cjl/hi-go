package config

import "time"

// GetJWTAccessTokenDuration 获取 JWT Access Token 过期时间
func GetJWTAccessTokenDuration() time.Duration {
	return time.Duration(Config.JWT.AccessTokenDuration) * time.Second
}

// GetJWTRefreshTokenDuration 获取 JWT Refresh Token 过期时间
func GetJWTRefreshTokenDuration() time.Duration {
	return time.Duration(Config.JWT.RefreshTokenDuration) * time.Second
}

// GetRedisTokenTTL 获取 Redis Token TTL
func GetRedisTokenTTL() time.Duration {
	return time.Duration(Config.Redis.TokenTTL) * time.Second
}

// GetRedisSessionTTL 获取 Redis Session TTL
func GetRedisSessionTTL() time.Duration {
	return time.Duration(Config.Redis.SessionTTL) * time.Second
}

// GetDBConnMaxLifetime 获取数据库连接最大生命周期
func GetDBConnMaxLifetime() time.Duration {
	return time.Duration(Config.Database.ConnMaxLifetime) * time.Second
}

// GetDBConnMaxIdleTime 获取数据库连接最大空闲时间
func GetDBConnMaxIdleTime() time.Duration {
	return time.Duration(Config.Database.ConnMaxIdleTime) * time.Second
}
