package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/image-recognition-engine/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/redis/go-redis/v9"
)

// 全局数据库连接实例
var (
	MySQLDB *sql.DB
	MongoDB *mongo.Database
	RedisClient *redis.Client
	dbStats *DatabaseStats
)

// DatabaseStats 数据库统计信息
type DatabaseStats struct {
	MySQLStats MySQLStats
	MongoStats MongoStats
}

// MySQLStats MySQL统计信息
type MySQLStats struct {
	OpenConnections int
	InUseConnections int
	IdleConnections int
	WaitCount int64
	WaitDuration time.Duration
	MaxIdleTimeClosed int64
}

// MongoStats MongoDB统计信息
type MongoStats struct {
	ActiveConnections int32
	AvailableConnections int32
	TotalCreated int64
}

// InitDatabase 初始化所有数据库连接
func InitDatabase(cfg *config.Config) error {
	// 初始化MySQL连接
	if err := initMySQL(cfg.Database.MySQL); err != nil {
		return fmt.Errorf("初始化MySQL失败: %w", err)
	}

	// 初始化MongoDB连接
	if err := initMongoDB(cfg.Database.Mongo); err != nil {
		return fmt.Errorf("初始化MongoDB失败: %w", err)
	}

	// 初始化Redis连接
	if err := initRedis(cfg.Redis); err != nil {
		return fmt.Errorf("初始化Redis失败: %w", err)
	}

	return nil
}

// OptimizeMySQL 优化MySQL配置和性能
func OptimizeMySQL() error {
	// 设置慢查询日志
	_, err := MySQLDB.Exec("SET GLOBAL slow_query_log = 1")
	if err != nil {
		log.Printf("启用慢查询日志失败: %v", err)
	}

	// 设置查询缓存
	_, err = MySQLDB.Exec("SET GLOBAL query_cache_size = 67108864") // 64MB
	if err != nil {
		log.Printf("设置查询缓存失败: %v", err)
	}

	return nil
}

// OptimizeMongoDB 优化MongoDB配置和性能
func OptimizeMongoDB() error {
	// 创建复合索引
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"created_at", 1},
			{"status", 1},
		},
		Options: options.Index().SetBackground(true),
	}

	_, err := MongoDB.Collection("recognition_results").Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Printf("创建MongoDB索引失败: %v", err)
		return err
	}

	return nil
}

// GetDatabaseStats 获取数据库统计信息
func GetDatabaseStats() *DatabaseStats {
	// 获取MySQL统计信息
	mysqlStats := MySQLStats{
		OpenConnections: MySQLDB.Stats().OpenConnections,
		InUseConnections: MySQLDB.Stats().InUse,
		IdleConnections: MySQLDB.Stats().Idle,
		WaitCount: MySQLDB.Stats().WaitCount,
		WaitDuration: MySQLDB.Stats().WaitDuration,
		MaxIdleTimeClosed: MySQLDB.Stats().MaxIdleTimeClosed,
	}

	// 获取MongoDB统计信息
	var result bson.M
	err := MongoDB.RunCommand(context.Background(), bson.D{{"serverStatus", 1}}).Decode(&result)
	if err != nil {
		log.Printf("获取MongoDB统计信息失败: %v", err)
	}

	mongoStats := MongoStats{}
	if connectionsRaw, ok := result["connections"]; ok {
		if connections, ok := connectionsRaw.(bson.M); ok {
			if current, ok := connections["current"].(int32); ok {
				mongoStats.ActiveConnections = current
			}
			if available, ok := connections["available"].(int32); ok {
				mongoStats.AvailableConnections = available
			}
			if totalCreated, ok := connections["totalCreated"].(int64); ok {
				mongoStats.TotalCreated = totalCreated
			}
		}
	}

	return &DatabaseStats{
		MySQLStats: mysqlStats,
		MongoStats: mongoStats,
	}
}

// initMySQL 初始化MySQL数据库连接
func initMySQL(cfg config.MySQLConfig) error {
	// 构建DSN连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// 设置连接池参数
	db.SetMaxOpenConns(cfg.MaxOpen)
	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetConnMaxLifetime(time.Duration(cfg.Lifetime) * time.Second)

	// 测试连接
	if err := db.Ping(); err != nil {
		return err
	}

	MySQLDB = db
	log.Println("MySQL数据库连接成功")
	return nil
}

// initMongoDB 初始化MongoDB数据库连接
func initMongoDB(cfg config.MongoConfig) error {
	// 创建MongoDB客户端选项
	clientOptions := options.Client().ApplyURI(cfg.URI)

	// 设置认证信息（如果有）
	if cfg.Username != "" && cfg.Password != "" {
		clientOptions.SetAuth(options.Credential{
			Username: cfg.Username,
			Password: cfg.Password,
		})
	}

	// 设置连接池大小
	clientOptions.SetMaxPoolSize(uint64(cfg.MaxPool))

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接到MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// 测试连接
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	// 获取数据库实例
	MongoDB = client.Database(cfg.Database)
	log.Println("MongoDB数据库连接成功")
	return nil
}

// initRedis 初始化Redis连接
func initRedis(cfg config.RedisConfig) error {
	// 创建Redis客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println("Redis连接成功")
	return nil
}

// CloseDatabase 关闭所有数据库连接
func CloseDatabase() {
	// 关闭MySQL连接
	if MySQLDB != nil {
		MySQLDB.Close()
		log.Println("MySQL连接已关闭")
	}

	// 关闭MongoDB连接
	if MongoDB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		client := MongoDB.Client()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("关闭MongoDB连接失败: %v\n", err)
		} else {
			log.Println("MongoDB连接已关闭")
		}
	}

	// 关闭Redis连接
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Printf("关闭Redis连接失败: %v\n", err)
		} else {
			log.Println("Redis连接已关闭")
		}
	}
}