# 问卷系统前端

## 项目概述

这是问卷系统的前端部分，基于Vue 3开发，采用Vant UI组件库打造移动端友好的用户界面。本项目使用Vite作为构建工具，提供了高效的开发体验和快速的热重载功能。

## 技术栈

- **核心框架**：Vue 3 + Composition API
- **构建工具**：Vite 4
- **UI组件库**：Vant 4
- **状态管理**：Pinia
- **路由管理**：Vue Router 4
- **HTTP客户端**：Axios
- **数据可视化**：ECharts 5
- **CSS预处理器**：Sass/SCSS

## 快速开始

### 环境要求

- Node.js 16.0.0+
- npm 7.0.0+ 或 yarn 1.22.0+

### 安装步骤

1. **克隆项目**（如果尚未克隆）
```bash
git clone <仓库地址>
cd questionnaire/frontend
```

2. **安装依赖**
```bash
cd questionnaire-app
npm install
# 或使用 yarn
yarn install
```

3. **启动开发服务器**
```bash
npm run dev
# 或使用 yarn
yarn dev
```

4. **构建生产版本**
```bash
npm run build
# 或使用 yarn
yarn build
```

## 项目结构

```
questionnaire-app/
├── public/                 # 静态资源
│   ├── favicon.ico         # 网站图标
│   └── images/             # 图片资源
├── src/                    # 源代码
│   ├── assets/             # 资源文件
│   │   ├── styles/         # 全局样式
│   │   └── images/         # 图片资源
│   ├── components/         # 通用组件
│   │   ├── common/         # 公共组件
│   │   ├── form/           # 表单组件
│   │   ├── layout/         # 布局组件
│   │   └── question/       # 问题组件
│   ├── router/             # 路由配置
│   │   └── index.js        # 路由定义
│   ├── stores/             # Pinia状态管理
│   │   ├── user.js         # 用户状态
│   │   └── questionnaire.js # 问卷状态
│   ├── utils/              # 工具函数
│   │   ├── request.js      # Axios封装
│   │   ├── auth.js         # 认证相关
│   │   └── validator.js    # 表单验证
│   ├── views/              # 页面视图
│   │   ├── auth/           # 认证相关页面
│   │   ├── home/           # 首页
│   │   ├── questionnaire/  # 问卷相关页面
│   │   └── user/           # 用户中心
│   ├── App.vue             # 根组件
│   └── main.js             # 入口文件
├── .env                    # 环境变量
├── .env.development        # 开发环境变量
├── .env.production         # 生产环境变量
├── index.html              # HTML模板
├── package.json            # 项目配置
└── vite.config.js          # Vite配置
```

## 核心功能模块

### 用户认证模块
- 登录/注册页面
- 密码重置功能
- 认证状态管理
- Token存储与刷新

### 问卷管理模块
- 问卷列表展示
- 问卷搜索与筛选
- 问卷创建与编辑
- 问卷预览功能

### 问卷填写模块
- 多种题型渲染
- 表单验证
- 答案保存与提交
- 填写进度保存

### 统计分析模块
- 数据可视化图表
- 答案详情查看
- 数据导出功能
- 实时统计更新

## Vant UI 组件使用

Vant是一个轻量、可靠的移动端组件库，本项目中广泛使用了Vant组件。

### 全局引入

在`main.js`中配置：

```js
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import vant from 'vant'
import 'vant/lib/index.css'

const app = createApp(App)
app.use(router)
app.use(createPinia())
app.use(vant)
app.mount('#app')
```

### 按需引入（推荐）

为了减小打包体积，推荐使用按需引入：

```js
// vite.config.js
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import Components from 'unplugin-vue-components/vite'
import { VantResolver } from 'unplugin-vue-components/resolvers'

export default defineConfig({
  plugins: [
    vue(),
    Components({
      resolvers: [VantResolver()],
    }),
  ],
})
```

## API 请求配置

项目使用Axios处理API请求，封装了请求拦截器和响应拦截器。

### 基本配置

```js
// src/utils/request.js
import axios from 'axios'
import { useUserStore } from '../stores/user'
import { showToast } from 'vant'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  timeout: 10000
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    const userStore = useUserStore()
    const token = userStore.token || localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 0) {
      showToast(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '未知错误'))
    }
    return res.data
  },
  error => {
    showToast(error.message || '网络错误')
    return Promise.reject(error)
  }
)

export default request
```

## 路由配置

使用Vue Router管理应用路由：

```js
// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/home/index.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/auth/login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/questionnaire/:id',
    name: 'QuestionnaireDetail',
    component: () => import('../views/questionnaire/detail.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/create',
    name: 'CreateQuestionnaire',
    component: () => import('../views/questionnaire/create.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/user',
    name: 'UserCenter',
    component: () => import('../views/user/index.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
})

export default router
```

## 环境变量配置

项目使用`.env`文件管理环境变量：

```
# .env.development
VITE_API_BASE_URL=http://localhost:8080/api
VITE_APP_TITLE=问卷系统(开发环境)

# .env.production
VITE_API_BASE_URL=/api
VITE_APP_TITLE=问卷系统
```

## 移动端适配

项目使用`viewport`单位和媒体查询实现移动端适配：

```js
// vite.config.js
import postcssPxToViewport from 'postcss-px-to-viewport'

export default defineConfig({
  css: {
    postcss: {
      plugins: [
        postcssPxToViewport({
          unitToConvert: 'px',
          viewportWidth: 375,
          unitPrecision: 5,
          propList: ['*'],
          viewportUnit: 'vw',
          fontViewportUnit: 'vw',
          selectorBlackList: [],
          minPixelValue: 1,
          mediaQuery: false,
          replace: true,
          exclude: [],
          landscape: false,
          landscapeUnit: 'vw',
          landscapeWidth: 568
        })
      ]
    }
  }
})
```

## 常见问题

### 跨域问题

开发环境下可能遇到跨域问题，可通过Vite的代理功能解决：

```js
// vite.config.js
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  }
})
```

### 移动端调试

推荐使用以下工具进行移动端调试：
- VConsole: 在移动端页面中插入调试面板
- Charles/Fiddler: 抓包工具，分析网络请求
- 浏览器开发者工具的设备模拟功能

## 构建与部署

### 构建优化

```js
// vite.config.js
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { visualizer } from 'rollup-plugin-visualizer'
import viteCompression from 'vite-plugin-compression'

export default defineConfig({
  plugins: [
    vue(),
    viteCompression(), // Gzip压缩
    visualizer() // 构建分析
  ],
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true
      }
    },
    rollupOptions: {
      output: {
        manualChunks: {
          vant: ['vant'],
          echarts: ['echarts']
        }
      }
    }
  }
})
```

### Nginx配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    root /path/to/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api/ {
        proxy_pass http://backend-server:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 贡献指南

欢迎为项目做出贡献！请遵循以下步骤：

1. Fork 项目仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建Pull Request 