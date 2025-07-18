# ques问卷系统

## 项目简介

这是一坨功能对标"问卷星"的问卷系统，支持用户便捷完成投票、调查问卷、考评等交互操作，同时提供完整的管理功能（问卷创建、数据统计、权限管理等）。系统采用前后端分离架构，具有良好的用户体验和扩展性。

## 技术栈

### 前端
- **核心框架**：Vue 3
- **UI组件库**：Vant UI
- **数据可视化**：ECharts
- **状态管理**：Pinia
- **路由管理**：Vue Router
- **HTTP客户端**：Axios

### 后端
- **开发语言**：Go 1.24.1
- **Web框架**：Gin
- **ORM框架**：GORM
- **数据库**：MySQL 8

## 系统架构

```
questionnaire/
  ├── backend/           # Go后端服务
  │   ├── cmd/           # 命令行工具
  │   ├── config/        # 配置文件
  │   ├── database/      # 数据库连接
  │   ├── handlers/      # HTTP请求处理
  │   ├── middleware/    # 中间件
  │   ├── models/        # 数据模型
  │   ├── scripts/       # 脚本工具
  │   └── main.go        # 主入口文件
  │
  └── frontend/          # Vue3前端应用
      ├── questionnaire-app/  # Vue项目
      │   ├── public/    # 静态资源
      │   └── src/       # 源代码
      │       ├── assets/     # 资源文件
      │       ├── components/ # 通用组件
      │       ├── router/     # 路由配置
      │       ├── stores/     # 状态管理
      │       ├── views/      # 页面视图
      │       ├── utils/      # 工具函数
      │       └── services/   # API服务
      └── dist/          # 构建输出目录
```

## 核心功能

### 用户功能
- 用户注册与登录
- 个人资料管理
- 问卷浏览与填写
- 历史记录查看

### 问卷管理
- 问卷创建与编辑
- 多种题型支持
  - 单选题
  - 多选题
  - 填空题
  - 评分题
  - 矩阵题
- 问卷发布与关闭
- 问卷分享（链接、二维码）

### 数据分析
- 实时统计数据
- 多维度图表展示
- 数据导出（CSV格式）
- 答卷详情查看

### 系统管理
- 用户权限管理
- 系统监控与日志
- 数据备份与恢复

## 开发环境准备

### 必备工具
1. **NodeJS**: v16.0.0+ https://nodejs.org/en/download/
2. **Go**: v1.24.1+ https://golang.org/dl/
3. **MySQL**: v8.0+ https://dev.mysql.com/downloads/mysql/
4. **Git**: https://git-scm.com/downloads
5. **IDE推荐**: VS Code https://code.visualstudio.com/

### 数据库配置

1. 创建MySQL数据库：
```sql
CREATE DATABASE questionnaire_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 配置数据库连接：
修改 `backend/config/config.yaml` 或 `backend/config.yaml` 中的数据库连接信息。

## 安装与启动

### 后端服务

1. 安装依赖：
```bash
cd backend
go mod tidy
```

2. 运行后端服务：
```bash
# 直接运行
go run main.go

# 或先构建再运行
go build -o questionnaire-server
./questionnaire-server
```

### 前端应用

1. 安装依赖：
```bash
cd frontend/questionnaire-app
npm install
```

2. 开发模式运行：
```bash
npm run dev
```

3. 构建生产版本：
```bash
npm run build
```

## 网络配置

### 本地开发
- 前端默认运行在: http://localhost:5173
- 后端默认运行在: http://localhost:8080

### 局域网访问
系统支持在局域网内访问，需要注意以下几点：

1. 后端服务配置监听所有网络接口：
```go
r.Run("0.0.0.0:8080")
```

2. 前端开发服务器配置：
```js
// vite.config.js
export default defineConfig({
  server: {
    host: '0.0.0.0',
    port: 5173
  }
})
```

3. 防火墙配置：
   - Windows系统需要允许应用通过防火墙
   - Linux系统可能需要配置iptables规则

## API文档

### 用户相关

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/user/register` | POST | 用户注册 |
| `/api/user/login` | POST | 用户登录 |
| `/api/user/reset-password` | POST | 重置密码 |

### 问卷相关

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/questionnaire/create` | POST | 创建问卷 |
| `/api/questionnaire/list` | GET | 获取问卷列表 |
| `/api/questionnaire/detail` | GET | 获取问卷详情 |
| `/api/questionnaire/submit` | POST | 提交问卷答案 |
| `/api/questionnaire/update` | PUT | 更新问卷 |
| `/api/questionnaire/update-status` | PUT | 更新问卷状态 |
| `/api/questionnaire/delete` | DELETE | 删除问卷 |
| `/api/questionnaire/results` | GET | 获取问卷结果 |
| `/api/questionnaire/stats` | GET | 获取系统统计数据 |

## 统计功能

系统提供了丰富的统计分析功能：

1. **图表展示**：使用ECharts实现多种可视化图表
   - 单选题：饼图、环形图
   - 多选题：柱状图、条形图
   - 评分题：柱状图、雷达图
   - 填空题：词云图

2. **数据分析**：
   - 回答总数与完成率
   - 选项分布与占比
   - 平均分与标准差
   - 答题时长分析
   - 用户画像分析

3. **导出功能**：
   - CSV格式导出
   - Excel格式导出
   - 图表图片导出

## 常见问题解决

### 网络连接问题

**问题**：无法连接到服务器或前端无法访问API
**解决方法**：
- 检查服务是否正常运行（可通过日志查看）
- 确认防火墙设置允许应用通信
- 验证API地址配置是否正确
- 使用`0.0.0.0`替代`localhost`监听所有网络接口

### 权限问题

**问题**：无法访问特定功能或资源
**解决方法**：
- 确保已正确登录系统
- 检查用户权限设置
- 验证token是否有效或已过期
- 查看后端日志中的权限检查记录

### 数据统计问题

**问题**：统计图表不显示或数据异常
**解决方法**：
- 确保问卷有足够的提交数据
- 检查ECharts配置是否正确
- 验证数据格式是否符合图表要求
- 添加适当的错误处理和空数据展示

## 贡献指南

欢迎贡献代码或提出建议！请遵循以下步骤：

1. Fork 项目仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 LICENSE 文件 
