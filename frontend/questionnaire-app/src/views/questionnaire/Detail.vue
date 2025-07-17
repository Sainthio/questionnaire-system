<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuestionnaireStore } from '../../stores/questionnaire'
import { useUserStore } from '../../stores/user'
import { Toast, Dialog } from 'vant'

const route = useRoute()
const router = useRouter()
const questionnaireStore = useQuestionnaireStore()
const userStore = useUserStore()

const loading = ref(false)
const questionnaire = ref(null)
const questions = ref([])
const hasSubmitted = ref(false)
const activeQuestions = ref([]) // 用于存储当前展开的问题

// 计算属性：判断当前用户是否有权限查看/填写问卷
const hasPermission = computed(() => {
  if (!questionnaire.value) return false
  
  // 如果问卷已发布，所有人都可以查看
  if (questionnaire.value.is_published) return true
  
  // 如果是管理员，可以查看所有问卷
  if (userStore.isAdmin) return true
  
  // 如果是问卷创建者，可以查看自己的问卷
  return userStore.isLoggedIn && userStore.userInfo.id === questionnaire.value.created_by
})

// 计算属性：判断当前用户是否可以编辑问卷
const canEdit = computed(() => {
  if (!questionnaire.value) return false
  
  // 只有未发布的问卷可以编辑
  if (questionnaire.value.is_published) return false
  
  // 管理员或创建者可以编辑
  return userStore.isAdmin || (userStore.isLoggedIn && userStore.userInfo.id === questionnaire.value.created_by)
})

// 计算属性：判断当前用户是否可以查看问卷结果
const canViewResults = computed(() => {
  if (!questionnaire.value) return false
  
  // 管理员或创建者可以查看结果
  return userStore.isAdmin || (userStore.isLoggedIn && userStore.userInfo.id === questionnaire.value.created_by)
})

// 检查用户是否已提交过问卷
const checkUserSubmission = async () => {
  if (!userStore.isLoggedIn || !questionnaire.value) return
  
  try {
    const response = await questionnaireStore.checkSubmission(
      questionnaire.value.id, 
      userStore.userInfo.id
    )
    
    if (response && response.has_submitted) {
      hasSubmitted.value = true
      console.log('用户已提交过该问卷')
    } else {
      hasSubmitted.value = false
      console.log('用户未提交过该问卷')
    }
  } catch (error) {
    console.error('检查提交状态失败:', error)
    hasSubmitted.value = false
  }
}

// 获取问卷详情
const loadQuestionnaireDetail = async () => {
  const id = route.params.id
  if (!id) {
    Toast('问卷ID不能为空')
    router.push('/questionnaire/list')
    return
  }
  
  loading.value = true
  Toast.loading({
    message: '加载中...',
    forbidClick: true,
    duration: 0
  })
  
  try {
    const response = await questionnaireStore.getQuestionnaireDetail(id)
    console.log('问卷详情数据:', response)
    
    // 检查响应数据结构
    if (!response || !response.data) {
      throw new Error('返回数据格式错误')
    }
    
    questionnaire.value = response.data.questionnaire || {}
    questions.value = response.data.questions || []
    
    // 检查用户是否已提交过该问卷
    await checkUserSubmission()
    
    // 如果问卷为空，显示错误
    if (!questionnaire.value.id) {
      throw new Error('问卷不存在或已被删除')
    }
    
    // 检查权限
    if (!hasPermission.value) {
      Dialog.alert({
        title: '无权访问',
        message: '该问卷尚未发布，您没有权限查看。',
      }).then(() => {
        router.push('/questionnaire/list')
      })
      return
    }
  } catch (error) {
    console.error('加载问卷详情失败:', error)
    Toast.clear()
    Toast.fail('加载失败: ' + (error.message || '未知错误'))
  } finally {
    loading.value = false
    Toast.clear()
  }
}

// 填写问卷
const fillQuestionnaire = () => {
  if (hasSubmitted.value) {
    Dialog.alert({
      title: '已提交',
      message: '您已经提交过该问卷，不能重复提交',
    })
    return
  }
  
  router.push(`/questionnaire/fill/${route.params.id}`)
}

// 编辑问卷
const editQuestionnaire = () => {
  if (!canEdit.value) {
    Dialog.alert({
      title: '无法编辑',
      message: '您没有权限编辑此问卷或问卷已发布。',
    })
    return
  }
  
  // 将问卷数据存储到store中，以便编辑页面使用
  questionnaireStore.currentQuestionnaire = questionnaire.value
  
  // 导航到编辑页面
  router.push(`/questionnaire/edit/${route.params.id}`)
}

// 发布问卷
const publishQuestionnaire = async () => {
  if (!canEdit.value) {
    Dialog.alert({
      title: '无法发布',
      message: '您没有权限发布此问卷。',
    })
    return
  }
  
  Dialog.confirm({
    title: '发布问卷',
    message: '发布后问卷将对所有用户可见，且不可再编辑。确定发布吗？',
  }).then(async () => {
    try {
      Toast.loading({
        message: '发布中...',
        forbidClick: true,
        duration: 0
      })
      
      // 确保ID是数字类型
      const questionnaireId = parseInt(route.params.id, 10)
      console.log('准备发布问卷:', { id: questionnaireId, is_published: true })
      
      // 检查用户登录状态
      if (!userStore.isLoggedIn) {
        throw new Error('您需要先登录才能发布问卷')
      }
      
      await questionnaireStore.updateQuestionnaireStatus(questionnaireId, true)
      
      Toast.success('发布成功')
      
      // 重新加载问卷详情
      loadQuestionnaireDetail()
    } catch (error) {
      console.error('发布问卷失败:', error)
      Toast.fail('发布失败: ' + (error.message || '未知错误'))
    } finally {
      Toast.clear()
    }
  }).catch(() => {
    // 用户取消操作
  })
}

// 查看结果
const viewResults = () => {
  if (!canViewResults.value) {
    Dialog.alert({
      title: '无法查看结果',
      message: '您没有权限查看此问卷的结果。',
    })
    return
  }
  
  router.push(`/questionnaire/results/${route.params.id}`)
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString()
}

// 处理问题折叠面板变化
const handleQuestionChange = (names) => {
  console.log('问题折叠面板变化:', names)
  activeQuestions.value = names
}

// 查看问题统计结果
const viewQuestionResults = (questionId) => {
  if (!canViewResults.value) {
    Dialog.alert({
      title: '无法查看结果',
      message: '您没有权限查看此问卷的结果。',
    })
    return
  }
  
  // 直接跳转到问卷结果页面并聚焦到特定问题
  router.push(`/questionnaire/results/${route.params.id}?focusQuestion=${questionId}`)
}

// 跳转到问卷管理页面（仅管理员）
const goToQuestionnaireManagement = () => {
  router.push('/questionnaire/list')
}

onMounted(() => {
  loadQuestionnaireDetail()
  // 默认展开第一个问题
  activeQuestions.value = [0]
})
</script>

<template>
  <div class="detail-container">
    <van-nav-bar
      title="问卷详情"
      left-text="返回"
      left-arrow
      @click-left="router.push('/questionnaire/list')"
    />
    
    <div class="content" v-if="questionnaire">
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
        
        <div class="action-buttons">
          <!-- 管理问卷按钮 - 仅管理员可见 -->
          <van-button
            v-if="userStore.isAdmin"
            type="primary"
            block
            @click="goToQuestionnaireManagement"
            class="manage-button"
          >
            管理问卷
          </van-button>
          
          <!-- 编辑按钮 - 仅管理员和创建者可见，且问卷未发布 -->
          <van-button 
            v-if="canEdit && !questionnaire.is_published" 
            type="primary" 
            block 
            @click="editQuestionnaire"
          >
            编辑问卷
          </van-button>
          
          <!-- 发布/取消发布按钮 - 仅管理员和创建者可见 -->
          <van-button 
            v-if="canEdit" 
            :type="questionnaire.is_published ? 'warning' : 'success'" 
            block 
            @click="publishQuestionnaire"
            class="publish-button"
          >
            {{ questionnaire.is_published ? '取消发布' : '发布问卷' }}
          </van-button>
          
          <!-- 查看结果按钮 - 仅管理员和创建者可见 -->
          <van-button 
            v-if="canViewResults"
            type="info" 
            block 
            @click="viewResults"
            class="results-button"
          >
            查看结果
          </van-button>
          
          <!-- 填写问卷按钮 - 已发布问卷所有用户可见 -->
          <van-button 
            v-if="questionnaire.is_published && userStore.isLoggedIn" 
            type="primary" 
            block 
            @click="fillQuestionnaire"
            :disabled="hasSubmitted"
            class="fill-button"
          >
            {{ hasSubmitted ? '已填写' : '填写问卷' }}
          </van-button>
        </div>
      </div>
      
      <van-divider>问题列表</van-divider>
      
      <div class="questions-list">
        <template v-if="questions.length > 0">
          <van-collapse v-model="activeQuestions" @change="handleQuestionChange">
            <van-collapse-item
              v-for="(question, index) in questions"
              :key="question.id"
              :title="`${index + 1}. ${question.title}`"
              :name="index"
            >
              <div class="question-content">
                <div class="question-info">
                  <van-tag type="primary">{{ question.type }}</van-tag>
                  <van-tag type="danger" v-if="question.required">必填</van-tag>
                </div>
                <div class="question-options" v-if="question.options">
                  <p>选项：{{ question.options }}</p>
                </div>
                <!-- 添加查看结果按钮 -->
                <div class="question-actions" v-if="questionnaire.is_published && canViewResults">
                  <van-button 
                    type="primary" 
                    size="small" 
                    @click.stop="viewQuestionResults(question.id)"
                  >
                    查看统计结果
                  </van-button>
                </div>
              </div>
            </van-collapse-item>
          </van-collapse>
        </template>
        <template v-else>
          <van-empty description="暂无问题" />
        </template>
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
.detail-container {
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

.questions-list {
  margin-bottom: 24px;
}

.question-content {
  padding: 8px 0;
}

.question-info {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.question-options {
  font-size: 14px;
  color: #666;
  margin-top: 8px;
}

.action-buttons {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.publish-button,
.results-button,
.fill-button {
  margin-top: 10px;
}

.loading, .error {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

.van-button--disabled {
  opacity: 0.6;
}

.question-actions {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}
</style> 