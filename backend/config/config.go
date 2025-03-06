package config

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

// Config 系统配置结构
type Config struct {
	Server   ServerConfig   `json:"server" validate:"required"`
	Database DatabaseConfig `json:"database" validate:"required"`
	Redis    RedisConfig    `json:"redis" validate:"required"`
	JWT      JWTConfig      `json:"jwt" validate:"required"`
	Storage  StorageConfig  `json:"storage" validate:"required"`
	Model    ModelConfig    `json:"model" validate:"required"`

	// 内部使用，不导出
	configPath string
	lastMod    time.Time
	mutex      sync.RWMutex
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port" validate:"port_range"`
	Mode         string `json:"mode" validate:"server_mode"` // development, production
	ReadTimeout  int    `json:"readTimeout"`  // 秒
	WriteTimeout int    `json:"writeTimeout"` // 秒
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MySQL MySQLConfig `json:"mysql"`
	Mongo MongoConfig  `json:"mongo"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	MaxOpen  int    `json:"maxOpen"`  // 最大连接数
	MaxIdle  int    `json:"maxIdle"`  // 最大空闲连接数
	Lifetime int    `json:"lifetime"` // 连接最大生存时间(秒)
}

// MongoConfig MongoDB配置
type MongoConfig struct {
	URI      string `json:"uri"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	MaxPool  int    `json:"maxPool"` // 最大连接池大小
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	PoolSize int    `json:"poolSize"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `json:"secret"`
	Expiration int    `json:"expiration"` // 过期时间(小时)
}

// StorageConfig 存储配置
type StorageConfig struct {
    Type      string    `json:"type" validate:"required,oneof=local s3"`
    LocalPath string    `json:"localPath"`
    S3        S3Config  `json:"s3"`
    MaxSize   int64     `json:"maxSize" validate:"required"` // 文件大小限制(MB)
    Timeout   int       `json:"timeout" validate:"required"` // 上传超时时间(秒)
}

// S3Config S3存储配置
type S3Config struct {
    Endpoint  string `json:"endpoint"`
    Bucket    string `json:"bucket"`
    AccessKey string `json:"accessKey"`
    SecretKey string `json:"secretKey"`
    Region    string `json:"region"`
}

// ModelConfig 模型配置
type ModelConfig struct {
	BasePath    string `json:"basePath"`
	DefaultModel string `json:"defaultModel"`
	MaxConcurrent int    `json:"maxConcurrent"`
}

var (
	config *Config
	once   sync.Once
	
	// 配置文件监控间隔
	configCheckInterval = 5 * time.Second
)

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	var err error
	once.Do(func() {
		config = &Config{}
		configPath := getConfigPath()
		
		// 读取配置文件
		data, err := os.ReadFile(configPath)
		if err != nil {
			log.Printf("读取配置文件失败: %v", err)
			return
		}
		
		// 解析JSON
		err = json.Unmarshal(data, config)
		if err != nil {
			log.Printf("解析配置文件失败: %v", err)
			return
		}
		
		// 记录配置文件路径和修改时间
		config.configPath = configPath
		fileInfo, err := os.Stat(configPath)
		if err == nil {
			config.lastMod = fileInfo.ModTime()
		}
		
		// 从环境变量覆盖配置
		overrideFromEnv(config)
		
		// 解密敏感信息
		err = decryptSensitiveInfo(config)
		if err != nil {
			log.Printf("解密敏感信息失败: %v", err)
			return
		}
		
		// 验证配置
		err = validateConfig(config)
		if err != nil {
			log.Printf("配置验证失败: %v", err)
			return
		}
		
		// 这里注释掉监控配置文件的代码，因为该函数未定义
		// go monitorConfigFile()
	})

	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	return config
}

// decrypt 解密字符串
func decrypt(encryptedData, key string) (string, error) {
	// 解码Base64
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("base64解码失败: %v", err)
	}
	
	// 创建解密器
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("创建解密器失败: %v", err)
	}
	
	// 解密
	mode := cipher.NewCBCDecrypter(block, []byte(key)[:block.BlockSize()])
	mode.CryptBlocks(data, data)
	
	// 去除填充
	unpadded, err := pkcs7Unpad(data)
	if err != nil {
		return "", fmt.Errorf("去除填充失败: %v", err)
	}
	return string(unpadded), nil
}

// decryptSensitiveInfo 解密敏感信息
func decryptSensitiveInfo(cfg *Config) error {
	// 使用环境变量中的加密密钥
	key := os.Getenv("CONFIG_ENCRYPTION_KEY")
	if key == "" {
		return nil
	}

	var err error

	// 解密数据库密码
	if cfg.Database.MySQL.Password != "" {
		cfg.Database.MySQL.Password, err = decrypt(cfg.Database.MySQL.Password, key)
		if err != nil {
			return fmt.Errorf("MySQL密码解密失败: %v", err)
		}
	}
	if cfg.Database.Mongo.Password != "" {
		cfg.Database.Mongo.Password, err = decrypt(cfg.Database.Mongo.Password, key)
		if err != nil {
			return fmt.Errorf("MongoDB密码解密失败: %v", err)
		}
	}

	// 解密Redis密码
	if cfg.Redis.Password != "" {
		cfg.Redis.Password, err = decrypt(cfg.Redis.Password, key)
		if err != nil {
			return fmt.Errorf("Redis密码解密失败: %v", err)
		}
	}

	// 解密JWT密钥
	if cfg.JWT.Secret != "" {
		cfg.JWT.Secret, err = decrypt(cfg.JWT.Secret, key)
		if err != nil {
			return fmt.Errorf("JWT密钥解密失败: %v", err)
		}
	}

	// 解密存储凭证
	if cfg.Storage.S3.SecretKey != "" {
		cfg.Storage.S3.SecretKey, err = decrypt(cfg.Storage.S3.SecretKey, key)
		if err != nil {
			return fmt.Errorf("S3密钥解密失败: %v", err)
		}
	}

	return nil
}

// validateConfig 验证配置
func validateConfig(cfg *Config) error {
	validate := validator.New()
	
	// 注册自定义验证器
	validate.RegisterValidation("port_range", validatePortRange)
	validate.RegisterValidation("server_mode", validateServerMode)
	validate.RegisterValidation("storage_type", validateStorageType)
	
	// 注册结构体验证
	validate.RegisterStructValidation(serverConfigValidation, ServerConfig{})
	
	// 注册结构体验证
	validate.RegisterStructValidation(databaseConfigValidation, DatabaseConfig{})
	
	// 验证整个配置
	err := validate.Struct(cfg)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			return fmt.Errorf("配置验证失败: %s.%s: %s", e.Namespace(), e.Field(), e.Tag())
		}
		return err
	}
	
	// 验证存储配置
	err = validateStorageConfig(cfg.Storage)
	if err != nil {
		return err
	}
	
	// 验证模型配置
	return validateModelConfig(cfg.Model)
}

// validatePortRange 验证端口范围
func validatePortRange(fl validator.FieldLevel) bool {
	port := fl.Field().Int()
	return port > 0 && port < 65536
}

// validateServerMode 验证服务器模式
func validateServerMode(fl validator.FieldLevel) bool {
	mode := fl.Field().String()
	return mode == "development" || mode == "production"
}

// validateStorageType 验证存储类型
func validateStorageType(fl validator.FieldLevel) bool {
	storageType := fl.Field().String()
	return storageType == "local" || storageType == "s3" || storageType == "oss"
}

// serverConfigValidation 服务器配置验证
func serverConfigValidation(sl validator.StructLevel) {
	config := sl.Current().Interface().(ServerConfig)

	// 验证超时设置
	if config.ReadTimeout <= 0 {
		sl.ReportError(config.ReadTimeout, "ReadTimeout", "readTimeout", "positive", "")
	}
	if config.WriteTimeout <= 0 {
		sl.ReportError(config.WriteTimeout, "WriteTimeout", "writeTimeout", "positive", "")
	}
}

// databaseConfigValidation 数据库配置验证
func databaseConfigValidation(sl validator.StructLevel) {
	config := sl.Current().Interface().(DatabaseConfig)

	// 验证MySQL配置
	if config.MySQL.MaxOpen < config.MySQL.MaxIdle {
		sl.ReportError(config.MySQL.MaxOpen, "MaxOpen", "maxOpen", "gtefield", "maxIdle")
	}
	if config.MySQL.Lifetime <= 0 {
		sl.ReportError(config.MySQL.Lifetime, "Lifetime", "lifetime", "positive", "")
	}

	// 验证MongoDB配置
	if config.Mongo.MaxPool <= 0 {
		sl.ReportError(config.Mongo.MaxPool, "MaxPool", "maxPool", "positive", "")
	}
}

// validateStorageConfig 验证存储配置
func validateStorageConfig(config StorageConfig) error {
	switch config.Type {
	case "local":
		if config.LocalPath == "" {
			return fmt.Errorf("本地存储路径不能为空")
		}
		// 验证本地路径是否存在
		if _, err := os.Stat(config.LocalPath); os.IsNotExist(err) {
			return fmt.Errorf("本地存储路径不存在: %s", config.LocalPath)
		}
	case "s3":
		if config.S3.Endpoint == "" {
			return fmt.Errorf("S3终端节点不能为空")
		}
		if config.S3.Bucket == "" {
			return fmt.Errorf("S3存储桶不能为空")
		}
		if config.S3.AccessKey == "" || config.S3.SecretKey == "" {
			return fmt.Errorf("S3访问凭证不能为空")
		}
	case "oss":
		return fmt.Errorf("暂不支持OSS存储")
	default:
		return fmt.Errorf("不支持的存储类型: %s", config.Type)
	}
	return nil
}

// validateModelConfig 验证模型配置
func validateModelConfig(config ModelConfig) error {
	if config.BasePath == "" {
		return fmt.Errorf("模型基础路径不能为空")
	}
	if config.DefaultModel == "" {
		return fmt.Errorf("默认模型不能为空")
	}
	if config.MaxConcurrent <= 0 {
		return fmt.Errorf("最大并发数必须大于0")
	}
	// 验证模型基础路径是否存在
	if _, err := os.Stat(config.BasePath); os.IsNotExist(err) {
		return fmt.Errorf("模型基础路径不存在: %s", config.BasePath)
	}
	return nil
}

// encrypt 加密字符串
func encrypt(data, key string) string {
	// 创建加密器
	block, _ := aes.NewCipher([]byte(key))
	
	// 填充数据
	padded := pkcs7Pad([]byte(data), block.BlockSize())
	
	// 加密
	encrypted := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, []byte(key)[:block.BlockSize()])
	mode.CryptBlocks(encrypted, padded)
	
	// 编码为Base64
	return base64.StdEncoding.EncodeToString(encrypted)
}

// getConfigPath 获取配置文件路径
func getConfigPath() string {
	// 首先检查环境变量
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}

	// 然后检查当前目录和上级目录
	paths := []string{
		"./config.json",
		"../config.json",
		"./config/config.json",
		"../config/config.json",
	}

	for _, path := range paths {
		absPath, _ := filepath.Abs(path)
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	return ""
}

// overrideFromEnv 从环境变量覆盖配置
func overrideFromEnv(cfg *Config) {
	// 服务器配置
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}
	if mode := os.Getenv("SERVER_MODE"); mode != "" {
		cfg.Server.Mode = mode
	}
	if timeout := os.Getenv("SERVER_READ_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.Server.ReadTimeout = t
		}
	}
	if timeout := os.Getenv("SERVER_WRITE_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.Server.WriteTimeout = t
		}
	}

	// MySQL配置
	if host := os.Getenv("MYSQL_HOST"); host != "" {
		cfg.Database.MySQL.Host = host
	}
	if port := os.Getenv("MYSQL_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Database.MySQL.Port = p
		}
	}
	if username := os.Getenv("MYSQL_USERNAME"); username != "" {
		cfg.Database.MySQL.Username = username
	}
	if password := os.Getenv("MYSQL_PASSWORD"); password != "" {
		cfg.Database.MySQL.Password = password
	}
	if database := os.Getenv("MYSQL_DATABASE"); database != "" {
		cfg.Database.MySQL.Database = database
	}
	if maxOpen := os.Getenv("MYSQL_MAX_OPEN"); maxOpen != "" {
		if mo, err := strconv.Atoi(maxOpen); err == nil {
			cfg.Database.MySQL.MaxOpen = mo
		}
	}
	if maxIdle := os.Getenv("MYSQL_MAX_IDLE"); maxIdle != "" {
		if mi, err := strconv.Atoi(maxIdle); err == nil {
			cfg.Database.MySQL.MaxIdle = mi
		}
	}
	if lifetime := os.Getenv("MYSQL_LIFETIME"); lifetime != "" {
		if lt, err := strconv.Atoi(lifetime); err == nil {
			cfg.Database.MySQL.Lifetime = lt
		}
	}

	// MongoDB配置
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		cfg.Database.Mongo.URI = uri
	}
	if database := os.Getenv("MONGO_DATABASE"); database != "" {
		cfg.Database.Mongo.Database = database
	}
	if username := os.Getenv("MONGO_USERNAME"); username != "" {
		cfg.Database.Mongo.Username = username
	}
	if password := os.Getenv("MONGO_PASSWORD"); password != "" {
		cfg.Database.Mongo.Password = password
	}
	if maxPool := os.Getenv("MONGO_MAX_POOL"); maxPool != "" {
		if mp, err := strconv.Atoi(maxPool); err == nil {
			cfg.Database.Mongo.MaxPool = mp
		}
	}

	// Redis配置
	if host := os.Getenv("REDIS_HOST"); host != "" {
		cfg.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Redis.Port = p
		}
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		cfg.Redis.Password = password
	}
	if db := os.Getenv("REDIS_DB"); db != "" {
		if d, err := strconv.Atoi(db); err == nil {
			cfg.Redis.DB = d
		}
	}
	if poolSize := os.Getenv("REDIS_POOL_SIZE"); poolSize != "" {
		if ps, err := strconv.Atoi(poolSize); err == nil {
			cfg.Redis.PoolSize = ps
		}
	}

	// JWT配置
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWT.Secret = secret
	}
	if expiration := os.Getenv("JWT_EXPIRATION"); expiration != "" {
		if exp, err := strconv.Atoi(expiration); err == nil {
			cfg.JWT.Expiration = exp
		}
	}

	// 存储配置
	if storageType := os.Getenv("STORAGE_TYPE"); storageType != "" {
		cfg.Storage.Type = storageType
	}
	if localPath := os.Getenv("STORAGE_LOCAL_PATH"); localPath != "" {
		cfg.Storage.LocalPath = localPath
	}
	if endpoint := os.Getenv("STORAGE_S3_ENDPOINT"); endpoint != "" {
		cfg.Storage.S3.Endpoint = endpoint
	}
	if bucket := os.Getenv("STORAGE_S3_BUCKET"); bucket != "" {
		cfg.Storage.S3.Bucket = bucket
	}
	if accessKey := os.Getenv("STORAGE_S3_ACCESS_KEY"); accessKey != "" {
		cfg.Storage.S3.AccessKey = accessKey
	}
	if secretKey := os.Getenv("STORAGE_S3_SECRET_KEY"); secretKey != "" {
		cfg.Storage.S3.SecretKey = secretKey
	}
	if region := os.Getenv("STORAGE_S3_REGION"); region != "" {
		cfg.Storage.S3.Region = region
	}

	// 模型配置
	if basePath := os.Getenv("MODEL_BASE_PATH"); basePath != "" {
		cfg.Model.BasePath = basePath
	}
	if defaultModel := os.Getenv("MODEL_DEFAULT_MODEL"); defaultModel != "" {
		cfg.Model.DefaultModel = defaultModel
	}
	if maxConcurrent := os.Getenv("MODEL_MAX_CONCURRENT"); maxConcurrent != "" {
		if mc, err := strconv.Atoi(maxConcurrent); err == nil {
			cfg.Model.MaxConcurrent = mc
		}
	}
}

// pkcs7Pad 填充数据
func pkcs7Pad(data []byte, blockSize int) []byte {
	padLen := blockSize - len(data)%blockSize
	pad := make([]byte, padLen)
	for i := 0; i < padLen; i++ {
		pad[i] = byte(padLen)
	}
	return append(data, pad...)
}

// pkcs7Unpad 去除填充
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) {
		return nil, errors.New("invalid padding size")
	}
	return data[:len(data)-padding], nil
}