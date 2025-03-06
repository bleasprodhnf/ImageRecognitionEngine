-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS image_recognition;

-- 使用数据库
USE image_recognition;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 创建图像识别记录表
CREATE TABLE IF NOT EXISTS recognition_records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    image_path VARCHAR(255) NOT NULL,
    result TEXT,
    model_used VARCHAR(50),
    confidence FLOAT,
    processing_time FLOAT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 创建模型表
CREATE TABLE IF NOT EXISTS models (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    version VARCHAR(20) NOT NULL,
    description TEXT,
    file_path VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 创建统计表
CREATE TABLE IF NOT EXISTS statistics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    date DATE NOT NULL,
    total_requests INT DEFAULT 0,
    successful_requests INT DEFAULT 0,
    failed_requests INT DEFAULT 0,
    avg_processing_time FLOAT,
    UNIQUE KEY (date)
);

-- 插入默认管理员用户（密码：password）
INSERT INTO users (username, password, email, role) VALUES
('admin', '$2a$10$kz7O8DhQg0RpRRK9.QVpJeQOTXUQD7dZwILGPXq8jYfHrpYzGdQiS', 'admin@example.com', 'admin')
ON DUPLICATE KEY UPDATE username=username;
