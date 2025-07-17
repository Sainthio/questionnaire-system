<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { Toast } from 'vant'

const router = useRouter()
const userStore = useUserStore()

const username = ref('')
const password = ref('')
const loading = ref(false)
const errorMsg = ref('')
const showAnimation = ref(false)

// 重置密码相关
const showResetForm = ref(false)
const resetUsername = ref('')
const resetPassword = ref('')
const resetLoading = ref(false)

onMounted(() => {
  // 添加进入动画
  setTimeout(() => {
    showAnimation.value = true
  }, 100)
})

const handleLogin = async () => {
  // 清除之前的错误信息
  errorMsg.value = ''
  
  // 表单验证
  if (!username.value) {
    Toast('请输入用户名')
    return
  }
  
  if (!password.value) {
    Toast('请输入密码')
    return
  }
  
  loading.value = true
  Toast.loading({
    message: '登录中...',
    forbidClick: true,
    duration: 0
  })
  
  try {
    console.log('开始登录请求...')
    console.log(`登录信息: 用户名=${username.value}, 密码长度=${password.value.length}`)
    
    // 记录请求开始时间
    const startTime = new Date().getTime()
    
    await userStore.login(username.value, password.value)
    
    // 计算请求耗时
    const endTime = new Date().getTime()
    const duration = endTime - startTime
    console.log(`登录请求耗时: ${duration}ms`)
    
    Toast.clear()
    Toast.success('登录成功')
    router.push('/')
  } catch (error) {
    console.error('登录错误详情:', error)
    Toast.clear()
    
    // 提取错误信息
    let message = '登录失败'
    if (error.message) {
      try {
        // 尝试解析JSON错误消息
        const errorObj = JSON.parse(error.message)
        if (errorObj.error) {
          message = errorObj.error
        }
      } catch (e) {
        // 如果不是JSON格式，直接使用错误消息
        message = error.message
      }
    }
    
    // 网络错误特殊处理
    if (error.name === 'NetworkError' || message.includes('Network Error')) {
      message = '网络错误，请检查后端服务是否启动'
    }
    
    // 显示错误信息
    errorMsg.value = message
    Toast.fail(message)
  } finally {
    loading.value = false
  }
}

const handleResetPassword = async () => {
  if (!resetUsername.value) {
    Toast('请输入用户名')
    return
  }
  
  if (!resetPassword.value) {
    Toast('请输入新密码')
    return
  }
  
  resetLoading.value = true
  Toast.loading({
    message: '重置密码中...',
    forbidClick: true,
    duration: 0
  })
  
  try {
    const response = await fetch('/api/user/reset-password', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: resetUsername.value,
        new_password: resetPassword.value
      })
    })
    
    const data = await response.json()
    
    if (!response.ok) {
      if (data && data.error) {
        throw new Error(data.error)
      } else {
        throw new Error('重置密码失败')
      }
    }
    
    Toast.clear()
    Toast.success(data.message || '密码重置成功')
    
    // 如果是测试用户，自动填充登录表单
    if (resetUsername.value === 'testuser') {
      username.value = resetUsername.value
      password.value = resetPassword.value
    }
    
    showResetForm.value = false
  } catch (error) {
    console.error('重置密码错误:', error)
    Toast.clear()
    Toast.fail('重置密码失败: ' + error.message)
  } finally {
    resetLoading.value = false
  }
}

const createTestUser = async () => {
  resetUsername.value = 'testuser'
  resetPassword.value = 'test123'
  await handleResetPassword()
}

const goToRegister = () => {
  router.push('/register')
}

const toggleResetForm = () => {
  showResetForm.value = !showResetForm.value
}
</script>

<template>
  <div class="login-container">
    <van-nav-bar
      title="用户登录"
      left-text="返回"
      left-arrow
      @click-left="router.push('/')"
      class="custom-nav"
    />
    
    <div class="form-container" :class="{ 'show-animation': showAnimation }">
      <div class="login-card">
        <div class="login-header">
          <van-icon name="user-circle-o" size="48" class="login-icon" />
          <h2 class="login-title">{{ showResetForm ? '重置密码' : '欢迎登录' }}</h2>
          <p class="login-subtitle">{{ showResetForm ? '输入用户名和新密码' : '请输入您的账号和密码' }}</p>
        </div>
        
        <van-form @submit="handleLogin" v-if="!showResetForm">
          <van-cell-group inset>
            <van-field
              v-model="username"
              name="username"
              label="用户名"
              placeholder="请输入用户名"
              :rules="[{ required: true, message: '请输入用户名' }]"
              left-icon="user-o"
            />
            <van-field
              v-model="password"
              type="password"
              name="password"
              label="密码"
              placeholder="请输入密码"
              :rules="[{ required: true, message: '请输入密码' }]"
              left-icon="lock"
            />
          </van-cell-group>
          
          <div v-if="errorMsg" class="error-message">{{ errorMsg }}</div>
          
          <div class="button-area">
            <van-button round block type="primary" native-type="submit" :loading="loading">
              登录
            </van-button>
            <div class="action-links">
              <div class="register-link">
                还没有账号？<a @click="goToRegister">立即注册</a>
              </div>
              <div class="reset-link">
                <a @click="toggleResetForm">忘记密码？</a>
              </div>
            </div>
            <div class="test-link">
              <van-button size="small" plain hairline type="primary" @click="createTestUser" icon="question-o">创建测试账号</van-button>
            </div>
          </div>
        </van-form>
        
        <!-- 重置密码表单 -->
        <van-form @submit="handleResetPassword" v-else>
          <van-cell-group inset>
            <van-field
              v-model="resetUsername"
              name="resetUsername"
              label="用户名"
              placeholder="请输入用户名"
              :rules="[{ required: true, message: '请输入用户名' }]"
              left-icon="user-o"
            />
            <van-field
              v-model="resetPassword"
              type="password"
              name="resetPassword"
              label="新密码"
              placeholder="请输入新密码"
              :rules="[{ required: true, message: '请输入新密码' }]"
              left-icon="lock"
            />
          </van-cell-group>
          
          <div class="button-area">
            <van-button round block type="primary" native-type="submit" :loading="resetLoading">
              重置密码
            </van-button>
            <div class="reset-link back-link">
              <a @click="toggleResetForm"><van-icon name="arrow-left" /> 返回登录</a>
            </div>
          </div>
        </van-form>
      </div>
      
      <div class="login-footer">
        <p>© 2025 问卷系统 · 安全登录</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  min-height: 100vh;
  background-color: #f7f8fa;
  background-image: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.custom-nav {
  background: transparent;
  box-shadow: none;
}

:deep(.custom-nav .van-nav-bar__title) {
  color: #333;
  font-weight: bold;
}

:deep(.custom-nav .van-icon),
:deep(.custom-nav .van-nav-bar__text) {
  color: #333 !important;
}

.form-container {
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 46px);
  opacity: 0;
  transform: translateY(20px);
  transition: all 0.5s ease;
}

.show-animation {
  opacity: 1;
  transform: translateY(0);
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.login-header {
  text-align: center;
  margin-bottom: 20px;
  padding: 10px 0;
}

.login-icon {
  color: #4481eb;
  margin-bottom: 15px;
  display: inline-block;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.1);
  }
  100% {
    transform: scale(1);
  }
}

.login-title {
  font-size: 22px;
  font-weight: bold;
  color: #323233;
  margin: 0 0 5px;
}

.login-subtitle {
  font-size: 14px;
  color: #969799;
  margin: 0;
}

.button-area {
  margin: 24px 16px 8px;
}

.action-links {
  display: flex;
  justify-content: space-between;
  margin-top: 16px;
}

.register-link, .reset-link, .test-link {
  text-align: center;
  margin-top: 12px;
  font-size: 14px;
  color: #969799;
}

.register-link a, .reset-link a {
  color: #4481eb;
  text-decoration: none;
}

.test-link {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.error-message {
  color: #ee0a24;
  font-size: 14px;
  text-align: center;
  margin: 10px 0;
  padding: 8px;
  background-color: #fff1f0;
  border-radius: 8px;
}

.back-link {
  margin-top: 16px;
  display: flex;
  justify-content: center;
}

.back-link a {
  display: flex;
  align-items: center;
}

.login-footer {
  margin-top: 30px;
  text-align: center;
  font-size: 12px;
  color: #666;
}

/* 桌面端样式 */
:deep(.desktop) .login-card {
  padding: 30px;
  max-width: 450px;
}

:deep(.desktop) .login-header {
  padding: 20px 0;
}

:deep(.desktop) .login-title {
  font-size: 24px;
}
</style> 