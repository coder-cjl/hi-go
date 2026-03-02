package main

import (
	"fmt"
	"hi-go/src/config"
	"hi-go/src/model"
	"hi-go/src/router"
	"hi-go/src/service/aiservice"
	"hi-go/src/utils/jwt"
	"hi-go/src/utils/logger"
	"hi-go/src/utils/mysql"
	"hi-go/src/utils/redis"
	"hi-go/src/utils/snowflake"
	"hi-go/src/utils/yapi"

	"go.uber.org/zap"
)

// @title           Hi-Go API 文档
// @version         1.0
// @description     Hi-Go 项目的 RESTful API 接口文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description					JWT 认证，格式：Bearer {token}

// initConfig 初始化配置
func initConfig() {
	// 从环境变量 GO_ENV 读取环境，默认为 dev
	env := config.GetEnv()
	if err := config.Init(env); err != nil {
		logger.Fatalf("配置初始化失败: %v", err)
	}
	// 更新向后兼容的变量
	config.UpdateLegacyVars()
	logger.Infof("配置加载成功 [环境: %s]", env)
}

// initJWT 初始化JWT管理器
func initJWT() {
	jwt.Init(nil)
}

// initLogger 初始化日志
func initLogger() {
	cfg := &logger.Config{
		Level:      config.Config.Log.Level,
		Env:        config.Config.Server.Mode,
		Topic:      "hi-go",
		FilePath:   config.Config.Log.Filename,
		MaxSize:    config.Config.Log.MaxSize,
		MaxBackups: config.Config.Log.MaxBackups,
		MaxAge:     config.Config.Log.MaxAge,
		Compress:   config.Config.Log.Compress,
		// Elasticsearch 配置
		ESEnabled:   config.Config.Elasticsearch.Enabled,
		ESAddrs:     config.Config.Elasticsearch.Addrs,
		ESUsername:  config.Config.Elasticsearch.Username,
		ESPassword:  config.Config.Elasticsearch.Password,
		ESIndex:     config.Config.Elasticsearch.Index,
		ESMaxRetry:  config.Config.Elasticsearch.MaxRetry,
		ESBatchSize: 100, // 默认批量大小
		ESFlushTime: 5,   // 默认5秒刷新一次
		// Logstash 配置
		LogstashEnabled:    config.Config.Logstash.Enabled,
		LogstashHost:       config.Config.Logstash.Host,
		LogstashPort:       config.Config.Logstash.Port,
		LogstashProtocol:   config.Config.Logstash.Protocol,
		LogstashTimeout:    config.Config.Logstash.Timeout,
		LogstashReconnect:  config.Config.Logstash.Reconnect,
		LogstashBufferSize: config.Config.Logstash.BufferSize,
	}
	if err := logger.Init(cfg); err != nil {
		logger.Fatalf("日志初始化失败: %v", err)
	}
	defer logger.Sync()
	logger.Info("日志初始化成功",
		zap.String("level", cfg.Level),
		zap.String("file", cfg.FilePath),
		zap.Bool("es_enabled", cfg.ESEnabled))
}

// initMySQL 初始化MySQL数据库
func initMySQL() {
	dbCfg := config.Config.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.DBName,
		dbCfg.Charset,
	)

	cfg := &mysql.Config{
		DSN:             dsn,
		MaxOpenConns:    dbCfg.MaxOpenConns,
		MaxIdleConns:    dbCfg.MaxIdleConns,
		ConnMaxLifetime: config.GetDBConnMaxLifetime(),
		ConnMaxIdleTime: config.GetDBConnMaxIdleTime(),
	}

	if err := mysql.Init(cfg); err != nil {
		logger.Error("MySQL初始化失败", zap.Error(err))
		panic(err)
	}
	logger.Info("MySQL初始化成功", zap.String("database", dbCfg.DBName))
}

// initRedis 初始化Redis
func initRedis() {
	redisCfg := config.Config.Redis
	cfg := &redis.Config{
		Addr:     fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	}

	if err := redis.Init(cfg); err != nil {
		logger.Error("Redis初始化失败", zap.Error(err))
		panic(err)
	}
	logger.Info("Redis初始化成功", zap.String("addr", cfg.Addr))
}

// initSnowflake 初始化雪花ID生成器
func initSnowflake() {
	machineID := config.Config.Snowflake.MachineID
	if err := snowflake.Init(machineID); err != nil {
		logger.Error("雪花ID生成器初始化失败", zap.Error(err))
		panic(err)
	}
	logger.Info("雪花ID生成器初始化成功", zap.Int64("machineID", machineID))
}

// initDB 初始化数据库（迁移表结构）
func initDB() {
	// 自动迁移数据库表
	if err := mysql.Database.AutoMigrate(&model.User{}, &model.Home{}, &model.Webhook{}); err != nil {
		logger.Error("数据库迁移失败", zap.Error(err))
		panic(err)
	}
	logger.Info("数据库迁移成功")
}

// initYApiSync 同步 Swagger 文档到 YApi
func initYApiSync() {
	if err := yapi.SyncToYApi(); err != nil {
		logger.Warn("YApi 同步失败，但不影响服务启动", zap.Error(err))
	}
}

// initRouter 初始化路由并启动服务
func initRouter() {
	r := router.Setup()

	// 从配置读取端口
	port := ":" + config.Config.Server.Port
	logger.Info("服务启动",
		zap.String("port", port),
		zap.String("mode", config.Config.Server.Mode),
		zap.String("env", config.GetEnv()))

	if err := r.Run(port); err != nil {
		logger.Error("服务启动失败", zap.Error(err))
		panic(err)
	}
}

func main() {
	// 1. 初始化配置（必须最先执行）
	initConfig()

	// 2. 初始化JWT
	initJWT()

	// 3. 初始化日志
	initLogger()
	logger.Info("应用启动中...", zap.String("env", config.GetEnv()))

	// 4. 初始化MySQL
	initMySQL()

	// 5. 初始化Redis
	initRedis()

	// 6. 初始化雪花ID生成器
	initSnowflake()

	// 7. 初始化数据库（迁移表结构）
	initDB()

	// 8. 初始化AI服务
	aiservice.Init()

	// 9. 同步 Swagger 文档到 YApi（可选）
	initYApiSync()

	// 10. 设置路由并启动服务
	initRouter()
}
