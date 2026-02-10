package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	log   *zap.Logger
	sugar *zap.SugaredLogger
)

// Config 日志配置结构
type Config struct {
	Level      string // 日志级别: debug, info, warn, error, fatal
	Env        string // 环境: dev, prod
	Topic      string // 日志主题
	FilePath   string // 日志文件路径（为空则只输出到控制台）
	MaxSize    int    // 单个日志文件最大 MB
	MaxBackups int    // 保留的旧日志文件数量
	MaxAge     int    // 保留的旧日志文件最大天数
	Compress   bool   // 是否压缩旧日志
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:      "debug",
		Env:        "dev",
		Topic:      "LUCA",
		FilePath:   "",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   false,
	}
}

// Init 使用自定义配置初始化 logger
func Init(cfg *Config) error {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// 解析日志级别
	level := zapcore.DebugLevel
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	// 编码器配置
	var encoderConfig zapcore.EncoderConfig
	if cfg.Env == "prod" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 开发环境使用彩色输出
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 使用 ISO8601 时间格式

	// 编码器
	var encoder zapcore.Encoder
	if cfg.Env == "prod" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 输出位置
	var cores []zapcore.Core

	// 控制台输出
	consoleCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)
	cores = append(cores, consoleCore)

	// 文件输出（如果配置了文件路径）
	if cfg.FilePath != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}

		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(fileWriter),
			level,
		)
		cores = append(cores, fileCore)
	}

	// 组合多个 Core
	core := zapcore.NewTee(cores...)

	// 创建 Logger
	log = zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel), // Error 级别及以上才打印堆栈
	).Named(cfg.Topic) // 使用 Named 将 topic 作为 logger 名称

	sugar = log.Sugar()
	return nil
}

func init() {
	// 使用默认配置初始化
	if err := Init(DefaultConfig()); err != nil {
		panic(err)
	}
}

// Sync 刷新日志缓冲区
func Sync() {
	if log != nil {
		log.Sync()
	}
	if sugar != nil {
		sugar.Sync()
	}
}

// GetLogger 获取底层 zap.Logger
func GetLogger() *zap.Logger {
	return log
}

// GetSugar 获取 SugaredLogger
func GetSugar() *zap.SugaredLogger {
	return sugar
}

// 结构化日志方法
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func Warning(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// 格式化日志方法（使用 SugaredLogger）
func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}

// With 创建带有预设字段的 logger
func With(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}
