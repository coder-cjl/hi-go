package main

import (
	"hi-go/src/config"
	"hi-go/src/model"
	"hi-go/src/router"
	"hi-go/src/utils/logger"
	"hi-go/src/utils/mysql"
	"hi-go/src/utils/snowflake"

	"go.uber.org/zap"
)

func initLogger() {
	defer logger.Sync()
}

func initSnowflake() {
	// 初始化雪花ID生成器
	if err := snowflake.Init(config.SnowflakeMachineID); err != nil {
		logger.Error("雪花ID生成器初始化失败", zap.Error(err))
		panic(err)
	}
	logger.Info("雪花ID生成器初始化成功", zap.Int64("machineID", config.SnowflakeMachineID))
}

func initDB() {
	// 自动迁移数据库表
	if err := mysql.Database.AutoMigrate(&model.User{}); err != nil {
		logger.Error("数据库迁移失败", zap.Error(err))
		panic(err)
	}
	logger.Info("数据库迁移成功")
}

func initRouter() {
	r := router.Setup()

	// 启动服务
	port := ":8000"
	logger.Info("服务启动", zap.String("port", port))
	if err := r.Run(port); err != nil {
		logger.Error("服务启动失败", zap.Error(err))
		panic(err)
	}
}

func main() {
	// 初始化日志
	initLogger()
	logger.Info("应用启动中...")

	// 初始化雪花ID生成器
	initSnowflake()

	// 初始化数据库
	initDB()

	// 设置路由
	initRouter()
}
