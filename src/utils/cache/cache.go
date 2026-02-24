package cache

import (
	"context"
	"time"
)

// Cache 缓存接口
type Cache interface {
	// Get 获取缓存
	Get(ctx context.Context, key string) (string, error)

	// Set 设置缓存
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// Delete 删除缓存
	Delete(ctx context.Context, key string) error

	// Exists 检查是否存在
	Exists(ctx context.Context, key string) (bool, error)
}
