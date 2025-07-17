<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'
import { statisticsApi } from '../../api/admin'
import { Toast, Dialog } from 'vant'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(true)
const stats = ref(null)

// 检查是否是管理员
onMounted(async () => {
  if (!userStore.isAdmin) {
    Dialog.confirm({
      title: '权限不足',
      message: '您需要管理员权限才能访问此页面',
      showCancelButton: false,
      confirmButtonText: '返回首页'
    }).then(() => {
      router.push('/')
    })
    return
  }
  
  await loadStatistics()
})

// 加载统计数据
const loadStatistics = async () => {
  loading.value = true
  try {
    const response = await statisticsApi.getSystemStatistics()
    stats.value = response.data.data
    console.log('系统统计数据:', stats.value)
  } catch (error) {
    console.error('获取统计数据失败:', error)
    Toast('获取统计数据失败')
  } finally {
    loading.value = false
  }
}

// 导航到不同的管理页面
const navigateTo = (path) => {
  router.push(path)
}
</script>

<template>
  <div class="admin-dashboard">
    <van-nav-bar
      title="管理员控制台"
      left-text="返回"
      left-arrow
      @click-left="router.push('/')"
    />
    
    <div class="dashboard-container">
      <h2 class="dashboard-title">系统概览</h2>
      
      <van-loading v-if="loading" size="24px" vertical>加载中...</van-loading>
      
      <template v-else-if="stats">
        <!-- 数据卡片 -->
        <div class="stats-cards">
          <!-- 用户统计 -->
          <div class="stats-card">
            <div class="card-header">
              <van-icon name="friends-o" size="24" />
              <h3>用户统计</h3>
            </div>
            <div class="card-content">
              <div class="stat-item">
                <span class="stat-label">总用户数</span>
                <span class="stat-value">{{ stats.user_statistics.total_users }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">管理员用户</span>
                <span class="stat-value">{{ stats.user_statistics.admin_users }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">普通用户</span>
                <span class="stat-value">{{ stats.user_statistics.normal_users }}</span>
              </div>
            </div>
            <div class="card-footer">
              <van-button type="primary" size="small" @click="navigateTo('/admin/users')">管理用户</van-button>
            </div>
          </div>
          
          <!-- 问卷统计 -->
          <div class="stats-card">
            <div class="card-header">
              <van-icon name="description" size="24" />
              <h3>问卷统计</h3>
            </div>
            <div class="card-content">
              <div class="stat-item">
                <span class="stat-label">总问卷数</span>
                <span class="stat-value">{{ stats.questionnaire_statistics.total_questionnaires }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">已发布问卷</span>
                <span class="stat-value">{{ stats.questionnaire_statistics.published_questionnaires }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">问题总数</span>
                <span class="stat-value">{{ stats.questionnaire_statistics.total_questions }}</span>
              </div>
            </div>
            <div class="card-footer">
              <van-button type="primary" size="small" @click="navigateTo('/questionnaire/list')">管理问卷</van-button>
            </div>
          </div>
          
          <!-- 提交统计 -->
          <div class="stats-card">
            <div class="card-header">
              <van-icon name="records" size="24" />
              <h3>提交统计</h3>
            </div>
            <div class="card-content">
              <div class="stat-item">
                <span class="stat-label">总提交数</span>
                <span class="stat-value">{{ stats.submission_statistics.total_submissions }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">总回答数</span>
                <span class="stat-value">{{ stats.submission_statistics.total_answers }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">最近7天提交</span>
                <span class="stat-value">{{ stats.submission_statistics.recent_submissions }}</span>
              </div>
            </div>
            <div class="card-footer">
              <van-button type="primary" size="small" @click="navigateTo('/admin/questionnaires')">查看详情</van-button>
            </div>
          </div>
        </div>
        
        <!-- 快速导航 -->
        <div class="quick-nav">
          <h3>快速导航</h3>
          <div class="nav-buttons">
            <van-button icon="manager-o" type="primary" @click="navigateTo('/admin/users')">用户管理</van-button>
            <van-button icon="todo-list-o" type="info" @click="navigateTo('/questionnaire/list')">问卷管理</van-button>
            <van-button icon="chart-trending-o" type="warning" @click="navigateTo('/admin/statistics')">数据分析</van-button>
          </div>
        </div>
      </template>
      
      <div v-else class="error-message">
        <van-empty description="无法加载统计数据">
          <template #bottom>
            <van-button round type="primary" @click="loadStatistics">重试</van-button>
          </template>
        </van-empty>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-dashboard {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.dashboard-container {
  padding: 16px;
}

.dashboard-title {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 20px;
  color: #323233;
  text-align: center;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stats-card {
  background-color: #fff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 12px rgba(100, 101, 102, 0.08);
}

.card-header {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  color: #323233;
}

.card-header h3 {
  margin: 0 0 0 8px;
  font-size: 16px;
  font-weight: bold;
}

.card-content {
  margin-bottom: 16px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
}

.stat-label {
  color: #646566;
}

.stat-value {
  font-weight: bold;
  color: #323233;
}

.card-footer {
  display: flex;
  justify-content: flex-end;
}

.quick-nav {
  background-color: #fff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 12px rgba(100, 101, 102, 0.08);
}

.quick-nav h3 {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 16px;
  color: #323233;
}

.nav-buttons {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 12px;
}

.error-message {
  margin-top: 40px;
}
</style> 