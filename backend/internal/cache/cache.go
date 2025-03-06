package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/image-recognition-engine/config"
)

// CacheManager 缓存管理器接口
type CacheManager interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	GetStats() *CacheStats
}

// CacheStats 缓存统计信息
type CacheStats struct {
	Hits        int64     `json:"hits"`
	Misses      int64     `json:"misses"`
	KeyCount    int64     `json:"keyCount"`
	MemoryUsage int64     `json:"memoryUsage"`
	LastUpdate  time.Time `json:"lastUpdate"`
}

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
	stats  *CacheStats
}

// NewRedisCache 创建新的Redis缓存管理器
func NewRedisCache(cfg config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接Redis失败: %w", err)
	}

	return &RedisCache{
		client: client,
		stats: &CacheStats{
			LastUpdate: time.Now(),
		},
	}, nil
}

// Set 设置缓存
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %w", err)
	}

	err = rc.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return fmt.Errorf("设置缓存失败: %w", err)
	}

	return nil
}

// Get 获取缓存
func (rc *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
	data, err := rc.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			rc.stats.Misses++
			return fmt.Errorf("缓存不存在")
		}
		return fmt.Errorf("获取缓存失败: %w", err)
	}

	err = json.Unmarshal(data, value)
	if err != nil {
		return fmt.Errorf("反序列化数据失败: %w", err)
	}

	rc.stats.Hits++
	return nil
}

// Delete 删除缓存
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	err := rc.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("删除缓存失败: %w", err)
	}
	return nil
}

// Exists 检查缓存是否存在
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("检查缓存是否存在失败: %w", err)
	}
	return result > 0, nil
}

// GetStats 获取缓存统计信息
func (rc *RedisCache) GetStats() *CacheStats {
	ctx := context.Background()

	// 更新统计信息
	info, err := rc.client.Info(ctx, "memory").Result()
	if err == nil {
		// 解析内存使用信息
		// 这里需要根据实际Redis INFO命令返回格式解析
		_ = info // TODO: 解析内存使用信息
	}

	// 获取键数量
	keyCount, err := rc.client.DBSize(ctx).Result()
	if err == nil {
		rc.stats.KeyCount = keyCount
	}

	rc.stats.LastUpdate = time.Now()
	return rc.stats
}