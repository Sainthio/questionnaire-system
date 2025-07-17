<script setup>
import { ref, onMounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuestionnaireStore } from '../../stores/questionnaire'
import { useUserStore } from '../../stores/user'
import { Toast, Dialog } from 'vant'
import * as echarts from 'echarts/core'
import { BarChart, PieChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

// 注册必须的组件
echarts.use([
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  BarChart,
  PieChart,
  CanvasRenderer
])

const route = useRoute()
const router = useRouter()
const questionnaireStore = useQuestionnaireStore()
const userStore = useUserStore()

const loading = ref(false)
const questionnaire = ref(null)
const questions = ref([])
const submissions = ref([])
const totalSubmissions = ref(0)
const activeNames = ref([])
const chartInstances = ref({})
const focusedQuestionId = ref(null)

// 统计数据
const statistics = ref({
  questionStats: [],
  userStats: {
    total: 0,
    uniqueUsers: 0
  }
})

// 获取问卷结果
const loadQuestionnaireResults = async () => {
  const id = route.params.id
  if (!id) {
    Toast('问卷ID不能为空')
    router.push('/questionnaire/list')
    return
  }
  
  // 检查是否有特定问题ID需要聚焦
  const focusQuestionId = route.query.focusQuestion
  if (focusQuestionId) {
    focusedQuestionId.value = parseInt(focusQuestionId)
    console.log('聚焦到问题ID:', focusedQuestionId.value)
  }
  
  // 检查用户是否登录
  if (!userStore.isLoggedIn) {
    Dialog.alert({
      title: '未登录',
      message: '请先登录再查看问卷结果',
    }).then(() => {
      router.push('/login')
    })
    return
  }
  
  // 检查用户是否有权限查看问卷结果
  const userInfo = userStore.userInfo
  if (!userInfo) {
    Dialog.alert({
      title: '用户信息错误',
      message: '无法获取用户信息，请重新登录',
    }).then(() => {
      router.push('/login')
    })
    return
  }
  
  loading.value = true
  Toast.loading({
    message: '加载中...',
    forbidClick: true,
    duration: 0
  })
  
  try {
    console.log('尝试获取问卷结果，用户ID:', userInfo.id, '问卷ID:', id)
    
    const response = await questionnaireStore.getQuestionnaireResults(id)
    console.log('问卷结果数据:', response)
    
    // 检查响应数据结构
    if (!response || !response.data) {
      throw new Error('返回数据格式错误')
    }
    
    questionnaire.value = response.data.questionnaire || {}
    questions.value = response.data.questions || []
    submissions.value = response.data.submissions || []
    totalSubmissions.value = response.data.total_submissions || 0
    
    // 如果问卷为空，显示错误
    if (!questionnaire.value.id) {
      throw new Error('问卷不存在或已被删除')
    }
    
    // 检查用户是否是问卷创建者
    if (questionnaire.value.created_by !== userInfo.id && !userInfo.is_admin) {
      throw new Error('您没有权限查看此问卷的结果')
    }
    
    console.log('问题数据:', questions.value)
    console.log('提交记录数据:', submissions.value)
    
    // 确保每个提交记录都有有效的结构
    submissions.value = submissions.value.map(submission => {
      // 添加调试日志
      console.log('处理提交记录:', submission.submission?.id, '答案数量:', submission.answers?.length)
      
      // 检查答案数据
      if (submission.answers && Array.isArray(submission.answers)) {
        submission.answers.forEach(answer => {
          console.log('答案数据:', answer.question_id, answer.content)
        })
      } else {
        console.warn('提交记录缺少答案数据:', submission.submission?.id)
      }
      
      return {
        submission: submission.submission || { id: 0, user_id: 0, submitted_at: new Date() },
        answers: Array.isArray(submission.answers) ? submission.answers : [],
        user_info: submission.user_info || { username: '未知用户' }
      }
    })
    
    // 初始化折叠面板，使用索引
    if (submissions.value.length > 0) {
      // 默认展开第一个提交项
      activeNames.value = [0]
      console.log('初始化折叠面板，展开索引:', activeNames.value)
    }
    
    // 生成统计数据
    generateStatistics()
    
    // 渲染图表
    await nextTick()
    renderCharts()
  } catch (error) {
    console.error('加载问卷结果失败:', error)
    Toast.clear()
    
    // 处理特定错误
    if (error.message && (error.message.includes('权限不足') || error.message.includes('没有权限'))) {
      Dialog.alert({
        title: '无权访问',
        message: '您没有权限查看此问卷的结果',
      }).then(() => {
        router.push('/questionnaire/list')
      })
    } else if (error.response && error.response.status === 403) {
      Dialog.alert({
        title: '无权访问',
        message: '您没有权限查看此问卷的结果',
      }).then(() => {
        router.push('/questionnaire/list')
      })
    } else {
      Toast.fail('加载失败: ' + (error.message || '未知错误'))
    }
  } finally {
    loading.value = false
    Toast.clear()
  }
}

// 生成统计数据
const generateStatistics = () => {
  if (!questions.value.length) {
    console.log('没有问题数据，无法生成统计')
    return
  }
  
  if (!submissions.value.length) {
    console.log('没有提交记录，生成空统计')
    // 即使没有提交，也初始化问题统计结构
    const questionStats = questions.value.map(question => {
      return {
        id: question.id,
        title: question.title,
        type: question.type,
        answers: [],
        answerCounts: {}
      }
    })
    
    statistics.value = {
      questionStats,
      userStats: {
        total: 0,
        uniqueUsers: 0
      }
    }
    return
  }
  
  // 初始化问题统计
  const questionStats = questions.value.map(question => {
    return {
      id: question.id,
      title: question.title,
      type: question.type,
      answers: [],
      answerCounts: {}
    }
  })
  
  console.log('开始生成统计数据，问题数:', questionStats.length, '提交数:', submissions.value.length)
  
  // 统计每个问题的答案
  submissions.value.forEach((submission, index) => {
    console.log(`处理第${index+1}条提交的答案统计`)
    
    // 检查answers是否存在
    if (submission.answers && Array.isArray(submission.answers)) {
      submission.answers.forEach(answer => {
        if (answer && answer.question_id) {
          const questionStat = questionStats.find(q => q.id === answer.question_id)
          if (questionStat) {
            const content = answer.content || '空'
            questionStat.answers.push(content)
            
            // 计算答案频率
            if (!questionStat.answerCounts[content]) {
              questionStat.answerCounts[content] = 0
            }
            questionStat.answerCounts[content]++
            
            console.log(`问题ID=${answer.question_id} 答案="${content}" 当前计数=${questionStat.answerCounts[content]}`)
          } else {
            console.warn(`未找到对应的问题统计，问题ID=${answer.question_id}`)
          }
        } else {
          console.warn('答案数据无效:', answer)
        }
      })
    } else {
      console.warn(`第${index+1}条提交缺少答案数据`)
    }
  })
  
  // 计算用户统计信息
  const uniqueUserIds = new Set(submissions.value.map(s => s.submission?.user_id).filter(id => id !== undefined))
  
  statistics.value = {
    questionStats,
    userStats: {
      total: submissions.value.length,
      uniqueUsers: uniqueUserIds.size
    }
  }
  
  console.log('统计数据生成完成:', statistics.value)
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString()
}

// 获取问题标题
const getQuestionTitle = (questionId) => {
  const question = questions.value.find(q => q.id === questionId)
  return question ? question.title : '未知问题'
}

// 导出结果为CSV
const exportToCSV = () => {
  if (!submissions.value.length || !questions.value.length) {
    Toast('没有数据可导出')
    return
  }
  
  try {
    // 创建CSV内容
    let csvContent = '提交时间,用户名,IP地址,'
    
    // 添加问题标题作为列头
    questions.value.forEach(question => {
      csvContent += `"${question.title || '未命名问题'}",`
    })
    csvContent += '\n'
    
    // 添加每一行数据
    submissions.value.forEach(submission => {
      if (!submission || !submission.submission) return
      
      // 添加基本信息
      csvContent += `"${formatDate(submission.submission.submitted_at)}",`
      csvContent += `"${submission.user_info?.username || '匿名用户'}",`
      csvContent += `"${submission.submission.ip_address || '-'}",`
      
      // 为每个问题添加答案
      questions.value.forEach(question => {
        let answerContent = '-'
        if (Array.isArray(submission.answers)) {
          const answer = submission.answers.find(a => a && a.question_id === question.id)
          if (answer && answer.content !== undefined) {
            answerContent = answer.content
          }
        }
        csvContent += `"${answerContent}",`
      })
      
      csvContent += '\n'
    })
    
    // 创建下载链接
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.setAttribute('href', url)
    link.setAttribute('download', `问卷结果_${questionnaire.value?.title || '未命名'}_${new Date().toISOString().slice(0, 10)}.csv`)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    
    Toast.success('导出成功')
  } catch (error) {
    console.error('导出CSV失败:', error)
    Toast.fail('导出失败: ' + (error.message || '未知错误'))
  }
}

// 渲染图表
const renderCharts = () => {
  // 清除所有现有图表
  Object.values(chartInstances.value).forEach(chart => {
    if (chart && chart.dispose) {
      chart.dispose()
    }
  })
  chartInstances.value = {}
  
  // 延迟执行，确保DOM已更新
  setTimeout(() => {
    statistics.value.questionStats.forEach(question => {
      if (Object.keys(question.answerCounts).length === 0) return
      
      const chartId = `chart-${question.id}`
      const chartDom = document.getElementById(chartId)
      if (!chartDom) {
        console.warn(`未找到图表DOM元素: ${chartId}`)
        return
      }
      
      try {
        // 检查是否已经有图表实例
        if (echarts.getInstanceByDom(chartDom)) {
          console.log(`销毁已存在的图表实例: ${chartId}`)
          echarts.dispose(chartDom)
        }
        
        console.log(`初始化图表: ${chartId}, 问题类型: ${question.type}`)
        const chart = echarts.init(chartDom)
        chartInstances.value[question.id] = chart
        
        // 根据问题类型选择不同图表
        if (question.type === '单选题') {
          renderPieChart(chart, question)
        } else if (question.type === '多选题') {
          renderBarChart(chart, question)
        } else if (question.type === '评分题') {
          renderScoreChart(chart, question)
        }
        
        // 添加窗口大小变化时自动调整图表大小
        window.addEventListener('resize', () => {
          if (chart && !chart.isDisposed()) {
            chart.resize()
          }
        })
      } catch (error) {
        console.error(`渲染图表失败 (问题ID: ${question.id}):`, error)
        Toast.fail(`图表渲染失败: ${error.message || '未知错误'}`)
      }
    })
    
    // 如果是聚焦的问题，滚动到该问题
    if (focusedQuestionId.value) {
      scrollToQuestion(focusedQuestionId.value)
    }
  }, 500) // 增加延迟时间，确保DOM完全准备好
}

// 滚动到指定问题
const scrollToQuestion = (questionId) => {
  setTimeout(() => {
    const element = document.getElementById(`question-${questionId}`)
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' })
      // 高亮显示
      element.classList.add('highlight-question')
      setTimeout(() => {
        element.classList.remove('highlight-question')
      }, 2000)
    }
  }, 300)
}

// 渲染饼图 (单选题)
const renderPieChart = (chart, question) => {
  const data = Object.entries(question.answerCounts).map(([name, value]) => ({ name, value }))
  
  // 计算总数和百分比
  const total = data.reduce((sum, item) => sum + item.value, 0)
  
  // 生成文字描述
  let description = `该问题共收到 ${total} 个回答。\n`
  
  // 获取前三名选项
  const sortedData = [...data].sort((a, b) => b.value - a.value)
  sortedData.slice(0, 3).forEach((item, index) => {
    const percent = ((item.value / total) * 100).toFixed(1)
    description += `${index + 1}. "${item.name}" 选项被选择 ${item.value} 次，占比 ${percent}%。\n`
  })
  
  // 如果有超过3个选项，添加总结信息
  if (data.length > 3) {
    description += `\n共有 ${data.length} 个不同选项。`
  }
  
  const option = {
    title: {
      text: question.title,
      left: 'center',
      subtext: description,
      subtextStyle: {
        align: 'left',
        color: '#666',
        lineHeight: 18
      }
    },
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left',
      top: 'middle',
      data: data.map(item => item.name)
    },
    series: [
      {
        name: question.title,
        type: 'pie',
        radius: '60%',
        center: ['60%', '60%'],
        data: data,
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        },
        label: {
          formatter: '{b}: {c} ({d}%)'
        }
      }
    ]
  }
  
  try {
    chart.setOption(option)
  } catch (error) {
    console.error('设置饼图选项失败:', error)
  }
}

// 渲染柱状图 (多选题)
const renderBarChart = (chart, question) => {
  const data = Object.entries(question.answerCounts).map(([name, value]) => ({ name, value }))
  data.sort((a, b) => b.value - a.value)
  
  // 生成文字描述
  const total = data.reduce((sum, item) => sum + item.value, 0)
  let description = `该多选题共收到 ${total} 个选择。\n`
  
  // 获取前三名选项
  data.slice(0, 3).forEach((item, index) => {
    const percent = ((item.value / total) * 100).toFixed(1)
    description += `${index + 1}. "${item.name}" 被选择 ${item.value} 次，占比 ${percent}%。\n`
  })
  
  const option = {
    title: {
      text: question.title,
      left: 'center',
      subtext: description,
      subtextStyle: {
        align: 'left',
        color: '#666',
        lineHeight: 18
      }
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
      type: 'value'
    },
    yAxis: {
      type: 'category',
      data: data.map(item => item.name),
      axisLabel: {
        interval: 0,
        rotate: 30
      }
    },
    series: [
      {
        name: '选择次数',
        type: 'bar',
        data: data.map(item => item.value)
      }
    ]
  }
  
  chart.setOption(option)
}

// 渲染评分图表
const renderScoreChart = (chart, question) => {
  const scores = Object.keys(question.answerCounts).map(Number).sort((a, b) => a - b)
  const counts = scores.map(score => question.answerCounts[score.toString()])
  
  // 计算平均分
  const total = scores.reduce((sum, score, index) => sum + score * counts[index], 0)
  const totalCount = counts.reduce((sum, count) => sum + count, 0)
  const average = totalCount > 0 ? (total / totalCount).toFixed(1) : 0
  
  // 生成文字描述
  let description = `该评分题共收到 ${totalCount} 个评分，平均分为 ${average}。\n`
  
  // 找出最高频率的分数
  let maxCount = 0
  let mostCommonScore = 0
  counts.forEach((count, index) => {
    if (count > maxCount) {
      maxCount = count
      mostCommonScore = scores[index]
    }
  })
  
  if (maxCount > 0) {
    const percent = ((maxCount / totalCount) * 100).toFixed(1)
    description += `最常见的评分是 ${mostCommonScore} 分，有 ${maxCount} 人选择，占比 ${percent}%。\n`
  }
  
  const option = {
    title: {
      text: `${question.title} (平均分: ${average})`,
      left: 'center',
      subtext: description,
      subtextStyle: {
        align: 'left',
        color: '#666',
        lineHeight: 18
      }
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
    xAxis: [
      {
        type: 'category',
        data: scores.map(String),
        name: '分数'
      }
    ],
    yAxis: [
      {
        type: 'value',
        name: '人数'
      }
    ],
    series: [
      {
        name: '评分分布',
        type: 'bar',
        data: counts,
        markLine: {
          data: [
            { type: 'average', name: '平均值' }
          ]
        }
      }
    ]
  }
  
  chart.setOption(option)
}

// 查看填空题所有回答
const viewTextAnswers = (questionId) => {
  const question = questions.value.find(q => q.id === questionId)
  if (!question) return
  
  const answers = []
  submissions.value.forEach(submission => {
    if (submission.answers) {
      const answer = submission.answers.find(a => a.question_id === questionId)
      if (answer && answer.content) {
        answers.push({
          content: answer.content,
          username: submission.user_info?.username || '匿名用户',
          time: formatDate(submission.submission?.submitted_at)
        })
      }
    }
  })
  
  if (answers.length === 0) {
    Toast('暂无回答')
    return
  }
  
  Dialog.alert({
    title: question.title,
    message: answers.map(a => `${a.username} (${a.time}):\n${a.content}`).join('\n\n'),
    messageAlign: 'left',
    confirmButtonText: '关闭'
  })
}

// 切换折叠面板
const handleCollapseChange = (names) => {
  console.log('折叠面板状态变化:', names, typeof names[0])
  activeNames.value = names
}

onMounted(() => {
  loadQuestionnaireResults()
})

// 监听统计数据变化，重新渲染图表
watch(() => statistics.value, async () => {
  await nextTick()
  renderCharts()
}, { deep: true })

// 监听路由参数变化，重新加载数据
watch(() => route.query.focusQuestion, (newVal) => {
  if (newVal) {
    focusedQuestionId.value = parseInt(newVal)
    if (questions.value.length > 0) {
      scrollToQuestion(focusedQuestionId.value)
    }
  }
})
</script>

<template>
  <div class="results-container">
    <van-nav-bar
      title="问卷结果"
      left-text="返回"
      left-arrow
      @click-left="router.push(`/questionnaire/detail/${route.params.id}`)"
    />
    
    <div class="content" v-if="questionnaire && !loading">
      <div class="questionnaire-header">
        <h1 class="title">{{ questionnaire.title }}</h1>
        <p class="description">{{ questionnaire.description }}</p>
        <div class="meta">
          <van-tag type="primary" v-if="questionnaire.is_published">已发布</van-tag>
          <van-tag type="warning" v-else>未发布</van-tag>
          <div class="time">
            <p>开始时间: {{ formatDate(questionnaire.start_time) }}</p>
            <p>结束时间: {{ formatDate(questionnaire.end_time) }}</p>
          </div>
        </div>
      </div>
      
      <div class="summary">
        <van-cell-group inset title="统计信息">
          <van-cell title="总提交数" :value="totalSubmissions + ' 份'" />
          <van-cell title="独立用户数" :value="statistics.userStats.uniqueUsers + ' 人'" />
        </van-cell-group>
        
        <div class="export-button">
          <van-button type="primary" block @click="exportToCSV">导出结果为CSV</van-button>
        </div>
      </div>
      
      <!-- 问题统计结果 -->
      <van-divider>问题统计</van-divider>
      
      <div class="question-stats" v-if="statistics.questionStats.length > 0">
        <van-cell-group 
          v-for="question in statistics.questionStats" 
          :key="question.id"
          inset
          :title="`${question.title} (${question.type})`"
          :id="`question-${question.id}`"
          class="question-stat-group"
        >
          <template v-if="Object.keys(question.answerCounts || {}).length > 0 && question.answers && question.answers.length > 0">
            <!-- 图表容器 -->
            <div v-if="question.type === '单选题' || question.type === '多选题' || question.type === '评分题'" 
                 :id="`chart-${question.id}`" 
                 class="chart-container">
            </div>
            
            <!-- 填空题查看按钮 -->
            <van-cell v-else-if="question.type === '填空题'" center>
              <template #title>
                <van-button type="primary" size="small" @click="viewTextAnswers(question.id)">
                  查看所有回答 ({{ question.answers.length }}条)
                </van-button>
              </template>
            </van-cell>
            
            <!-- 数据表格 -->
            <van-cell v-for="(count, answer) in question.answerCounts" :key="answer">
              <template #title>
                <span class="answer-text">{{ answer }}</span>
              </template>
              <template #value>
                <div class="answer-stat">
                  <div class="progress-bar">
                    <div 
                      class="progress" 
                      :style="{width: `${(count / (question.answers.length || 1)) * 100}%`}"
                    ></div>
                  </div>
                  <span class="count">{{ count }}次 ({{ Math.round((count / (question.answers.length || 1)) * 100) }}%)</span>
                </div>
              </template>
            </van-cell>
          </template>
          <template v-else>
            <van-cell title="暂无回答" />
          </template>
        </van-cell-group>
      </div>
      
      <van-divider>提交列表</van-divider>
      
      <div class="submissions-list" v-if="submissions.length > 0">
        <van-collapse v-model="activeNames" @change="handleCollapseChange">
          <van-collapse-item 
            v-for="(submission, index) in submissions" 
            :key="submission.submission?.id || index"
            :title="`提交 #${index + 1} - ${formatDate(submission.submission?.submitted_at)}`"
            :name="index"
          >
            <van-cell-group inset title="提交信息">
              <van-cell title="提交时间" :value="formatDate(submission.submission?.submitted_at)" />
              <van-cell title="提交用户" :value="submission.user_info?.username || '匿名用户'" />
              <van-cell title="IP地址" :value="submission.submission?.ip_address || '-'" />
            </van-cell-group>
            
            <van-cell-group inset title="回答内容">
              <template v-if="Array.isArray(submission.answers) && submission.answers.length > 0">
                <template v-for="answer in submission.answers" :key="answer?.id || Math.random()">
                  <van-cell v-if="answer">
                    <template #title>
                      <div class="question-title">{{ getQuestionTitle(answer.question_id) }}</div>
                    </template>
                    <template #value>
                      <div class="answer-content">{{ answer.content || '-' }}</div>
                    </template>
                  </van-cell>
                </template>
              </template>
              <template v-else>
                <van-cell title="暂无回答数据" />
              </template>
            </van-cell-group>
          </van-collapse-item>
        </van-collapse>
      </div>
      
      <div class="empty-state" v-else>
        <van-empty description="暂无提交记录" />
      </div>
    </div>
    
    <div class="loading" v-else-if="loading">
      <van-loading type="spinner" />
    </div>
    
    <div class="error" v-else>
      <van-empty description="加载失败，请重试" />
    </div>
  </div>
</template>

<style scoped>
.results-container {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.content {
  padding: 16px;
}

.questionnaire-header {
  background-color: #fff;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 16px;
}

.title {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 8px;
}

.description {
  font-size: 14px;
  color: #666;
  margin-bottom: 16px;
}

.meta {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.time {
  font-size: 12px;
  color: #999;
  margin-top: 8px;
}

.summary {
  margin-bottom: 24px;
}

.export-button {
  margin-top: 16px;
}

.question-stats {
  margin-bottom: 24px;
}

.question-stat-group {
  margin-bottom: 16px;
}

.answer-stat {
  display: flex;
  align-items: center;
  gap: 8px;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background-color: #eee;
  border-radius: 4px;
  overflow: hidden;
}

.progress {
  height: 100%;
  background-color: #1989fa;
  border-radius: 4px;
}

.count {
  font-size: 12px;
  color: #666;
  white-space: nowrap;
}

.answer-text {
  font-size: 14px;
  color: #333;
}

.submissions-list {
  margin-bottom: 24px;
}

.question-title {
  font-weight: bold;
  color: #333;
}

.answer-content {
  word-break: break-all;
}

.loading, .error {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

.chart-container {
  height: 300px;
  margin: 16px 0;
}

.empty-state {
  padding: 24px 0;
}

/* 高亮样式 */
.highlight-question {
  animation: highlight-animation 2s ease;
}

@keyframes highlight-animation {
  0% { background-color: rgba(24, 144, 255, 0.1); }
  50% { background-color: rgba(24, 144, 255, 0.2); }
  100% { background-color: transparent; }
}
</style> 