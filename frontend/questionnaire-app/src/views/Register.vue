<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { Toast } from 'vant'

const router = useRouter()
const userStore = useUserStore()

const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const email = ref('')
const phone = ref('')
const loading = ref(false)
const errorMsg = ref('')

const handleRegister = async () => {
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
  
  if (password.value !== confirmPassword.value) {
    Toast('两次输入的密码不一致')
    return
  }
  
  if (!email.value) {
    Toast('请输入邮箱')
    return
  }
  
  loading.value = true
  Toast.loading({
    message: '注册中...',
    forbidClick: true,
    duration: 0
  })
  
  try {
    console.log('开始注册请求...')
    await userStore.register({
      username: username.value,
      password: password.value,
      email: email.value,
      phone: phone.value
    })
    Toast.clear()
    Toast.success('注册成功')
    router.push('/login')
  } catch (error) {
    console.error('注册错误详情:', error)
    Toast.clear()
    
    // 提取错误信息
    let message = '注册失败'
    if (error.message) {
      // 尝试从错误信息中提取具体原因
      if (error.message.includes('用户名已存在')) {
        message = '用户名已存在'
      } else if (error.message.includes('邮箱已存在')) {
        message = '邮箱已存在'
      } else {
        message = error.message
      }
    }
    
    // 显示错误信息
    errorMsg.value = message
    Toast.fail(message)
  } finally {
    loading.value = false
  }
}

const goToLogin = () => {
  router.push('/login')
}
</script>

<template>
  <div class="register-container">
    <van-nav-bar
      title="用户注册"
      left-text="返回"
      left-arrow
      @click-left="router.push('/')"
    />
    
    <div class="form-container">
      <van-form @submit="handleRegister">
        <van-cell-group inset>
          <van-field
            v-model="username"
            name="username"
            label="用户名"
            placeholder="请输入用户名"
            :rules="[{ required: true, message: '请输入用户名' }]"
          />
          <van-field
            v-model="password"
            type="password"
            name="password"
            label="密码"
            placeholder="请输入密码"
            :rules="[{ required: true, message: '请输入密码' }]"
          />
          <van-field
            v-model="confirmPassword"
            type="password"
            name="confirmPassword"
            label="确认密码"
            placeholder="请再次输入密码"
            :rules="[{ required: true, message: '请再次输入密码' }]"
          />
          <van-field
            v-model="email"
            name="email"
            label="邮箱"
            placeholder="请输入邮箱"
            :rules="[{ required: true, message: '请输入邮箱' }]"
          />
          <van-field
            v-model="phone"
            name="phone"
            label="手机号"
            placeholder="请输入手机号（选填）"
          />
        </van-cell-group>
        
        <div v-if="errorMsg" class="error-message">{{ errorMsg }}</div>
        
        <div style="margin: 16px;">
          <van-button round block type="primary" native-type="submit" :loading="loading">
            注册
          </van-button>
          <div class="login-link">
            已有账号？<a @click="goToLogin">立即登录</a>
          </div>
        </div>
      </van-form>
    </div>
  </div>
</template>

<style scoped>
.register-container {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.form-container {
  padding: 20px;
}

.login-link {
  text-align: center;
  margin-top: 16px;
  font-size: 14px;
  color: #969799;
}

.login-link a {
  color: #1989fa;
  text-decoration: none;
}

.error-message {
  color: #ee0a24;
  font-size: 14px;
  text-align: center;
  margin: 10px 0;
  padding: 8px;
  background-color: #fff1f0;
  border-radius: 4px;
}
</style> 