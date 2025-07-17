# 问卷系统前端

这是问卷系统的前端部分，使用Vue3 + Vant开发。

## 安装步骤

1. 安装Node.js (https://nodejs.org/)

2. 创建Vue项目:
```bash
npm create vite@latest questionnaire-app -- --template vue
cd questionnaire-app
```

3. 安装依赖:
```bash
npm install
npm install vant@next
npm install axios
npm install vue-router@4
npm install pinia
```

4. 启动开发服务器:
```bash
npm run dev
```

## 项目结构

创建后的项目结构应该如下:

```
questionnaire-app/
  ├── public/              # 静态资源
  ├── src/                 # 源代码
  │   ├── assets/          # 资源文件
  │   ├── components/      # 组件
  │   ├── router/          # 路由
  │   ├── stores/          # 状态管理
  │   ├── views/           # 页面视图
  │   ├── App.vue          # 根组件
  │   └── main.js          # 入口文件
  ├── index.html           # HTML模板
  ├── package.json         # 项目配置
  └── vite.config.js       # Vite配置
```

## 主要功能页面

1. 用户登录/注册
2. 问卷列表
3. 问卷详情/填写
4. 问卷创建/编辑
5. 统计结果查看

## 使用Vant组件

在main.js中引入Vant:

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

## API配置

在src/utils/request.js中配置axios:

```js
import axios from 'axios'

const request = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 5000
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
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
    return response.data
  },
  error => {
    return Promise.reject(error)
  }
)

export default request
``` 