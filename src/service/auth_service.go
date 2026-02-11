package service

import (
	"context"
	"errors"
	"fmt"
	"hi-go/src/config"
	"hi-go/src/model"
	"hi-go/src/repository"
	"hi-go/src/utils/jwt"
	"hi-go/src/utils/redis"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrUserDisabled       = errors.New("用户已被禁用")
	ErrUserExists         = errors.New("用户名已存在")
	ErrEmailExists        = errors.New("邮箱已存在")
)

// 认证服务
type AuthService struct {
	userRepo *repository.UserRepository
}

// 创建认证服务实例
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

// 登录用户
func (s *AuthService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	// 1. 查找用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 2. 检查用户状态
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}

	// 3. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 4. 生成 JWT Token
	userID := fmt.Sprintf("%d", user.ID)

	// 可以从数据库获取角色列表，这里简化为固定角色
	roles := []string{"user"}

	// 使用 GenerateTokenPair 生成 access token 和 refresh token
	accessToken, refreshToken, err := jwt.GenerateTokenPair(userID, user.Username, roles, nil)
	if err != nil {
		return nil, err
	}

	// 5. 要将accessToken信息存储到redis中，方便后续验证和刷新token时使用
	redisKey := fmt.Sprintf("jwt:access_token:%s", userID)
	err = redis.Set(context.Background(), redisKey, accessToken, config.GetRedisTokenTTL())
	if err != nil {
		return nil, err
	}

	// 6. 返回登录响应
	return &model.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    config.JWTAccessTokenDuration, // 使用配置常量
	}, nil
}

// 注册用户
func (s *AuthService) Register(req *model.RegisterRequest) (*model.User, error) {
	// 1. 检查用户名是否存在
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	// 2. 检查邮箱是否存在
	if req.Email != "" {
		exists, err := s.userRepo.ExistsByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEmailExists
		}
	}

	// 3. 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 4. 创建用户
	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// 根据ID获取用户信息
func (s *AuthService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
