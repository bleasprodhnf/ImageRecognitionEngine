package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

// TestRedisConnection 测试Redis连接和基本操作
func TestRedisConnection(t *testing.T) {
	config, err := LoadTestConfig()
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: "", // 如果有密码，在此设置
		DB:       0,  // 使用默认DB
	})

	ctx := context.Background()

	// 测试Redis连接
	_, err = client.Ping(ctx).Result()
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}

	// 测试基本的Redis操作
	// 测试Set操作
	key := "test_key"
	value := "test_value"
	err = client.Set(ctx, key, value, 5*time.Minute).Err()
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	// 测试Get操作
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	if result != value {
		t.Errorf("Expected value '%s', got '%s'", value, result)
	}

	// 测试过期时间
	ttl, err := client.TTL(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to get TTL: %v", err)
	}

	if ttl.Minutes() <= 0 {
		t.Error("Expected TTL to be greater than 0")
	}

	// 测试删除操作
	err = client.Del(ctx, key).Err()
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	// 验证键已被删除
	_, err = client.Get(ctx, key).Result()
	if err != redis.Nil {
		t.Errorf("Expected key to be deleted, but it still exists")
	}

	// 关闭连接
	err = client.Close()
	if err != nil {
		t.Fatalf("Failed to close Redis connection: %v", err)
	}
}