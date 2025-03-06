package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

// TestCachePerformance 测试缓存性能
func TestCachePerformance(t *testing.T) {
	// 初始化Redis客户端
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// 测试写入性能
	startWrite := time.Now()
	for i := 0; i < 1000; i++ {
		err := rdb.Set(ctx, fmt.Sprintf("test_key_%d", i), "test_value", time.Hour).Err()
		assert.NoError(t, err)
	}
	writeDuration := time.Since(startWrite)

	// 测试读取性能
	startRead := time.Now()
	for i := 0; i < 1000; i++ {
		_, err := rdb.Get(ctx, fmt.Sprintf("test_key_%d", i)).Result()
		assert.NoError(t, err)
	}
	readDuration := time.Since(startRead)

	// 验证性能指标
	assert.Less(t, writeDuration.Milliseconds(), int64(2000)) // 写入时间小于2秒
	assert.Less(t, readDuration.Milliseconds(), int64(1000))  // 读取时间小于1秒
}

// TestDatabaseQueryPerformance 测试数据库查询性能
func TestDatabaseQueryPerformance(t *testing.T) {
	// 初始化MongoDB客户端
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err)
	defer client.Disconnect(ctx)

	coll := client.Database("test_db").Collection("test_collection")

	// 测试批量插入性能
	docs := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		docs[i] = bson.M{"test_field": i}
	}

	startInsert := time.Now()
	_, err = coll.InsertMany(ctx, docs)
	assert.NoError(t, err)
	insertDuration := time.Since(startInsert)

	// 测试查询性能
	startQuery := time.Now()
	cursor, err := coll.Find(ctx, bson.M{"test_field": bson.M{"$gte": 0}})
	assert.NoError(t, err)
	var results []bson.M
	err = cursor.All(ctx, &results)
	assert.NoError(t, err)
	queryDuration := time.Since(startQuery)

	// 验证性能指标
	assert.Less(t, insertDuration.Milliseconds(), int64(3000)) // 插入时间小于3秒
	assert.Less(t, queryDuration.Milliseconds(), int64(1000))  // 查询时间小于1秒
	assert.Equal(t, 1000, len(results))
}

// TestIndexPerformance 测试索引性能
func TestIndexPerformance(t *testing.T) {
	// 初始化MongoDB客户端
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err)
	defer client.Disconnect(ctx)

	coll := client.Database("test_db").Collection("test_index")

	// 创建索引
	_, err = coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"test_field", 1}},
	})
	assert.NoError(t, err)

	// 插入测试数据
	docs := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		docs[i] = bson.M{"test_field": i, "data": fmt.Sprintf("data_%d", i)}
	}
	_, err = coll.InsertMany(ctx, docs)
	assert.NoError(t, err)

	// 测试带索引的查询性能
	startQuery := time.Now()
	cursor, err := coll.Find(ctx, bson.M{"test_field": bson.M{"$gt": 500}})
	assert.NoError(t, err)
	var results []bson.M
	err = cursor.All(ctx, &results)
	assert.NoError(t, err)
	indexedQueryDuration := time.Since(startQuery)

	// 测试非索引字段的查询性能
	startQuery = time.Now()
	cursor, err = coll.Find(ctx, bson.M{"data": bson.M{"$regex": "^data_5"}})
	assert.NoError(t, err)
	results = nil
	err = cursor.All(ctx, &results)
	assert.NoError(t, err)
	nonIndexedQueryDuration := time.Since(startQuery)

	// 验证索引查询性能优势
	assert.Less(t, indexedQueryDuration, nonIndexedQueryDuration/2)
}