import axios from 'axios'

// 创建带有超时设置的axios实例
export const api = axios.create({
  baseURL: '/api',
  //// 修改为您的服务器域名或IP
//baseURL: 'http://your-server-domain-or-ip/api'
  timeout: 30000, // 30秒超时
  headers: {
    'Content-Type': 'application/json',
    'Cache-Control': 'no-cache',
    'Pragma': 'no-cache'
  }
});

// 添加请求拦截器
api.interceptors.request.use(
  config => {
    const startTime = new Date().getTime();
    config.metadata = { startTime }; // 记录请求开始时间
    
    console.log(`发送请求: ${config.url}, 方法: ${config.method}, 时间: ${new Date().toISOString()}`);
    
    // 自动添加token到请求头
    const token = localStorage.getItem('token');
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    
    return config;
  },
  error => {
    console.error('请求错误:', error);
    return Promise.reject(error);
  }
);

// 添加响应拦截器
api.interceptors.response.use(
  response => {
    const endTime = new Date().getTime();
    const startTime = response.config.metadata ? response.config.metadata.startTime : 0;
    const duration = endTime - startTime;
    
    console.log(`收到响应: ${response.config.url}, 状态: ${response.status}, 耗时: ${duration}ms`);
    return response.data;
  },
  async error => {
    console.error(`响应错误: ${error.message || error}, URL: ${error.config?.url || '未知'}`);
    
    // 如果是超时错误，并且请求配置中没有设置重试标记，则尝试重试
    const { config } = error;
    if (config && error.code === 'ECONNABORTED' && error.message.includes('timeout') && !config._retryCount) {
      // 设置重试计数
      config._retryCount = 1;
      console.log(`请求超时，正在重试 (${config.url})...`);
      
      try {
        // 使用fetch API作为备选方案重试
        const url = `${config.baseURL || ''}${config.url}`.replace(/\/\//g, '/');
        console.log(`使用fetch API重试请求: ${url}`);
        
        const fetchResponse = await fetch(url, {
          method: config.method.toUpperCase(),
          headers: {
            ...config.headers,
            'Authorization': localStorage.getItem('token') ? `Bearer ${localStorage.getItem('token')}` : ''
          },
          body: config.data ? JSON.stringify(config.data) : undefined
        });
        
        if (!fetchResponse.ok) {
          throw new Error(`Fetch请求失败: ${fetchResponse.status} ${fetchResponse.statusText}`);
        }
        
        const data = await fetchResponse.json();
        console.log(`Fetch重试成功: ${url}`);
        return data;
      } catch (retryError) {
        console.error('重试失败:', retryError);
        return Promise.reject(retryError);
      }
    }
    
    return Promise.reject(error);
  }
);

export default api; 