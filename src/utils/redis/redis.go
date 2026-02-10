package redis

import (
	"context"
	"errors"
	"hi-go/src/utils/logger"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// 全局Redis客户端实例
var Client *redis.Client

// 错误定义
var (
	ErrNilValue    = errors.New("redis: 空值")
	ErrKeyNotFound = errors.New("redis: 键不存在")
	ErrClientNil   = errors.New("redis: 客户端未初始化")
	ErrInvalidType = errors.New("redis: 无效的类型")
	ErrEmptyKey    = errors.New("redis: 键不能为空")
)

// Redis配置结构
type Config struct {
	Addr         string        // Redis地址 (host:port)
	Password     string        // 密码（如果没有则为空）
	DB           int           // 数据库编号 (0-15)
	PoolSize     int           // 连接池大小
	MinIdleConns int           // 最小空闲连接数
	MaxRetries   int           // 最大重试次数
	DialTimeout  time.Duration // 连接超时
	ReadTimeout  time.Duration // 读取超时
	WriteTimeout time.Duration // 写入超时
	PoolTimeout  time.Duration // 连接池超时
}

// 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Addr:         "localhost:6379",
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 2,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	}
}

// 使用配置初始化Redis客户端
func Init(cfg *Config) error {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	Client = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolTimeout:  cfg.PoolTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		logger.Error("Redis 连接失败", zap.Error(err))
		return err
	}

	logger.Info("Redis 连接成功", zap.String("addr", cfg.Addr), zap.Int("db", cfg.DB))
	return nil
}

// 自动初始化默认Redis客户端
func init() {
	if err := Init(DefaultConfig()); err != nil {
		logger.Error("Redis 初始化失败", zap.Error(err))
	}
}

// 获取全局Redis客户端
func GetClient() *redis.Client {
	return Client
}

// 关闭Redis连接
func Close() error {
	if Client == nil {
		return nil
	}
	return Client.Close()
}

// 检查Redis连接是否正常（健康检查）
func Ping(ctx context.Context) error {
	if Client == nil {
		return ErrClientNil
	}
	return Client.Ping(ctx).Err()
}

// ==================== String 操作 ====================

//	设置键值对
//
// 参数:
//   - ctx: 上下文
//   - key: 键
//   - value: 值
//   - expiration: 过期时间，0表示永不过期
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if key == "" {
		return ErrEmptyKey
	}
	if Client == nil {
		return ErrClientNil
	}
	return Client.Set(ctx, key, value, expiration).Err()
}

//	获取键对应的值
//
// 参数:
//   - ctx: 上下文
//   - key: 键
//
// 返回:
//   - string: 值
//   - error: 错误信息
func Get(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}
	if Client == nil {
		return "", ErrClientNil
	}

	val, err := Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrKeyNotFound
	}
	return val, err
}

// 设置新值并返回旧值
func GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}
	if Client == nil {
		return "", ErrClientNil
	}
	return Client.GetSet(ctx, key, value).Result()
}

// 仅当键不存在时设置（分布式锁常用）
func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	if key == "" {
		return false, ErrEmptyKey
	}
	if Client == nil {
		return false, ErrClientNil
	}
	return Client.SetNX(ctx, key, value, expiration).Result()
}

// 设置键值对并指定过期时间（秒）
func SetEX(ctx context.Context, key string, value interface{}, seconds int) error {
	if key == "" {
		return ErrEmptyKey
	}
	if Client == nil {
		return ErrClientNil
	}
	return Client.SetEx(ctx, key, value, time.Duration(seconds)*time.Second).Err()
}

// 批量获取多个键的值
func MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	if len(keys) == 0 {
		return nil, ErrEmptyKey
	}
	if Client == nil {
		return nil, ErrClientNil
	}
	return Client.MGet(ctx, keys...).Result()
}

// 批量设置多个键值对
func MSet(ctx context.Context, pairs ...interface{}) error {
	if len(pairs) == 0 {
		return errors.New("pairs cannot be empty")
	}
	if Client == nil {
		return ErrClientNil
	}
	return Client.MSet(ctx, pairs...).Err()
}

// 将键的整数值加1
func Incr(ctx context.Context, key string) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.Incr(ctx, key).Result()
}

// 将键的整数值增加指定数值
func IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.IncrBy(ctx, key, value).Result()
}

// 将键的整数值减1
func Decr(ctx context.Context, key string) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.Decr(ctx, key).Result()
}

// 将键的整数值减少指定数值
func DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.DecrBy(ctx, key, value).Result()
}

// ==================== Key 操作 ====================

// 删除一个或多个键
func Del(ctx context.Context, keys ...string) (int64, error) {
	if len(keys) == 0 {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.Del(ctx, keys...).Result()
}

// 检查键是否存在
func Exists(ctx context.Context, keys ...string) (int64, error) {
	if len(keys) == 0 {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.Exists(ctx, keys...).Result()
}

// 设置键的过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	if key == "" {
		return false, ErrEmptyKey
	}
	if Client == nil {
		return false, ErrClientNil
	}
	return Client.Expire(ctx, key, expiration).Result()
}

// 设置键在指定时间过期
func ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	if key == "" {
		return false, ErrEmptyKey
	}
	if Client == nil {
		return false, ErrClientNil
	}
	return Client.ExpireAt(ctx, key, tm).Result()
}

// 获取键的剩余生存时间
func TTL(ctx context.Context, key string) (time.Duration, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.TTL(ctx, key).Result()
}

// 移除键的过期时间
func Persist(ctx context.Context, key string) (bool, error) {
	if key == "" {
		return false, ErrEmptyKey
	}
	if Client == nil {
		return false, ErrClientNil
	}
	return Client.Persist(ctx, key).Result()
}

// 查找所有符合给定模式的键
func Keys(ctx context.Context, pattern string) ([]string, error) {
	if Client == nil {
		return nil, ErrClientNil
	}
	return Client.Keys(ctx, pattern).Result()
}

// 重命名键
func Rename(ctx context.Context, key, newKey string) error {
	if key == "" || newKey == "" {
		return ErrEmptyKey
	}
	if Client == nil {
		return ErrClientNil
	}
	return Client.Rename(ctx, key, newKey).Err()
}

// 返回键存储的值的类型
func Type(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}
	if Client == nil {
		return "", ErrClientNil
	}
	return Client.Type(ctx, key).Result()
}

// ==================== Hash 操作 ====================

// 设置哈希表字段的值
func HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.HSet(ctx, key, values...).Result()
}

// 获取哈希表中指定字段的值
func HGet(ctx context.Context, key, field string) (string, error) {
	if key == "" || field == "" {
		return "", ErrEmptyKey
	}
	if Client == nil {
		return "", ErrClientNil
	}

	val, err := Client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", ErrKeyNotFound
	}
	return val, err
}

// 获取哈希表中所有字段和值
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}
	if Client == nil {
		return nil, ErrClientNil
	}
	return Client.HGetAll(ctx, key).Result()
}

// HMGet 批量获取哈希表中多个字段的值
func HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	if key == "" || len(fields) == 0 {
		return nil, ErrEmptyKey
	}
	if Client == nil {
		return nil, ErrClientNil
	}
	return Client.HMGet(ctx, key, fields...).Result()
}

// 批量设置哈希表中多个字段的值
func HMSet(ctx context.Context, key string, values ...interface{}) (bool, error) {
	if key == "" {
		return false, ErrEmptyKey
	}
	if Client == nil {
		return false, ErrClientNil
	}
	return Client.HMSet(ctx, key, values...).Result()
}

// 删除哈希表中一个或多个字段
func HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	if key == "" || len(fields) == 0 {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.HDel(ctx, key, fields...).Result()
}

// 检查哈希表中字段是否存在
func HExists(ctx context.Context, key, field string) (bool, error) {
	if key == "" || field == "" {
		return false, ErrEmptyKey
	}
	if Client == nil {
		return false, ErrClientNil
	}
	return Client.HExists(ctx, key, field).Result()
}

// 获取哈希表中所有字段
func HKeys(ctx context.Context, key string) ([]string, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}
	if Client == nil {
		return nil, ErrClientNil
	}
	return Client.HKeys(ctx, key).Result()
}

// 获取哈希表中所有值
func HVals(ctx context.Context, key string) ([]string, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}
	if Client == nil {
		return nil, ErrClientNil
	}
	return Client.HVals(ctx, key).Result()
}

// 获取哈希表中字段的数量
func HLen(ctx context.Context, key string) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.HLen(ctx, key).Result()
}

// 为哈希表中字段的整数值加上增量
func HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	if key == "" || field == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.HIncrBy(ctx, key, field, incr).Result()
}

// ==================== List 操作 ====================

// 将一个或多个值插入列表头部
func LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.LPush(ctx, key, values...).Result()
}

// 将一个或多个值插入列表尾部
func RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.RPush(ctx, key, values...).Result()
}

// 移除并返回列表的第一个元素
func LPop(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}
	if Client == nil {
		return "", ErrClientNil
	}

	val, err := Client.LPop(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrKeyNotFound
	}
	return val, err
}

// 移除并返回列表的最后一个元素
func RPop(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}
	if Client == nil {
		return "", ErrClientNil
	}

	val, err := Client.RPop(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrKeyNotFound
	}
	return val, err
}

// 获取列表指定范围内的元素
func LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}
	if Client == nil {
		return nil, ErrClientNil
	}
	return Client.LRange(ctx, key, start, stop).Result()
}

// 获取列表长度
func LLen(ctx context.Context, key string) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.LLen(ctx, key).Result()
}

// 移除列表中与参数值相等的元素
func LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.LRem(ctx, key, count, value).Result()
}

// 修剪列表，只保留指定区间内的元素
func LTrim(ctx context.Context, key string, start, stop int64) error {
	if key == "" {
		return ErrEmptyKey
	}
	if Client == nil {
		return ErrClientNil
	}
	return Client.LTrim(ctx, key, start, stop).Err()
}

// ==================== Set 操作 ====================

// 向集合添加一个或多个成员
func SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.SAdd(ctx, key, members...).Result()
}

// 获取集合中所有成员
func SMembers(ctx context.Context, key string) ([]string, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}
	if Client == nil {
		return nil, ErrClientNil
	}
	return Client.SMembers(ctx, key).Result()
}

// 判断元素是否是集合成员
func SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	if key == "" {
		return false, ErrEmptyKey
	}
	if Client == nil {
		return false, ErrClientNil
	}
	return Client.SIsMember(ctx, key, member).Result()
}

// 获取集合的成员数
func SCard(ctx context.Context, key string) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.SCard(ctx, key).Result()
}

// 移除集合中一个或多个成员
func SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.SRem(ctx, key, members...).Result()
}

// 移除并返回集合中的一个随机元素
func SPop(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}
	if Client == nil {
		return "", ErrClientNil
	}
	return Client.SPop(ctx, key).Result()
}

// 返回集合中一个或多个随机元素
func SRandMember(ctx context.Context, key string, count int64) ([]string, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}
	if Client == nil {
		return nil, ErrClientNil
	}
	return Client.SRandMemberN(ctx, key, count).Result()
}

// ==================== Sorted Set 操作 ====================

// 向有序集合添加一个或多个成员
func ZAdd(ctx context.Context, key string, members ...redis.Z) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.ZAdd(ctx, key, members...).Result()
}

// 返回有序集合中指定区间内的成员（按分数从小到大）
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}
	if Client == nil {
		return nil, ErrClientNil
	}
	return Client.ZRange(ctx, key, start, stop).Result()
}

// 返回有序集合中指定区间内的成员及分数
func ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}
	if Client == nil {
		return nil, ErrClientNil
	}
	return Client.ZRangeWithScores(ctx, key, start, stop).Result()
}

// 返回有序集合中指定区间内的成员（按分数从大到小）
func ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}
	if Client == nil {
		return nil, ErrClientNil
	}
	return Client.ZRevRange(ctx, key, start, stop).Result()
}

// 获取有序集合的成员数
func ZCard(ctx context.Context, key string) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.ZCard(ctx, key).Result()
}

// 获取有序集合中成员的分数
func ZScore(ctx context.Context, key, member string) (float64, error) {
	if key == "" || member == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.ZScore(ctx, key, member).Result()
}

// 移除有序集合中一个或多个成员
func ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.ZRem(ctx, key, members...).Result()
}

// 有序集合中对指定成员的分数加上增量
func ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	if key == "" || member == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.ZIncrBy(ctx, key, increment, member).Result()
}

// 返回有序集合中指定成员的排名（按分数从小到大）
func ZRank(ctx context.Context, key, member string) (int64, error) {
	if key == "" || member == "" {
		return 0, ErrEmptyKey
	}
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.ZRank(ctx, key, member).Result()
}

// ==================== 高级操作 ====================

// 执行管道操作
func Pipeline(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	if Client == nil {
		return nil, ErrClientNil
	}

	pipe := Client.Pipeline()
	if err := fn(pipe); err != nil {
		return nil, err
	}
	return pipe.Exec(ctx)
}

// 执行事务管道操作
func TxPipeline(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	if Client == nil {
		return nil, ErrClientNil
	}

	pipe := Client.TxPipeline()
	if err := fn(pipe); err != nil {
		return nil, err
	}
	return pipe.Exec(ctx)
}

// 监视一个或多个键（用于事务）
func Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	if Client == nil {
		return ErrClientNil
	}
	return Client.Watch(ctx, fn, keys...)
}

// 迭代当前数据库中的键
func Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	if Client == nil {
		return nil, 0, ErrClientNil
	}
	return Client.Scan(ctx, cursor, match, count).Result()
}

// 清空当前数据库
func FlushDB(ctx context.Context) error {
	if Client == nil {
		return ErrClientNil
	}
	return Client.FlushDB(ctx).Err()
}

// 清空所有数据库
func FlushAll(ctx context.Context) error {
	if Client == nil {
		return ErrClientNil
	}
	return Client.FlushAll(ctx).Err()
}

// 返回当前数据库的键数量
func DBSize(ctx context.Context) (int64, error) {
	if Client == nil {
		return 0, ErrClientNil
	}
	return Client.DBSize(ctx).Result()
}
