package mysql

import (
	"context"
	"errors"
	"hi-go/src/config"
	"hi-go/src/utils/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局数据库实例
var Database *gorm.DB

// 数据库配置结构
type Config struct {
	DSN             string        // 数据源名称
	MaxOpenConns    int           // 最大打开连接数
	MaxIdleConns    int           // 最大空闲连接数
	ConnMaxLifetime time.Duration // 连接最大存活时间
	ConnMaxIdleTime time.Duration // 空闲连接最大存活时间
}

// 默认配置
func DefaultConfig() *Config {
	return &Config{
		DSN:             "root:123456@~!@tcp(localhost:3306)/dev?charset=utf8mb4&parseTime=True&loc=Local",
		MaxOpenConns:    config.DBMaxOpenConns,
		MaxIdleConns:    config.DBMaxIdleConns,
		ConnMaxLifetime: config.DBConnMaxLifetime,
		ConnMaxIdleTime: config.DBConnMaxIdleTime,
	}
}

// 使用配置初始化数据库连接
func Init(cfg *Config) error {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		logger.Error("MySQL 连接失败", zap.Error(err))
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("获取底层 SQL DB 失败", zap.Error(err))
		return err
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	Database = db
	logger.Info("MySQL 连接成功")
	return nil
}

// ==================== 连接管理 ====================

// 获取数据库实例
func GetDB() *gorm.DB {
	return Database
}

// 关闭数据库连接
func Close() error {
	if Database == nil {
		return nil
	}
	sqlDB, err := Database.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// 检查数据库连接是否正常（健康检查）
func Ping() error {
	if Database == nil {
		return errors.New("数据库未初始化")
	}
	sqlDB, err := Database.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// 重新连接数据库
func Reconnect(cfg *Config) error {
	if err := Close(); err != nil {
		logger.Warn("关闭旧连接失败", zap.Error(err))
	}
	return Init(cfg)
}

// ==================== 事务封装 ====================

//	执行事务，自动提交或回滚
//
// fn 返回 error 时回滚，返回 nil 时提交
func Transaction(fn func(tx *gorm.DB) error) error {
	return Database.Transaction(fn)
}

//	开始一个事务，返回事务实例
//
// 需要手动调用 Commit() 或 Rollback()
func BeginTx() *gorm.DB {
	return Database.Begin()
}

// ==================== 通用 CRUD ====================

// 创建单条记录
func Create(data interface{}) error {
	return Database.Create(data).Error
}

// 批量创建记录
func CreateBatch(data interface{}, batchSize int) error {
	return Database.CreateInBatches(data, batchSize).Error
}

// 根据 ID 查询单条记录
func FindByID(dest interface{}, id interface{}) error {
	return Database.First(dest, id).Error
}

// 根据条件查询单条记录
func FindOne(dest interface{}, query interface{}, args ...interface{}) error {
	return Database.Where(query, args...).First(dest).Error
}

// 根据条件查询所有符合条件的记录
func FindAll(dest interface{}, query interface{}, args ...interface{}) error {
	return Database.Where(query, args...).Find(dest).Error
}

// 根据 ID 更新记录
func Update(model interface{}, id interface{}, updates map[string]interface{}) error {
	return Database.Model(model).Where("id = ?", id).Updates(updates).Error
}

// 根据条件更新记录
func UpdateByCondition(model interface{}, query interface{}, updates map[string]interface{}, args ...interface{}) error {
	return Database.Model(model).Where(query, args...).Updates(updates).Error
}

// 保存记录（存在则更新，不存在则创建）
func Save(data interface{}) error {
	return Database.Save(data).Error
}

// 根据 ID 删除记录（软删除，如果模型有 DeletedAt 字段）
func Delete(model interface{}, id interface{}) error {
	return Database.Delete(model, id).Error
}

// 根据条件删除记录
func DeleteByCondition(model interface{}, query interface{}, args ...interface{}) error {
	return Database.Where(query, args...).Delete(model).Error
}

// 硬删除记录（永久删除）
func HardDelete(model interface{}, id interface{}) error {
	return Database.Unscoped().Delete(model, id).Error
}

// 检查记录是否存在
func Exists(model interface{}, query interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := Database.Model(model).Where(query, args...).Count(&count).Error
	return count > 0, err
}

// 统计符合条件的记录数
func Count(model interface{}, query interface{}, args ...interface{}) (int64, error) {
	var count int64
	err := Database.Model(model).Where(query, args...).Count(&count).Error
	return count, err
}

// ==================== 分页查询 ====================

// 分页结果
type PageResult struct {
	Total    int64       // 总记录数
	Page     int         // 当前页码
	PageSize int         // 每页大小
	Pages    int         // 总页数
	Data     interface{} // 数据列表
}

//	分页查询
//
// dest: 结果切片指针，page: 页码（从1开始），pageSize: 每页数量
func Paginate(dest interface{}, model interface{}, page, pageSize int, query interface{}, args ...interface{}) (*PageResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大每页数量
	}

	var total int64
	db := Database.Model(model)

	if query != nil {
		db = db.Where(query, args...)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 计算总页数
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}

	// 查询数据
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(dest).Error; err != nil {
		return nil, err
	}

	return &PageResult{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
		Data:     dest,
	}, nil
}

// 带排序的分页查询
func PaginateWithOrder(dest interface{}, model interface{}, page, pageSize int, order string, query interface{}, args ...interface{}) (*PageResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var total int64
	db := Database.Model(model)

	if query != nil {
		db = db.Where(query, args...)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}

	offset := (page - 1) * pageSize
	if order != "" {
		db = db.Order(order)
	}
	if err := db.Offset(offset).Limit(pageSize).Find(dest).Error; err != nil {
		return nil, err
	}

	return &PageResult{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
		Data:     dest,
	}, nil
}

// ==================== 配置优化 ====================

// 设置最大打开连接数
func SetMaxOpenConns(n int) error {
	sqlDB, err := Database.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(n)
	return nil
}

// 设置最大空闲连接数
func SetMaxIdleConns(n int) error {
	sqlDB, err := Database.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(n)
	return nil
}

// 设置连接最大存活时间
func SetConnMaxLifetime(d time.Duration) error {
	sqlDB, err := Database.DB()
	if err != nil {
		return err
	}
	sqlDB.SetConnMaxLifetime(d)
	return nil
}

// 设置空闲连接最大存活时间
func SetConnMaxIdleTime(d time.Duration) error {
	sqlDB, err := Database.DB()
	if err != nil {
		return err
	}
	sqlDB.SetConnMaxIdleTime(d)
	return nil
}

// ==================== 数据库迁移 ====================

//	自动迁移表结构
//
// 根据模型结构自动创建或更新表
func AutoMigrate(models ...interface{}) error {
	return Database.AutoMigrate(models...)
}

// ==================== 原生 SQL ====================

// 执行原生 SQL 查询
func RawQuery(dest interface{}, sql string, values ...interface{}) error {
	return Database.Raw(sql, values...).Scan(dest).Error
}

// 执行原生 SQL（INSERT、UPDATE、DELETE）
func Exec(sql string, values ...interface{}) error {
	return Database.Exec(sql, values...).Error
}

// ==================== 辅助方法 ====================

// 应用查询作用域
func Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	return Database.Scopes(funcs...)
}

// 使用 context 创建新的 DB 实例
func WithContext(ctx context.Context) *gorm.DB {
	return Database.WithContext(ctx)
}

// 开启调试模式（打印 SQL）
func Debug() *gorm.DB {
	return Database.Debug()
}
