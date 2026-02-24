package aiservice

import (
	"hi-go/src/config"
	"hi-go/src/utils/cache"
	"hi-go/src/utils/logger"
	"hi-go/src/utils/redis"

	"go.uber.org/zap"
)

var (
	// GlobalService 全局AI服务实例
	GlobalService *Service
)

// Init 初始化AI服务
func Init() {
	if !config.Config.AI.Enabled {
		logger.Info("AI服务未启用")
		return
	}

	logger.Info("正在初始化AI服务...",
		zap.String("provider", config.Config.AI.Provider),
		zap.String("model", config.Config.AI.Model))

	// 创建缓存（使用Redis）
	redisCache := cache.NewRedisCache(redis.Client)

	// 创建技能注册表
	registry := NewSkillRegistry()

	// 注册天气技能
	if config.Config.AI.Skills.Weather.Enabled {
		weatherSkill := NewWeatherSkill(config.Config.AI.Skills.Weather, redisCache)
		registry.Register(weatherSkill)
		logger.Info("天气查询技能已注册",
			zap.String("provider", config.Config.AI.Skills.Weather.Provider))
	}

	// 创建AI客户端
	var client AIClient
	switch config.Config.AI.Provider {
	case "deepseek":
		client = NewDeepSeekClient(config.Config.AI)
		logger.Info("使用DeepSeek作为AI提供商")
	case "openai":
		// 可以扩展支持OpenAI
		logger.Warn("OpenAI提供商暂未实现，使用DeepSeek代替")
		client = NewDeepSeekClient(config.Config.AI)
	default:
		logger.Warn("未知的AI提供商，跳过AI服务初始化",
			zap.String("provider", config.Config.AI.Provider))
		return
	}

	// 创建AI服务
	GlobalService = NewService(config.Config.AI, registry, client)

	logger.Info("AI服务初始化成功",
		zap.String("provider", config.Config.AI.Provider),
		zap.String("model", config.Config.AI.Model))
}
