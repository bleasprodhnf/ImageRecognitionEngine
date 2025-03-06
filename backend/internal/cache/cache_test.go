package cache

import (
	"context"
	"testing"
	"time"

	"github.com/image-recognition-engine/config"
	"github.com/stretchr/testify/assert"
)

func TestRedisCache(t *testing.T) {
	// 测试配置
	cfg := config.RedisConfig{
		Host:     "tp-redis-1",
		Port:     6379,
		Password: "",
		DB:       0,
		PoolSize: 10,
	}

	// 创建缓存管理器
	cache, err := NewRedisCache(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, cache)

	// 测试上下文
	ctx := context.Background()

	// 测试数据
	testKey := "test_key"
	testValue := map[string]interface{}{
		"name": "测试数据",
		"age":  float64(25),
		"time": time.Now(),
	}

	// 测试Set方法
	t.Run("Test Set", func(t *testing.T) {
		err := cache.Set(ctx, testKey, testValue, time.Minute)
		assert.NoError(t, err)
	})

	// 测试Exists方法
	t.Run("Test Exists", func(t *testing.T) {
		exists, err := cache.Exists(ctx, testKey)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	// 测试Get方法
	t.Run("Test Get", func(t *testing.T) {
		var result map[string]interface{}
		err := cache.Get(ctx, testKey, &result)
		assert.NoError(t, err)
		assert.Equal(t, testValue["name"], result["name"])
		assert.Equal(t, testValue["age"], result["age"])
	})

	// 测试Delete方法
	t.Run("Test Delete", func(t *testing.T) {
		err := cache.Delete(ctx, testKey)
		assert.NoError(t, err)

		// 验证删除后不存在
		exists, err := cache.Exists(ctx, testKey)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	// 测试GetStats方法
	t.Run("Test GetStats", func(t *testing.T) {
		stats := cache.GetStats()
		assert.NotNil(t, stats)
		assert.GreaterOrEqual(t, stats.Hits, int64(0))
		assert.GreaterOrEqual(t, stats.Misses, int64(0))
		assert.NotZero(t, stats.LastUpdate)
	})
}