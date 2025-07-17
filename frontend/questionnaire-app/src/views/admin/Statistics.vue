<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'
import { statisticsApi, questionnaireApi } from '../../api/admin'
import { Toast, Dialog } from 'vant'
import * as echarts from 'echarts/core'
import { PieChart, BarChart, LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

// 注册 ECharts 组件
echarts.use([
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  PieChart,
  BarChart,
  LineChart,
  CanvasRenderer
])

const router = useRouter()
const userStore = useUserStore()
const loading = ref(true)
const stats = ref(null)

// 问卷数据
const topQuestionnaires = ref([])

// 图表实例
const userChartRef = ref(null)
const questionnaireChartRef = ref(null)
const submissionChartRef = ref(null)
const submissionTrendChartRef = ref(null)
let userChart = null
let questionnaireChart = null
let submissionChart = null
let submissionTrendChart = null

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
  await loadTopQuestionnaires()
  
  // 初始化图表
  setTimeout(() => {
    initCharts()
  }, 500)
})

// 加载统计数据
const loadStatistics = async () => {
  loading.value = true
  try {
    const response = await statisticsApi.getSystemStatistics()
    stats.value = response.data.data
    
    // 确保所有必要的对象结构存在
    if (!stats.value) {
      stats.value = {}
    }
    
    if (!stats.value.user_statistics) {
      stats.value.user_statistics = {
        total_users: 0,
        admin_users: 0,
        normal_users: 0
      }
    }
    
    if (!stats.value.questionnaire_statistics) {
      stats.value.questionnaire_statistics = {
        total_questionnaires: 0,
        published_questionnaires: 0,
        unpublished_questionnaires: 0,
        total_questions: 0
      }
    }
    
    if (!stats.value.submission_statistics) {
      stats.value.submission_statistics = {
        total_submissions: 0,
        total_answers: 0,
        recent_submissions: 0,
        average_answers_per_submission: 0
      }
    }
    
    console.log('系统统计数据:', stats.value)
  } catch (error) {
    console.error('获取统计数据失败:', error)
    Toast('获取统计数据失败')
    
    // 创建默认数据结构
    stats.value = {
      user_statistics: {
        total_users: 0,
        admin_users: 0,
        normal_users: 0
      },
      questionnaire_statistics: {
        total_questionnaires: 0,
        published_questionnaires: 0,
        unpublished_questionnaires: 0,
        total_questions: 0
      },
      submission_statistics: {
        total_submissions: 0,
        total_answers: 0,
        recent_submissions: 0,
        average_answers_per_submission: 0
      }
    }
  } finally {
    loading.value = false
  }
}

// 加载热门问卷数据（假数据，实际应该从API获取）
const loadTopQuestionnaires = async () => {
  try {
    const response = await questionnaireApi.getAllQuestionnaires(1, 5)
    topQuestionnaires.value = response.data.data.questionnaires || []
    console.log('热门问卷:', topQuestionnaires.value)
  } catch (error) {
    console.error('获取热门问卷失败:', error)
    topQuestionnaires.value = []
  }
}

// 计算属性：用户统计数据
const userData = computed(() => {
  if (!stats.value || !stats.value.user_statistics) return []
  return [
    { name: '管理员', value: stats.value.user_statistics.admin_users || 0 },
    { name: '普通用户', value: stats.value.user_statistics.normal_users || 0 }
  ]
})

// 计算属性：问卷统计数据
const questionnaireData = computed(() => {
  if (!stats.value || !stats.value.questionnaire_statistics) return []
  return [
    { name: '已发布', value: stats.value.questionnaire_statistics.published_questionnaires || 0 },
    { name: '未发布', value: stats.value.questionnaire_statistics.unpublished_questionnaires || 0 }
  ]
})

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

// 初始化图表
const initCharts = () => {
  if (!stats.value) return
  
  // 用户分布图表
  if (userChartRef.value) {
    userChart = echarts.init(userChartRef.value)
    userChart.setOption({
      title: {
        text: '用户分布',
        left: 'center'
      },
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        left: 'left',
        data: ['管理员', '普通用户']
      },
      series: [
        {
          name: '用户类型',
          type: 'pie',
          radius: '50%',
          data: userData.value,
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          }
        }
      ]
    })
  }
  
  // 问卷状态图表
  if (questionnaireChartRef.value) {
    questionnaireChart = echarts.init(questionnaireChartRef.value)
    questionnaireChart.setOption({
      title: {
        text: '问卷状态',
        left: 'center'
      },
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        left: 'left',
        data: ['已发布', '未发布']
      },
      series: [
        {
          name: '问卷状态',
          type: 'pie',
          radius: '50%',
          data: questionnaireData.value,
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          }
        }
      ]
    })
  }
  
  // 提交统计图表
  if (submissionChartRef.value && stats.value && stats.value.submission_statistics) {
    submissionChart = echarts.init(submissionChartRef.value)
    submissionChart.setOption({
      title: {
        text: '提交统计',
        left: 'center'
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: ['总提交数', '最近7天提交']
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '提交数量',
          type: 'bar',
          data: [
            stats.value.submission_statistics.total_submissions || 0,
            stats.value.submission_statistics.recent_submissions || 0
          ],
          itemStyle: {
            color: function(params) {
              const colorList = ['#1989fa', '#07c160']
              return colorList[params.dataIndex]
            }
          }
        }
      ]
    })
  }
  
  // 提交趋势图表（模拟数据）
  if (submissionTrendChartRef.value) {
    // 生成模拟的最近7天数据
    const days = []
    const data = []
    
    const today = new Date()
    for (let i = 6; i >= 0; i--) {
      const date = new Date(today)
      date.setDate(date.getDate() - i)
      days.push(formatDate(date))
      
      // 生成随机数据，实际应该从后端获取
      const randomValue = Math.floor(Math.random() * 10) + 1
      data.push(randomValue)
    }
    
    submissionTrendChart = echarts.init(submissionTrendChartRef.value)
    submissionTrendChart.setOption({
      title: {
        text: '最近7天提交趋势',
        left: 'center'
      },
      tooltip: {
        trigger: 'axis'
      },
      xAxis: {
        type: 'category',
        data: days
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '提交数量',
          type: 'line',
          data: data,
          smooth: true,
          itemStyle: {
            color: '#1989fa'
          },
          lineStyle: {
            color: '#1989fa'
          },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: 'rgba(25, 137, 250, 0.5)'
                },
                {
                  offset: 1,
                  color: 'rgba(25, 137, 250, 0.05)'
                }
              ]
            }
          }
        }
      ]
    })
  }
}

// 窗口大小变化时重绘图表
window.addEventListener('resize', () => {
  userChart?.resize()
  questionnaireChart?.resize()
  submissionChart?.resize()
  submissionTrendChart?.resize()
})

// 导航到问卷管理页面
const goToQuestionnaireManagement = () => {
  router.push('/questionnaire/list')
}

// 导航到用户管理页面
const goToUserManagement = () => {
  router.push('/admin/users')
}

// 查看问卷详情
const viewQuestionnaireDetail = (id) => {
  router.push(`/questionnaire/detail/${id}`)
}
</script>

<template>
  <div class="admin-statistics">
    <van-nav-bar
      title="系统数据统计"
      left-text="返回"
      left-arrow
      @click-left="router.push('/admin')"
    />
    
    <div class="statistics-container">
      <van-loading v-if="loading" size="24px" vertical>加载中...</van-loading>
      
      <template v-else-if="stats">
        <!-- 系统概览 -->
        <div class="section-title-bar">
          <h2>系统概览</h2>
        </div>
        
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
                <span class="stat-value">{{ stats.user_statistics?.total_users || 0 }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">管理员用户</span>
                <span class="stat-value">{{ stats.user_statistics?.admin_users || 0 }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">普通用户</span>
                <span class="stat-value">{{ stats.user_statistics?.normal_users || 0 }}</span>
              </div>
            </div>
            <div class="card-footer">
              <van-button type="primary" size="small" @click="goToUserManagement">管理用户</van-button>
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
                <span class="stat-value">{{ stats.questionnaire_statistics?.total_questionnaires || 0 }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">已发布问卷</span>
                <span class="stat-value">{{ stats.questionnaire_statistics?.published_questionnaires || 0 }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">问题总数</span>
                <span class="stat-value">{{ stats.questionnaire_statistics?.total_questions || 0 }}</span>
              </div>
            </div>
            <div class="card-footer">
              <van-button type="primary" size="small" @click="goToQuestionnaireManagement">管理问卷</van-button>
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
                <span class="stat-value">{{ stats.submission_statistics?.total_submissions || 0 }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">总回答数</span>
                <span class="stat-value">{{ stats.submission_statistics?.total_answers || 0 }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">最近7天提交</span>
                <span class="stat-value">{{ stats.submission_statistics?.recent_submissions || 0 }}</span>
              </div>
            </div>
            <div class="card-footer">
              <van-button type="primary" size="small" @click="router.push('/admin/questionnaires')">查看详情</van-button>
            </div>
          </div>
        </div>
        
        <!-- 图表区域 -->
        <div class="section-title-bar">
          <h2>数据可视化</h2>
        </div>
        
        <div class="charts-grid">
          <!-- 用户分布图表 -->
          <div class="chart-card">
            <div class="chart" ref="userChartRef"></div>
          </div>
          
          <!-- 问卷状态图表 -->
          <div class="chart-card">
            <div class="chart" ref="questionnaireChartRef"></div>
          </div>
          
          <!-- 提交统计图表 -->
          <div class="chart-card">
            <div class="chart" ref="submissionChartRef"></div>
          </div>
          
          <!-- 提交趋势图表 -->
          <div class="chart-card full-width">
            <div class="chart" ref="submissionTrendChartRef"></div>
          </div>
        </div>
        
        <!-- 热门问卷 -->
        <div class="section-title-bar">
          <h2>热门问卷统计</h2>
        </div>
        
        <div class="top-questionnaires">
          <van-empty 
            v-if="topQuestionnaires.length === 0" 
            description="暂无问卷数据" 
            image="search"
          />
          
          <div v-else class="questionnaire-list">
            <van-swipe-cell 
              v-for="item in topQuestionnaires" 
              :key="item.questionnaire.id"
              class="questionnaire-item"
            >
              <van-cell 
                :title="item.questionnaire.title" 
                :label="`创建者: ${item.creator_name || '未知'} | 提交数: ${item.submission_count}`"
                @click="viewQuestionnaireDetail(item.questionnaire.id)"
              >
                <template #right-icon>
                  <van-tag 
                    :type="item.questionnaire.is_published ? 'primary' : 'warning'"
                  >
                    {{ item.questionnaire.is_published ? '已发布' : '未发布' }}
                  </van-tag>
                </template>
              </van-cell>
              
              <template #right>
                <van-button 
                  square 
                  type="primary" 
                  class="view-button" 
                  @click="viewQuestionnaireDetail(item.questionnaire.id)"
                >
                  查看
                </van-button>
              </template>
            </van-swipe-cell>
          </div>
          
          <div class="view-more">
            <van-button plain hairline type="primary" size="small" @click="goToQuestionnaireManagement">
              查看更多问卷
            </van-button>
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
.admin-statistics {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.statistics-container {
  padding: 16px;
}

.section-title-bar {
  margin: 24px 0 16px;
  padding-left: 10px;
  border-left: 3px solid #1989fa;
}

.section-title-bar h2 {
  font-size: 18px;
  margin: 0;
  color: #323233;
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

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.chart-card {
  background-color: #fff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 12px rgba(100, 101, 102, 0.08);
}

.full-width {
  grid-column: 1 / -1;
}

.chart {
  height: 300px;
  width: 100%;
}

.top-questionnaires {
  background-color: #fff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 12px rgba(100, 101, 102, 0.08);
  margin-bottom: 24px;
}

.questionnaire-list {
  margin-bottom: 16px;
}

.questionnaire-item {
  margin-bottom: 8px;
}

.questionnaire-item:last-child {
  margin-bottom: 0;
}

.view-button {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 65px;
}

.view-more {
  text-align: center;
  margin-top: 16px;
}

.error-message {
  margin-top: 40px;
}
</style> 