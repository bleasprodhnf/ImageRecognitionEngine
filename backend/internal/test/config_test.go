package test

import (
	"os"
	"testing"
)

// TestConfig 测试配置结构体
type TestConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisHost  string
	RedisPort  string
}

// LoadTestConfig 加载测试配置
func LoadTestConfig() (*TestConfig, error) {
	// 从环境变量或配置文件加载测试配置
	return &TestConfig{
		DBHost:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		DBPort:     getEnvOrDefault("TEST_DB_PORT", "3306"),
		DBUser:     getEnvOrDefault("TEST_DB_USER", "test"),
		DBPassword: getEnvOrDefault("TEST_DB_PASSWORD", "test"),
		DBName:     getEnvOrDefault("TEST_DB_NAME", "image_recognition_test"),
		RedisHost:  getEnvOrDefault("TEST_REDIS_HOST", "localhost"),
		RedisPort:  getEnvOrDefault("TEST_REDIS_PORT", "6379"),
	}, nil
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestLoadTestConfig 测试配置加载功能
func TestLoadTestConfig(t *testing.T) {
	config, err := LoadTestConfig()
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	if config.DBHost == "" {
		t.Error("DBHost should not be empty")
	}

	if config.DBPort == "" {
		t.Error("DBPort should not be empty")
	}

	// 其他配置项的测试...
}