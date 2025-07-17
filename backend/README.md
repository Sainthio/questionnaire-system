# 问卷系统后端

这是问卷系统的后端部分，使用Go语言开发。

## 技术栈

- Go 1.24.1
- GORM (ORM库)
- MySQL 8

## 项目结构

```
backend/
  ├── config/          # 配置文件
  │   └── config.go    # 应用配置
  ├── database/        # 数据库相关
  │   └── db.go        # 数据库连接
  ├── handlers/        # HTTP处理器
  │   ├── questionnaire_handler.go  # 问卷相关处理
  │   └── user_handler.go           # 用户相关处理
  ├── models/          # 数据模型
  │   └── questionnaire.go          # 问卷相关模型
  ├── go.mod           # Go模块文件
  └── main.go          # 主入口文件
```

## 安装与运行

### 前置条件

1. 安装Go 1.24或更高版本
2. 安装MySQL 8
3. 创建数据库

```sql
CREATE DATABASE questionnaire_db;
```

### 配置数据库

编辑 `config/config.go` 文件，修改数据库连接信息：

```go
DatabaseConfig{
    Host:     "localhost",
    Port:     "3306",
    User:     "root",
    Password: "your_password", // 修改为您的实际密码
    DBName:   "questionnaire_db",
}
```

### 安装依赖

```bash
go mod tidy
```

### 运行服务

```bash
go run main.go
```

服务将在 http://localhost:8080 启动。

## API接口

### 用户相关

#### 注册用户

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
- **成功响应**: 
  ```json
  {
    "id": 1,
    "username": "test_user",
    "email": "test@example.com",
    "phone": "13800138000",
    "is_admin": false,
    "created_at": "2024-06-01T12:00:00Z",
    "updated_at": "2024-06-01T12:00:00Z"
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
- **成功响应**: 
  ```json
  {
    "user_id": 1,
    "username": "test_user",
    "email": "test@example.com",
    "is_admin": false,
    "token": "sample_token_test_user"
  }
  ```

### 问卷相关

#### 创建问卷

- **URL**: `/api/questionnaire/create`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "title": "测试问卷",
    "description": "这是一个测试问卷",
    "created_by": 1,
    "start_time": "2024-06-01T00:00:00Z",
    "end_time": "2024-06-30T23:59:59Z"
  }
  ```
- **成功响应**: 返回创建的问卷信息

#### 获取问卷详情

- **URL**: `/api/questionnaire/get?id=1`
- **方法**: `GET`
- **成功响应**: 返回问卷详情及问题列表

#### 获取问卷列表

- **URL**: `/api/questionnaire/list?page=1`
- **方法**: `GET`
- **成功响应**: 返回分页的问卷列表

#### 提交问卷答案

- **URL**: `/api/questionnaire/submit`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "questionnaire_id": 1,
    "user_id": 1,
    "answers": [
      {
        "question_id": 1,
        "content": "选项A"
      },
      {
        "question_id": 2,
        "content": "这是填空题答案"
      }
    ]
  }
  ```
- **成功响应**: 
  ```json
  {
    "message": "问卷提交成功",
    "submission_id": 1
  }
  ```

## 数据库结构

系统使用以下数据表：

1. **users**: 用户信息
2. **questionnaires**: 问卷信息
3. **questions**: 问题信息
4. **answers**: 答案记录
5. **submissions**: 提交记录

数据表结构在 `models/questionnaire.go` 文件中定义，使用GORM自动迁移功能创建。 