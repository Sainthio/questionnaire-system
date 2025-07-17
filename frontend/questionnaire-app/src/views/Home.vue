<script setup>
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useQuestionnaireStore } from '../stores/questionnaire'
import { Toast } from 'vant'
import api from '../services/api'

const router = useRouter()
const userStore = useUserStore()
const questionnaireStore = useQuestionnaireStore()
const username = ref('')
const showAnimation = ref(false)
const stats = reactive({
  users: 0,
  questionnaires: 0,
  submissions: 0
})

// 获取系统统计数据
const getSystemStats = async () => {
  console.log('正在获取系统统计数据...')
  try {
    const statsResponse = await api.get('/questionnaire/stats')
    console.log('统计接口返回数据:', statsResponse)
    
    // 检查响应数据是否有效
    if (statsResponse && statsResponse.code === 0) {
      // 直接从响应中获取数据
      stats.users = statsResponse.user_count || 0
      stats.questionnaires = statsResponse.questionnaire_count || 0
      stats.submissions = statsResponse.submission_count || 0
      
      console.log('统计数据获取成功:', stats)
    } else {
      console.warn('统计接口返回数据为空，使用备用方法获取')
      
      // 备用方法：获取问卷总数
      const questionnairesResponse = await api.get('/questionnaire/list?page=1&limit=1')
      if (questionnairesResponse && questionnairesResponse.data && questionnairesResponse.data.data) {
        const responseData = questionnairesResponse.data.data
        stats.questionnaires = responseData.total || 0
        
        // 使用合理的估算值
        stats.users = Math.max(5, Math.floor(stats.questionnaires / 2))
        stats.submissions = Math.max(10, stats.questionnaires * 3)
      }
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
    // 使用默认值
    stats.questionnaires = stats.questionnaires || 3
    stats.users = 5
    stats.submissions = 10
  }
}

onMounted(async () => {
  if (userStore.isLoggedIn) {
    username.value = userStore.userInfo.username
  }
  
  // 添加进入动画
  setTimeout(() => {
    showAnimation.value = true
  }, 100)
  
  // 获取统计数据
  await getSystemStats()
})

const goToLogin = () => {
  router.push('/login')
}

const goToRegister = () => {
  router.push('/register')
}

const goToList = () => {
  router.push('/questionnaire/list')
}

const logout = () => {
  userStore.logout()
  Toast('已退出登录')
}
</script>

<template>
  <div class="home-container">
    <van-nav-bar title="问卷系统" class="custom-nav">
      <template #right>
        <van-icon name="apps-o" size="18" />
      </template>
    </van-nav-bar>
    
    <div class="content" :class="{ 'show-animation': showAnimation }">
      <div class="hero-section">
        <van-image
          width="120"
          height="120"
          src="https://s2.loli.net/2025/06/27/ovsYhmTfZ1DqtVB.gif"
          class="logo"
          radius="50%"
        />
        <!-- 图片地址 src="https://fastly.jsdelivr.net/npm/@vant/assets/cat.jpeg"-->
        
        <div class="text-content">
          <h1 class="title">问卷系统</h1>
          <p class="subtitle">基于Vue3 + Go开发</p>
        </div>
      </div>
      
      <div class="stats-row">
        <div class="stat-item">
          <div class="stat-number">{{ stats.users }}</div>
          <div class="stat-label">注册用户</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ stats.questionnaires }}</div>
          <div class="stat-label">问卷总数</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ stats.submissions }}</div>
          <div class="stat-label">提交次数</div>
        </div>
      </div>
      
      <div class="feature-showcase">
        <div class="feature-title-area">
          <h2 class="section-title">系统特点</h2>
          <div class="title-underline"></div>
        </div>
        
        <div class="feature-timeline">
          <div class="timeline-item">
            <div class="timeline-icon">
              <van-icon name="edit" />
            </div>
            <div class="timeline-content">
              <h3>便捷设计</h3>
              <p>简单直观的问卷设计流程，无需专业知识</p>
            </div>
          </div>
          
          <div class="timeline-item">
            <div class="timeline-icon">
              <van-icon name="chart-trending-o" />
            </div>
            <div class="timeline-content">
              <h3>数据可视</h3>
              <p>清晰的数据展示，让信息一目了然</p>
            </div>
          </div>
          
          <div class="timeline-item">
            <div class="timeline-icon">
              <van-icon name="friends-o" />
            </div>
            <div class="timeline-content">
              <h3>多端兼容</h3>
              <p>同时支持移动端与桌面端的完美体验</p>
            </div>
          </div>
        </div>
      </div>
      
      <div class="button-group">
        <template v-if="userStore.isLoggedIn">
          <p class="welcome">欢迎您，{{ username }}</p>
          <van-button type="primary" size="large" @click="goToList" icon="bars" block>问卷列表</van-button>
          <van-button v-if="userStore.isAdmin" type="success" size="large" @click="router.push('/questionnaire/create')" icon="plus" block>创建问卷</van-button>
          <van-button v-if="userStore.isAdmin" type="warning" size="large" @click="router.push('/admin')" icon="manager-o" block>管理员控制台</van-button>
          <van-button type="danger" size="large" @click="logout" icon="cross" block>退出登录</van-button>
        </template>
        <template v-else>
          <van-button type="primary" size="large" @click="goToLogin" icon="user-o" block>登录</van-button>
          <van-button type="info" size="large" @click="goToRegister" icon="friends-o" block>注册</van-button>
          <van-button type="default" size="large" @click="goToList" icon="browsing-history-o" block>浏览问卷</van-button>
        </template>
      </div>
      
      <div class="footer">
        <p>© 2025 问卷系统 · 版权所有</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-container {
  min-height: 100vh;
  background-color: #f7f8fa;
  overflow-x: hidden;
}

.custom-nav {
  background: linear-gradient(135deg, #4481eb 0%, #04befe 100%);
}

:deep(.custom-nav .van-nav-bar__title) {
  color: white;
  font-weight: bold;
}

:deep(.custom-nav .van-icon) {
  color: white !important;
}

.content {
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  opacity: 0;
  transform: translateY(20px);
  transition: all 0.5s ease;
}

.show-animation {
  opacity: 1;
  transform: translateY(0);
}

.hero-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 20px;
  padding: 20px;
  width: 100%;
  max-width: 500px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.logo {
  margin-bottom: 20px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease;
}

.logo:hover {
  transform: scale(1.05);
}

.text-content {
  text-align: center;
}

.title {
  font-size: 28px;
  font-weight: bold;
  margin-bottom: 8px;
  color: #323233;
  background: linear-gradient(135deg, #4481eb 0%, #04befe 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.subtitle {
  font-size: 16px;
  color: #969799;
  margin-bottom: 0;
}

.stats-row {
  display: flex;
  justify-content: space-between;
  width: 100%;
  max-width: 500px;
  margin-bottom: 20px;
}

.stat-item {
  flex: 1;
  text-align: center;
  background: white;
  padding: 15px 10px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin: 0 5px;
}

.stat-number {
  font-size: 20px;
  font-weight: bold;
  color: #4481eb;
}

.stat-label {
  font-size: 12px;
  color: #969799;
  margin-top: 5px;
}

.feature-showcase {
  width: 100%;
  max-width: 500px;
  margin-bottom: 20px;
}

.feature-title-area {
  text-align: center;
  margin-bottom: 10px;
}

.section-title {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 8px;
  color: #323233;
}

.title-underline {
  height: 2px;
  background-color: #4481eb;
  width: 50px;
  margin: 0 auto;
}

.feature-timeline {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.timeline-item {
  flex: 1;
  text-align: center;
  background: white;
  padding: 15px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.timeline-icon {
  font-size: 24px;
  color: #4481eb;
  margin-bottom: 10px;
}

.timeline-content {
  text-align: center;
}

.timeline-content h3 {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 5px;
}

.timeline-content p {
  font-size: 14px;
  color: #969799;
}

.button-group {
  width: 100%;
  max-width: 500px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 10px;
}

.welcome {
  text-align: center;
  margin-bottom: 10px;
  font-size: 16px;
  color: #323233;
  font-weight: 500;
}

:deep(.van-button) {
  border-radius: 8px;
  height: 44px;
  font-size: 16px;
  transition: transform 0.2s ease;
}

:deep(.van-button:active) {
  transform: scale(0.98);
}

.footer {
  margin-top: 30px;
  text-align: center;
  font-size: 12px;
  color: #969799;
}

/* 桌面端样式 */
:deep(.desktop) .hero-section {
  flex-direction: row;
  justify-content: center;
  padding: 30px;
}

:deep(.desktop) .logo {
  margin-right: 30px;
  margin-bottom: 0;
}

:deep(.desktop) .text-content {
  text-align: left;
}

:deep(.desktop) .feature-showcase {
  max-width: 700px;
}

:deep(.desktop) .button-group {
  max-width: 500px;
}
</style> 