# 图像识别引擎系统 - 数据库设计文档

## 1. 数据库架构

系统采用混合数据库架构，结合了关系型数据库(MySQL)和文档型数据库(MongoDB)的优势：

- **MySQL**: 用于存储结构化数据，如用户信息、权限配置等
- **MongoDB**: 用于存储非结构化数据，如模型文件、识别结果等

### 1.1 MySQL数据库

#### 用户认证相关表

##### admin_users（管理员用户表）
```sql
CREATE TABLE admin_users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    real_name VARCHAR(50),
    email VARCHAR(100),
    phone VARCHAR(20),
    status TINYINT DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    last_login_time DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_status (status)
);
```

##### admin_user_roles（管理员角色关联表）
```sql
CREATE TABLE admin_user_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    admin_user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_admin_role (admin_user_id, role_id),
    INDEX idx_admin_user_id (admin_user_id),
    INDEX idx_role_id (role_id),
    FOREIGN KEY (admin_user_id) REFERENCES admin_users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);
```

##### role_permissions（角色权限关联表）
```sql
CREATE TABLE role_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_role_permission (role_id, permission_id),
    INDEX idx_role_id (role_id),
    INDEX idx_permission_id (permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);
```

##### roles（角色表）
```sql
CREATE TABLE roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(200),
    status TINYINT DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name)
);
```

##### permissions（权限表）
```sql
CREATE TABLE permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    code VARCHAR(100) NOT NULL UNIQUE COMMENT '权限编码',
    name VARCHAR(50) NOT NULL COMMENT '权限名称',
    module VARCHAR(50) NOT NULL COMMENT '所属模块',
    description VARCHAR(200),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_code (code),
    INDEX idx_module (module)
);
```

#### 客户相关表

##### customers（客户表）
```sql
CREATE TABLE customers (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    company_name VARCHAR(100),
    contact_name VARCHAR(50),
    email VARCHAR(100),
    phone VARCHAR(20),
    status TINYINT DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    package_id BIGINT COMMENT '套餐ID',
    api_key VARCHAR(64) NOT NULL UNIQUE COMMENT 'API密钥',
    api_secret VARCHAR(128) NOT NULL COMMENT 'API密钥密文',
    last_login_time DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_status (status),
    INDEX idx_api_key (api_key)
);
```

##### packages（套餐表）
```sql
CREATE TABLE packages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    api_calls_limit INT COMMENT 'API调用次数限制',
    storage_limit INT COMMENT '存储空间限制(GB)',
    concurrency_limit INT COMMENT '并发请求限制',
    status TINYINT DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name)
);
```

##### customer_package_records（客户套餐变更记录表）
```sql
CREATE TABLE customer_package_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    customer_id BIGINT NOT NULL,
    package_id BIGINT NOT NULL,
    change_type TINYINT NOT NULL COMMENT '变更类型：1-新购，2-续费，3-升级，4-降级',
    original_package_id BIGINT COMMENT '原套餐ID',
    price DECIMAL(10,2) NOT NULL COMMENT '实付金额',
    start_time DATETIME NOT NULL COMMENT '生效时间',
    end_time DATETIME NOT NULL COMMENT '到期时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_customer_id (customer_id),
    INDEX idx_package_id (package_id),
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (package_id) REFERENCES packages(id)
);
```

### 1.2 MongoDB数据库

#### 模型相关集合

##### model_versions
```javascript
{
    _id: ObjectId,
    version: String,          // 版本号
    description: String,      // 版本描述
    status: String,          // 状态(development/testing/production)
    file_path: String,       // 模型文件路径
    accuracy: Number,        // 准确率
    parameters: {            // 模型参数
        batchSize: Number,
        learningRate: Number,
        epochs: Number
    },
    metrics: {               // 性能指标
        latency: Number,     // 延迟
        throughput: Number   // 吞吐量
    },
    created_at: ISODate,
    updated_at: ISODate
}
```

##### recognition_records
```javascript
{
    _id: ObjectId,
    customer_id: Long,       // 客户ID
    model_version: String,   // 模型版本
    image_url: String,       // 图像URL
    image_storage_type: String, // 存储类型(local/s3/oss)
    result: {               // 识别结果
        labels: [String],    // 标签列表
        confidence: Number   // 置信度
    },
    processing_time: Number, // 处理时间(ms)
    status: String,         // 状态(success/failed)
    error_message: String,  // 错误信息
    created_at: ISODate
}
```

## 2. 索引设计

### 2.1 MySQL索引

- **admin_users表**
  - 主键索引：id
  - 唯一索引：username
  - 普通索引：status

- **roles表**
  - 主键索引：id
  - 唯一索引：name

- **permissions表**
  - 主键索引：id
  - 唯一索引：code
  - 普通索引：module

- **customers表**
  - 主键索引：id
  - 唯一索引：username, api_key
  - 普通索引：status, package_id

- **packages表**
  - 主键索引：id
  - 唯一索引：name

### 2.2 MongoDB索引

- **model_versions集合**
  ```javascript
  db.model_versions.createIndex({ "version": 1 }, { unique: true })
  db.model_versions.createIndex({ "status": 1 })
  db.model_versions.createIndex({ "created_at": -1 })
  ```

- **recognition_records集合**
  ```javascript
  db.recognition_records.createIndex({ "customer_id": 1, "created_at": -1 })
  db.recognition_records.createIndex({ "model_version": 1 })
  db.recognition_records.createIndex({ "status": 1 })
  ```

## 3. 数据库优化策略

### 3.1 MySQL优化

1. **连接池配置**
   - 最小连接数：10
   - 最大连接数：100
   - 连接超时：30000ms

2. **查询优化**
   - 合理使用索引
   - 避免使用SELECT *
   - 使用EXPLAIN分析查询计划
   - 优化JOIN操作
   - 合理设置WHERE条件

3. **表结构优化**
   - 选择合适的字段类型
   - 适当冗余设计
   - 大字段拆分
   - 历史数据归档

4. **事务优化**
   - 控制事务大小
   - 合理设置隔离级别
   - 避免长事务

### 3.2 MongoDB优化

1. **索引优化**
   - 创建复合索引
   - 避免过多索引
   - 定期分析索引使用情况

2. **查询优化**
   - 使用投影操作
   - 避免大规模的skip操作
   - 合理使用聚合管道

3. **文档结构优化**
   - 避免过深的嵌套
   - 适当反范式化
   - 控制文档大小

4. **分片策略**
   - 选择合适的片键
   - 均衡数据分布
   - 定期检查分片状态

## 4. 性能监控

### 4.1 MySQL监控指标

1. **基础指标**
   - QPS (Queries Per Second)
   - TPS (Transactions Per Second)
   - 连接数
   - 缓存命中率
   - API调用频率
   - 存储空间使用率

2. **资源使用**
   - CPU使用率
   - 内存使用
   - 磁盘IO
   - 网络带宽

3. **慢查询**
   - 慢查询阈值：500ms
   - 慢查询日志分析
   - 定期优化

### 4.2 MongoDB监控指标

1. **操作指标**
   - 读写操作数
   - 平均响应时间
   - 并发连接数
   - 游标数量

2. **存储指标**
   - 数据大小
   - 索引大小
   - 文档数量
   - 碎片率

3. **性能指标**
   - 页面错误率
   - 锁等待时间
   - 网络延迟
   - 复制延迟

## 5. 数据备份策略

### 5.1 MySQL备份

1. **全量备份**
   - 周期：每天凌晨
   - 工具：mysqldump
   - 保留时间：30天
   - 压缩存储
   - 存储位置：独立备份服务器
   - 访问权限：仅DBA和系统管理员
   - 加密方案：AES-256加密

2. **增量备份**
   - 周期：每4小时
   - 基于binlog
   - 保留时间：7天

3. **备份验证**
   - 定期恢复测试
   - 完整性检查
   - 性能测试

### 5.2 MongoDB备份

1. **全量备份**
   - 周期：每天凌晨
   - 工具：mongodump
   - 保留时间：30天
   - 压缩存储

2. **增量备份**
   - 周期：每6小时
   - 基于oplog
   - 保留时间：7天

3. **备份验证**
   - 定期恢复测试
   - 数据一致性检查
   - 性能验证

## 6. 灾难恢复

### 6.1 故障类型

1. **硬件故障**
   - 磁盘损坏
   - 服务器宕机
   - 网络中断

2. **软件故障**
   - 数据库崩溃
   - 数据损坏
   - 人为误操作

### 6.2 恢复策略

1. **即时恢复**
   - 主从切换
   - 读写分离
   - 负载均衡

2. **数据恢复**
   - 全量恢复
   - 增量恢复
   - 时间点恢复

3. **应急预案**
   - 故障报警
   - 自动切换
   - 手动干预

### 6.3 RTO和RPO

- RTO (Recovery Time Objective): < 30分钟
- RPO (Recovery Point Objective): < 5分钟