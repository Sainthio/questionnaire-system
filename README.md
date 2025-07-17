# 移动端问卷系统

这是一个功能对标"问卷星"的移动端问卷系统，支持用户在手机端便捷完成投票、调查问卷、考评等交互操作，同时涵盖配套的管理功能（如问卷创建、数据统计、权限管理等）。

## 技术栈

- 前端：Vue3 + Vant + ECharts
- 后端：Go + Gin
- 数据库：MySQL 8

## 项目结构

```
project_1/
  ├── backend/           # Go后端
  │   ├── config/        # 配置
  │   ├── database/      # 数据库连接
  │   ├── handlers/      # 请求处理器
  │   ├── models/        # 数据模型
  │   └── main.go        # 主入口
  └── frontend/          # Vue3前端
      └── questionnaire-app/ # Vue项目
          ├── src/       # 源代码
          │   ├── components/ # 组件
          │   ├── views/      # 页面
          │   ├── stores/     # 状态管理
          │   └── services/   # API服务
          └── public/    # 静态资源
```

## 开发环境准备

1. **NodeJS**: https://nodejs.org/en/download/
2. **Go**: https://golang.org/dl/
3. **MySQL 8**: https://dev.mysql.com/downloads/mysql/
4. **Navicat**: 数据库管理工具
5. **VS Code**: https://code.visualstudio.com/

## 安装与启动

### 数据库配置

1. 创建MySQL数据库：
```sql
CREATE DATABASE questionnaire_db;
```

2. 配置数据库连接：
修改 `backend/config/config.go` 中的数据库连接信息，包括主机、端口、用户名和密码。

### 后端启动

1. 安装依赖：
```bash
cd backend
go mod tidy
```

2. 运行后端服务：
```bash
# Windows系统（推荐以管理员权限运行）
run_server.bat

# Linux/Mac系统
go run main.go
```

### 前端启动

1. 安装依赖：
```bash
cd frontend/questionnaire-app
npm install
```

2. 运行开发服务器：
```bash
npm run dev
```

### 一键启动（Windows）

项目根目录提供了一键启动脚本，可同时启动前后端服务：
```bash
# 以管理员权限运行
start_server.bat
```

## 网络配置

系统支持在局域网内访问，需要注意以下几点：

1. 后端服务监听所有网络接口（0.0.0.0:8080）
2. 前端服务监听所有网络接口（0.0.0.0:5173）
3. Windows系统需要配置防火墙规则（启动脚本会自动处理）
4. 访问时可使用本机IP地址替代localhost

## 功能特性

- 用户注册与登录
- 问卷创建与编辑
- 多种题型支持（单选、多选、填空、评分等）
- 问卷数据统计与分析
- 移动端友好的UI设计

## 统计功能

系统提供了丰富的统计分析功能：

1. **图表展示**：使用ECharts展示问卷结果，包括：
   - 单选题：饼图
   - 多选题：柱状图
   - 评分题：柱状图

2. **文字描述**：每个图表附带详细的文字统计描述，包括：
   - 回答总数
   - 选项占比
   - 平均分（评分题）
   - 最常见选项

3. **导出功能**：支持将问卷结果导出为CSV格式

## 常见问题解决

### 1. 网络连接问题

- **症状**：后端启动时卡住，显示Windows防火墙提示
- **解决方法**：
  - 以管理员权限运行启动脚本
  - 允许应用通过防火墙
  - 使用`0.0.0.0`替代`localhost`监听所有网络接口

### 2. 权限问题

- **症状**：无法查看问卷统计结果，显示权限不足
- **解决方法**：
  - 确保已登录
  - 检查是否为问卷创建者或管理员
  - 后端已更新权限检查逻辑

### 3. ECharts图表问题

- **症状**：图表不显示或显示错误
- **解决方法**：
  - 已添加GridComponent组件
  - 优化了图表渲染逻辑
  - 增加了错误处理和调试信息
  - 添加了窗口大小变化时自动调整图表大小

### 4. 统计数据问题

- **症状**：首页统计数据显示不正确
- **解决方法**：
  - 优化了统计数据获取逻辑
  - 添加了备用数据获取方法
  - 增加了日志输出以便调试

## API文档

### 用户相关

- `POST /api/user/register`: 用户注册
- `POST /api/user/login`: 用户登录
- `POST /api/user/reset-password`: 重置密码

### 问卷相关

- `POST /api/questionnaire/create`: 创建问卷
- `GET /api/questionnaire/list`: 获取问卷列表
- `GET /api/questionnaire/detail`: 获取问卷详情
- `POST /api/questionnaire/submit`: 提交问卷答案
- `PUT /api/questionnaire/update`: 更新问卷
- `PUT /api/questionnaire/update-status`: 更新问卷状态
- `DELETE /api/questionnaire/delete`: 删除问卷
- `GET /api/questionnaire/results`: 获取问卷结果
- `GET /api/questionnaire/check-submission`: 检查提交状态
- `GET /api/questionnaire/stats`: 获取系统统计数据 