package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Manager 全局JWT管理器实例
var Manager *JWTManager

// 错误定义
var (
	ErrTokenExpired     = errors.New("token已过期")
	ErrTokenInvalid     = errors.New("token无效")
	ErrTokenMalformed   = errors.New("token格式错误")
	ErrTokenNotValidYet = errors.New("token尚未生效")
	ErrSecretKeyEmpty   = errors.New("密钥不能为空")
)

// Claims 自定义JWT载荷结构
type Claims struct {
	UserID               string                 `json:"user_id"`  // 用户ID
	Username             string                 `json:"username"` // 用户名
	Roles                []string               `json:"roles"`    // 用户角色列表
	Extra                map[string]interface{} `json:"extra"`    // 额外的自定义字段
	jwt.RegisteredClaims                        // JWT标准声明
}

// Config JWT配置
type Config struct {
	SecretKey            string        // 签名密钥
	AccessTokenDuration  time.Duration // 访问令牌过期时间
	RefreshTokenDuration time.Duration // 刷新令牌过期时间
	Issuer               string        // 签发者
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		SecretKey:            "your-secret-key-change-it", // 生产环境请务必修改
		AccessTokenDuration:  15 * time.Minute,            // 访问令牌15分钟
		RefreshTokenDuration: 7 * 24 * time.Hour,          // 刷新令牌7天
		Issuer:               "hi-go-app",
	}
}

// JWTManager JWT管理器
type JWTManager struct {
	config *Config
}

// Init 使用配置初始化JWT管理器
func Init(cfg *Config) {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	Manager = &JWTManager{config: cfg}
}

// init 自动初始化默认JWT管理器
func init() {
	Init(DefaultConfig())
}

// GetManager 获取全局JWT管理器实例
func GetManager() *JWTManager {
	return Manager
}

// NewJWTManager 创建JWT管理器
// 参数:
//   - config: JWT配置，如果为nil则使用默认配置
//
// 返回:
//   - *JWTManager: JWT管理器实例
func NewJWTManager(config *Config) *JWTManager {
	if config == nil {
		config = DefaultConfig()
	}
	return &JWTManager{config: config}
}

// SetSecretKey 设置签名密钥
// 参数:
//   - key: 签名密钥
func (m *JWTManager) SetSecretKey(key string) {
	m.config.SecretKey = key
}

// SetAccessTokenDuration 设置访问令牌过期时间
// 参数:
//   - duration: 过期时间
func (m *JWTManager) SetAccessTokenDuration(duration time.Duration) {
	m.config.AccessTokenDuration = duration
}

// SetRefreshTokenDuration 设置刷新令牌过期时间
// 参数:
//   - duration: 过期时间
func (m *JWTManager) SetRefreshTokenDuration(duration time.Duration) {
	m.config.RefreshTokenDuration = duration
}

// GenerateToken 生成访问令牌
// 参数:
//   - userID: 用户ID
//   - username: 用户名
//   - roles: 用户角色列表
//   - extra: 额外的自定义字段
//
// 返回:
//   - string: 生成的token字符串
//   - error: 错误信息
func (m *JWTManager) GenerateToken(userID, username string, roles []string, extra map[string]interface{}) (string, error) {
	if m.config.SecretKey == "" {
		return "", ErrSecretKeyEmpty
	}

	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		Extra:    extra,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.AccessTokenDuration)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.SecretKey))
}

// GenerateRefreshToken 生成刷新令牌
// 参数:
//   - userID: 用户ID
//
// 返回:
//   - string: 生成的刷新token字符串
//   - error: 错误信息
func (m *JWTManager) GenerateRefreshToken(userID string) (string, error) {
	if m.config.SecretKey == "" {
		return "", ErrSecretKeyEmpty
	}

	now := time.Now()
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.RefreshTokenDuration)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.SecretKey))
}

// GenerateTokenPair 同时生成访问令牌和刷新令牌
// 参数:
//   - userID: 用户ID
//   - username: 用户名
//   - roles: 用户角色列表
//   - extra: 额外的自定义字段
//
// 返回:
//   - accessToken: 访问令牌
//   - refreshToken: 刷新令牌
//   - error: 错误信息
func (m *JWTManager) GenerateTokenPair(userID, username string, roles []string, extra map[string]interface{}) (accessToken, refreshToken string, err error) {
	accessToken, err = m.GenerateToken(userID, username, roles, extra)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = m.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ParseToken 解析并验证token
// 参数:
//   - tokenString: token字符串
//
// 返回:
//   - *Claims: 解析后的声明
//   - error: 错误信息
func (m *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return []byte(m.config.SecretKey), nil
	})

	if err != nil {
		// 详细的错误处理
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// ValidateToken 验证token是否有效
// 参数:
//   - tokenString: token字符串
//
// 返回:
//   - bool: 是否有效
//   - error: 错误信息
func (m *JWTManager) ValidateToken(tokenString string) (bool, error) {
	_, err := m.ParseToken(tokenString)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetUserIDFromToken 从token中提取用户ID
// 参数:
//   - tokenString: token字符串
//
// 返回:
//   - string: 用户ID
//   - error: 错误信息
func (m *JWTManager) GetUserIDFromToken(tokenString string) (string, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// IsTokenExpired 检查token是否已过期
// 参数:
//   - tokenString: token字符串
//
// 返回:
//   - bool: 是否已过期
func (m *JWTManager) IsTokenExpired(tokenString string) bool {
	_, err := m.ParseToken(tokenString)
	return errors.Is(err, ErrTokenExpired)
}

// RefreshAccessToken 使用刷新令牌生成新的访问令牌
// 参数:
//   - refreshToken: 刷新令牌
//   - username: 用户名
//   - roles: 用户角色列表
//   - extra: 额外的自定义字段
//
// 返回:
//   - string: 新的访问令牌
//   - error: 错误信息
func (m *JWTManager) RefreshAccessToken(refreshToken, username string, roles []string, extra map[string]interface{}) (string, error) {
	// 验证刷新令牌
	claims, err := m.ParseToken(refreshToken)
	if err != nil {
		return "", err
	}

	// 使用刷新令牌中的用户ID生成新的访问令牌
	return m.GenerateToken(claims.UserID, username, roles, extra)
}

// GetClaims 获取token中的所有声明
// 参数:
//   - tokenString: token字符串
//
// 返回:
//   - *Claims: 声明对象
//   - error: 错误信息
func (m *JWTManager) GetClaims(tokenString string) (*Claims, error) {
	return m.ParseToken(tokenString)
}

// HasRole 检查token中是否包含指定角色
// 参数:
//   - tokenString: token字符串
//   - role: 要检查的角色
//
// 返回:
//   - bool: 是否包含该角色
//   - error: 错误信息
func (m *JWTManager) HasRole(tokenString, role string) (bool, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return false, err
	}

	for _, r := range claims.Roles {
		if r == role {
			return true, nil
		}
	}
	return false, nil
}

// ==================== 全局便捷方法 ====================
// 以下方法直接使用全局Manager实例，简化调用

// GenerateToken 使用全局管理器生成访问令牌
func GenerateToken(userID, username string, roles []string, extra map[string]interface{}) (string, error) {
	return Manager.GenerateToken(userID, username, roles, extra)
}

// GenerateRefreshToken 使用全局管理器生成刷新令牌
func GenerateRefreshToken(userID string) (string, error) {
	return Manager.GenerateRefreshToken(userID)
}

// GenerateTokenPair 使用全局管理器同时生成访问令牌和刷新令牌
func GenerateTokenPair(userID, username string, roles []string, extra map[string]interface{}) (accessToken, refreshToken string, err error) {
	return Manager.GenerateTokenPair(userID, username, roles, extra)
}

// ParseToken 使用全局管理器解析并验证token
func ParseToken(tokenString string) (*Claims, error) {
	return Manager.ParseToken(tokenString)
}

// ValidateToken 使用全局管理器验证token是否有效
func ValidateToken(tokenString string) (bool, error) {
	return Manager.ValidateToken(tokenString)
}

// GetUserIDFromToken 使用全局管理器从token中提取用户ID
func GetUserIDFromToken(tokenString string) (string, error) {
	return Manager.GetUserIDFromToken(tokenString)
}

// IsTokenExpired 使用全局管理器检查token是否已过期
func IsTokenExpired(tokenString string) bool {
	return Manager.IsTokenExpired(tokenString)
}

// RefreshAccessToken 使用全局管理器刷新访问令牌
func RefreshAccessToken(refreshToken, username string, roles []string, extra map[string]interface{}) (string, error) {
	return Manager.RefreshAccessToken(refreshToken, username, roles, extra)
}

// GetClaims 使用全局管理器获取token中的所有声明
func GetClaims(tokenString string) (*Claims, error) {
	return Manager.GetClaims(tokenString)
}

// HasRole 使用全局管理器检查token中是否包含指定角色
func HasRole(tokenString, role string) (bool, error) {
	return Manager.HasRole(tokenString, role)
}
