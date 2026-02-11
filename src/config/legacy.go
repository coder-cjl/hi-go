package config

import "time"

// 向后兼容的全局变量
// 这些变量用于兼容旧代码，新代码建议直接使用 Config 实例
var (
	JWTSecretKey            string
	JWTIssuer               string
	JWTAccessTokenDuration  int
	JWTRefreshTokenDuration int
	RedisTokenTTL           int
	RedisSessionTTL         int
	DBMaxOpenConns          int
	DBMaxIdleConns          int
	DBConnMaxLifetime       time.Duration
	DBConnMaxIdleTime       time.Duration
	DefaultPageSize         int
	MaxPageSize             int
	PasswordMinLength       int
	UsernameMinLength       int
	SnowflakeMachineID      int64
)

// UpdateLegacyVars 更新向后兼容的全局变量
// 在配置加载后调用，将配置值同步到全局变量中
// 注意：这是为了向后兼容，新代码应该直接使用 config.Config
func UpdateLegacyVars() {
	if Config == nil {
		return
	}
	JWTSecretKey = Config.JWT.SecretKey
	JWTIssuer = Config.JWT.Issuer
	JWTAccessTokenDuration = Config.JWT.AccessTokenDuration
	JWTRefreshTokenDuration = Config.JWT.RefreshTokenDuration
	RedisTokenTTL = Config.Redis.TokenTTL
	RedisSessionTTL = Config.Redis.SessionTTL
	DBMaxOpenConns = Config.Database.MaxOpenConns
	DBMaxIdleConns = Config.Database.MaxIdleConns
	DBConnMaxLifetime = GetDBConnMaxLifetime()
	DBConnMaxIdleTime = GetDBConnMaxIdleTime()
	DefaultPageSize = Config.Business.DefaultPageSize
	MaxPageSize = Config.Business.MaxPageSize
	PasswordMinLength = Config.Business.PasswordMinLength
	UsernameMinLength = Config.Business.UsernameMinLength
	SnowflakeMachineID = Config.Snowflake.MachineID
}
