import { defineStore } from 'pinia'
import axios from 'axios'

// 创建axios实例
const api = axios.create({
  baseURL: '/api', // 使用相对路径，会通过Vite代理转发
  timeout: 30000, // 增加超时时间到30秒
  headers: {
    'Content-Type': 'application/json'
  }
});

// 添加请求拦截器
api.interceptors.request.use(
  config => {
    console.log('发送请求:', config.url, config.data)
    return config
  },
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 添加响应拦截器
api.interceptors.response.use(
  response => {
    console.log('收到响应:', response.data)
    return response
  },
  error => {
    console.error('响应错误:', error.response || error.message || error)
    return Promise.reject(error)
  }
)

// 定义用户状态管理
export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: JSON.parse(localStorage.getItem('userInfo') || '{}')
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.userInfo.is_admin || false
  },
  
  actions: {
    // 登录
    async login(username, password) {
      try {
        console.log(`尝试登录: 用户名=${username}`)
        
        // 使用axios发送请求
        const response = await api.post('/user/login', {
          username,
          password
        });
        
        console.log('登录响应:', response.data);
        
        const data = response.data;
        
        this.token = data.token
        this.userInfo = {
          id: data.user_id,
          username: data.username,
          email: data.email,
          is_admin: data.is_admin
        }
        
        // 保存到本地存储
        localStorage.setItem('token', this.token)
        localStorage.setItem('userInfo', JSON.stringify(this.userInfo))
        
        console.log('登录成功:', this.userInfo)
        return Promise.resolve(data)
      } catch (error) {
        console.error('登录失败:', error)
        
        // 提取错误信息
        let errorMessage = '登录失败';
        
        if (error.response) {
          // 服务器返回了错误响应
          console.error('错误响应状态:', error.response.status);
          console.error('错误响应数据:', error.response.data);
          
          if (error.response.data && error.response.data.error) {
            errorMessage = error.response.data.error;
          } else {
            errorMessage = `服务器错误 (${error.response.status})`;
          }
        } else if (error.request) {
          // 请求已发送但没有收到响应
          console.error('没有收到响应:', error.request);
          errorMessage = '服务器无响应，请检查后端服务是否启动';
        } else {
          // 请求配置出错
          console.error('请求错误:', error.message);
          errorMessage = error.message;
        }
        
        return Promise.reject(new Error(errorMessage));
      }
    },
    
    // 注册
    async register(userData) {
      try {
        console.log('尝试注册:', userData)
        
        // 使用axios发送请求
        const response = await api.post('/user/register', userData);
        
        console.log('注册响应:', response.data);
        const data = response.data;
        
        console.log('注册成功:', data)
        return Promise.resolve(data)
      } catch (error) {
        console.error('注册失败:', error)
        
        // 提取错误信息
        let errorMessage = '注册失败';
        
        if (error.response) {
          // 服务器返回了错误响应
          console.error('错误响应状态:', error.response.status);
          console.error('错误响应数据:', error.response.data);
          
          if (error.response.data && error.response.data.error) {
            errorMessage = error.response.data.error;
          } else {
            errorMessage = `服务器错误 (${error.response.status})`;
          }
        } else if (error.request) {
          // 请求已发送但没有收到响应
          console.error('没有收到响应:', error.request);
          errorMessage = '服务器无响应，请检查后端服务是否启动';
        } else {
          // 请求配置出错
          console.error('请求错误:', error.message);
          errorMessage = error.message;
        }
        
        return Promise.reject(new Error(errorMessage));
      }
    },
    
    // 登出
    logout() {
      this.token = ''
      this.userInfo = {}
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      console.log('用户已登出')
    }
  }
}) 