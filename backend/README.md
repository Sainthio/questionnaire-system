# 问卷系统后端

## 项目概述

这是问卷系统的后端部分，使用Go语言和Gin框架开发，提供RESTful API接口支持前端应用。系统采用分层架构设计，具有良好的可扩展性和可维护性。

## 技术栈

- **开发语言**：Go 1.24.1
- **Web框架**：Gin
- **ORM框架**：GORM
- **数据库**：MySQL 8.0
- **认证**：JWT (JSON Web Token)
- **文档**：Swagger/OpenAPI
- **日志**：Zap
- **配置**：Viper

## 项目结构

```
backend/
├── cmd/                  # 命令行工具
│   └── migration/        # 数据库迁移工具
├── config/               # 配置文件
│   ├── config.go         # 配置结构定义
│   └── config.yaml       # 配置文件
├── database/             # 数据库相关
│   ├── db.go             # 数据库连接
│   └── migration.go      # 数据库迁移
├── handlers/             # HTTP处理器
│   ├── auth_handler.go   # 认证相关处理
│   ├── user_handler.go   # 用户相关处理
│   └── questionnaire_handler.go  # 问卷相关处理
├── middleware/           # 中间件
│   ├── auth.go           # 认证中间件
│   ├── cors.go           # 跨域处理
│   └── logger.go         # 日志中间件
├── models/               # 数据模型
│   ├── base.go           # 基础模型
│   ├── user.go           # 用户模型
│   └── questionnaire.go  # 问卷相关模型
├── scripts/              # 脚本工具
│   └── deploy.sh         # 部署脚本
├── utils/                # 工具函数
│   ├── jwt.go            # JWT工具
│   ├── password.go       # 密码处理
│   └── response.go       # 响应处理
├── main.go               # 主入口文件
├── go.mod                # Go模块文件
├── go.sum                # 依赖校验文件
└── config.yaml           # 主配置文件
```

## 核心功能

### 用户管理
- 用户注册与登录
- JWT认证与授权
- 用户信息管理
- 密码加密与验证

### 问卷管理
- 问卷CRUD操作
- 问题类型支持
- 问卷状态管理
- 权限控制

### 数据处理
- 答卷提交与验证
- 数据统计与分析
- 导出功能
- 缓存优化

### 系统功能
- 日志记录
- 错误处理
- 性能监控
- 安全防护

## 开发环境准备

### 必备工具
1. **Go**: v1.24.1+ [下载地址](https://golang.org/dl/)
2. **MySQL**: v8.0+ [下载地址](https://dev.mysql.com/downloads/mysql/)
3. **Git**: [下载地址](https://git-scm.com/downloads)
4. **IDE推荐**: GoLand 或 VS Code + Go插件

### 依赖安装

本项目使用Go Modules进行依赖管理，主要依赖包括：

```bash
# 核心框架
go get -u github.com/gin-gonic/gin

# ORM框架
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql

# 配置管理
go get -u github.com/spf13/viper

# JWT认证
go get -u github.com/golang-jwt/jwt/v5

# 日志
go get -u go.uber.org/zap

# 其他工具
go get -u github.com/gin-contrib/cors
go get -u golang.org/x/crypto/bcrypt
```

## 安装与运行

### 数据库配置

1. 创建MySQL数据库：

```sql
CREATE DATABASE questionnaire_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 配置数据库连接：

编辑 `config/config.yaml` 文件，修改数据库连接信息：

```yaml
database:
  driver: mysql
  host: localhost
  port: 3306
  username: root
  password: your_password  # 修改为您的实际密码
  dbname: questionnaire_db
  charset: utf8mb4
  parseTime: true
  loc: Local
  max_idle_conns: 10
  max_open_conns: 100
```

### 运行服务

1. 安装依赖：
```bash
go mod tidy
```

2. 运行服务：
```bash
# 开发模式
go run main.go

# 或构建后运行
go build -o questionnaire-server
./questionnaire-server
```

服务默认在 http://localhost:8080 启动。

### 环境变量配置

除了配置文件外，也可以通过环境变量覆盖配置：

```bash
# 设置运行环境
export APP_ENV=development  # development, production, test

# 设置端口
export APP_PORT=8080

# 设置数据库连接
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=questionnaire_db

# JWT密钥
export JWT_SECRET=your_jwt_secret_key
```

## API文档

### 认证相关

#### 用户注册

- **URL**: `/api/user/register`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "username": "test_user",
    "password": "password123",
    "email": "test@example.com",
    "phone": "13800138000"
  }
  ```
- **成功响应** (200 OK): 
  ```json
  {
    "code": 0,
    "message": "注册成功",
    "data": {
      "id": 1,
      "username": "test_user",
      "email": "test@example.com",
      "phone": "13800138000",
      "created_at": "2024-06-01T12:00:00Z"
    }
  }
  ```

#### 用户登录

- **URL**: `/api/user/login`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "username": "test_user",
    "password": "password123"
  }
  ```
- **成功响应** (200 OK): 
  ```json
  {
    "code": 0,
    "message": "登录成功",
    "data": {
      "user_id": 1,
      "username": "test_user",
      "email": "test@example.com",
      "is_admin": false,
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
  }
  ```

### 问卷相关

#### 创建问卷

- **URL**: `/api/questionnaire/create`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "title": "测试问卷",
    "description": "这是一个测试问卷",
    "start_time": "2024-06-01T00:00:00Z",
    "end_time": "2024-06-30T23:59:59Z",
    "questions": [
      {
        "type": "single_choice",
        "title": "您的性别是？",
        "options": ["男", "女", "其他"],
        "required": true
      },
      {
        "type": "multiple_choice",
        "title": "您喜欢的水果有哪些？",
        "options": ["苹果", "香蕉", "橙子", "西瓜"],
        "required": false
      }
    ]
  }
  ```
- **成功响应** (200 OK): 
  ```json
  {
    "code": 0,
    "message": "问卷创建成功",
    "data": {
      "id": 1,
      "title": "测试问卷",
      "description": "这是一个测试问卷",
      "created_by": 1,
      "start_time": "2024-06-01T00:00:00Z",
      "end_time": "2024-06-30T23:59:59Z",
      "status": "draft",
      "created_at": "2024-06-01T12:00:00Z"
    }
  }
  ```

#### 获取问卷列表

- **URL**: `/api/questionnaire/list?page=1&size=10`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}`
- **成功响应** (200 OK): 
  ```json
  {
    "code": 0,
    "message": "获取成功",
    "data": {
      "total": 25,
      "page": 1,
      "size": 10,
      "items": [
        {
          "id": 1,
          "title": "测试问卷",
          "description": "这是一个测试问卷",
          "status": "published",
          "created_by": 1,
          "start_time": "2024-06-01T00:00:00Z",
          "end_time": "2024-06-30T23:59:59Z",
          "submission_count": 15,
          "created_at": "2024-06-01T12:00:00Z"
        },
        // ...更多问卷
      ]
    }
  }
  ```

#### 获取问卷详情

- **URL**: `/api/questionnaire/detail?id=1`
- **方法**: `GET`
- **请求头**: `Authorization: Bearer {token}` (可选，公开问卷无需认证)
- **成功响应** (200 OK): 
  ```json
  {
    "code": 0,
    "message": "获取成功",
    "data": {
      "id": 1,
      "title": "测试问卷",
      "description": "这是一个测试问卷",
      "status": "published",
      "created_by": 1,
      "creator_name": "test_user",
      "start_time": "2024-06-01T00:00:00Z",
      "end_time": "2024-06-30T23:59:59Z",
      "questions": [
        {
          "id": 1,
          "type": "single_choice",
          "title": "您的性别是？",
          "options": ["男", "女", "其他"],
          "required": true
        },
        {
          "id": 2,
          "type": "multiple_choice",
          "title": "您喜欢的水果有哪些？",
          "options": ["苹果", "香蕉", "橙子", "西瓜"],
          "required": false
        }
      ],
      "created_at": "2024-06-01T12:00:00Z",
      "updated_at": "2024-06-01T12:00:00Z"
    }
  }
  ```

#### 提交问卷答案

- **URL**: `/api/questionnaire/submit`
- **方法**: `POST`
- **请求头**: `Authorization: Bearer {token}` (可选，匿名提交无需认证)
- **请求体**:
  ```json
  {
    "questionnaire_id": 1,
    "answers": [
      {
        "question_id": 1,
        "content": "男"
      },
      {
        "question_id": 2,
        "content": ["苹果", "西瓜"]
      }
    ]
  }
  ```
- **成功响应** (200 OK): 
  ```json
  {
    "code": 0,
    "message": "问卷提交成功",
    "data": {
      "submission_id": 1,
      "submission_time": "2024-06-10T15:30:45Z"
    }
  }
  ```

## 数据库设计

系统使用以下主要数据表：

### users 表
```sql
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(100) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `is_admin` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### questionnaires 表
```sql
CREATE TABLE `questionnaires` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `description` text,
  `created_by` bigint(20) NOT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'draft',
  `start_time` datetime DEFAULT NULL,
  `end_time` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_questionnaires_user` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### questions 表
```sql
CREATE TABLE `questions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `questionnaire_id` bigint(20) NOT NULL,
  `type` varchar(50) NOT NULL,
  `title` varchar(255) NOT NULL,
  `options` text,
  `required` tinyint(1) NOT NULL DEFAULT '0',
  `sort_order` int(11) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_questionnaire_id` (`questionnaire_id`),
  CONSTRAINT `fk_questions_questionnaire` FOREIGN KEY (`questionnaire_id`) REFERENCES `questionnaires` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### submissions 表
```sql
CREATE TABLE `submissions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `questionnaire_id` bigint(20) NOT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `ip_address` varchar(50) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_questionnaire_id` (`questionnaire_id`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `fk_submissions_questionnaire` FOREIGN KEY (`questionnaire_id`) REFERENCES `questionnaires` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_submissions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### answers 表
```sql
CREATE TABLE `answers` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `submission_id` bigint(20) NOT NULL,
  `question_id` bigint(20) NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_submission_id` (`submission_id`),
  KEY `idx_question_id` (`question_id`),
  CONSTRAINT `fk_answers_submission` FOREIGN KEY (`submission_id`) REFERENCES `submissions` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_answers_question` FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 部署指南

### 构建二进制文件

```bash
# 适用于Linux
GOOS=linux GOARCH=amd64 go build -o questionnaire-server-linux main.go

# 适用于Windows
GOOS=windows GOARCH=amd64 go build -o questionnaire-server.exe main.go

# 适用于MacOS
GOOS=darwin GOARCH=amd64 go build -o questionnaire-server-mac main.go
```

### Docker部署

1. 创建Dockerfile:

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o questionnaire-server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/questionnaire-server .
COPY --from=builder /app/config.yaml .
EXPOSE 8080
CMD ["./questionnaire-server"]
```

2. 构建和运行Docker容器:

```bash
# 构建镜像
docker build -t questionnaire-server:latest .

# 运行容器
docker run -d -p 8080:8080 --name questionnaire-api \
  -e DB_HOST=db_host \
  -e DB_USER=root \
  -e DB_PASSWORD=password \
  -e DB_NAME=questionnaire_db \
  questionnaire-server:latest
```

### 使用Docker Compose

创建docker-compose.yml文件:

```yaml
version: '3'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_USER=root
      - DB_PASSWORD=password
      - DB_NAME=questionnaire_db
    depends_on:
      - db
    restart: always

  db:
    image: mysql:8
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=password
      - MYSQL_DATABASE=questionnaire_db
    volumes:
      - mysql_data:/var/lib/mysql
    restart: always

volumes:
  mysql_data:
```

启动服务:

```bash
docker-compose up -d
```

## 性能优化

系统已实现以下性能优化措施：

1. **数据库优化**
   - 合理的索引设计
   - 连接池配置
   - 查询优化

2. **缓存策略**
   - 热点数据缓存
   - 统计结果缓存

3. **并发处理**
   - Goroutine池
   - 异步处理

4. **API优化**
   - 分页查询
   - 字段过滤
   - 压缩响应

## 常见问题解决

### 数据库连接问题

**问题**: 无法连接到MySQL数据库
**解决方法**:
- 检查数据库服务是否运行
- 验证连接信息是否正确
- 确认数据库用户权限
- 检查网络连接和防火墙设置

### 权限验证问题

**问题**: JWT验证失败
**解决方法**:
- 检查JWT密钥配置
- 验证token格式是否正确
- 确认token是否过期
- 查看日志中的详细错误信息

## 贡献指南

欢迎为项目做出贡献！请遵循以下步骤：

1. Fork 项目仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 LICENSE 文件 