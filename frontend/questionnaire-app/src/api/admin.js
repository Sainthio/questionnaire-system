import axios from 'axios'

// 创建axios实例
const api = axios.create({
  baseURL: '/api/admin', // 管理员API路径
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器 - 添加token
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response
  },
  error => {
    // 处理401和403错误
    if (error.response && (error.response.status === 401 || error.response.status === 403)) {
      // 清除本地存储的token和用户信息
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      
      // 跳转到登录页面
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// 用户管理API
export const userApi = {
  // 获取所有用户
  getAllUsers(page = 1, pageSize = 10) {
    return api.get('/users', { params: { page, page_size: pageSize } })
  },
  
  // 获取用户详情
  getUserDetail(id) {
    return api.get('/user/detail', { params: { id } })
  },
  
  // 更新用户信息
  updateUser(userData) {
    return api.put('/user/update', userData)
  },
  
  // 删除用户
  deleteUser(id) {
    return api.delete('/user/delete', { params: { id } })
  }
}

// 问卷管理API
export const questionnaireApi = {
  // 获取所有问卷
  getAllQuestionnaires(page = 1, pageSize = 10) {
    return api.get('/questionnaires', { params: { page, page_size: pageSize } })
  },
  
  // 获取问卷提交详情
  getQuestionnaireSubmissions(id) {
    return api.get('/questionnaire/submissions', { params: { id } })
  }
}

// 系统统计API
export const statisticsApi = {
  // 获取系统统计信息
  getSystemStatistics() {
    return api.get('/statistics')
  }
}

export default {
  userApi,
  questionnaireApi,
  statisticsApi
} 