package config

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"hi-go/src/utils/logger"
)

// Config 全局配置实例
var Config *AppConfig

// Init 初始化配置
// env 参数指定环境：dev, test, uat, prod
// 如果 env 为空，则从环境变量 GO_ENV 读取，默认为 dev
func Init(env string) error {
	if env == "" {
		env = GetEnv()
	}

	// 设置配置文件
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("../../configs")

	// 允许环境变量覆盖配置文件
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("读取配置文件失败", zap.Error(err))
		return err
	}

	// 将配置解析到结构体
	Config = &AppConfig{}
	if err := viper.Unmarshal(Config); err != nil {
		logger.Error("解析配置文件失败", zap.Error(err))
		return err
	}

	return nil
}

// GetEnv 获取当前运行环境
// 从环境变量 GO_ENV 读取，默认为 dev
func GetEnv() string {
	env := os.Getenv("GO_ENV")
	if env == "" {
		return "dev"
	}
	return env
}

// Reload 重新加载配置
// 用于运行时热更新配置
func Reload() error {
	return Init(GetEnv())
}
