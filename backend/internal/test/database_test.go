package test

import (
	"database/sql"
	"fmt"
	"testing"
	
	_ "github.com/go-sql-driver/mysql"
)

// TestDatabaseConnection 测试数据库连接
func TestDatabaseConnection(t *testing.T) {
	config, err := LoadTestConfig()
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	// 尝试连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// 测试基本的数据库操作
	// 创建测试表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS test_table (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(50) NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// 插入测试数据
	result, err := db.Exec("INSERT INTO test_table (name) VALUES (?)", "test_record")
	if err != nil {
		t.Fatalf("Failed to insert test record: %v", err)
	}

	// 验证插入结果
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert ID: %v", err)
	}

	// 查询测试数据
	var name string
	err = db.QueryRow("SELECT name FROM test_table WHERE id = ?", id).Scan(&name)
	if err != nil {
		t.Fatalf("Failed to query test record: %v", err)
	}

	if name != "test_record" {
		t.Errorf("Expected name to be 'test_record', got '%s'", name)
	}

	// 清理测试数据
	_, err = db.Exec("DROP TABLE test_table")
	if err != nil {
		t.Fatalf("Failed to drop test table: %v", err)
	}
}