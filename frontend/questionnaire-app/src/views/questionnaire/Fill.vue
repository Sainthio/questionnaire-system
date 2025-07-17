<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuestionnaireStore } from '../../stores/questionnaire'
import { useUserStore } from '../../stores/user'
import { Toast, Dialog } from 'vant'
import { api } from '../../services/api'

const route = useRoute()
const router = useRouter()
const questionnaireStore = useQuestionnaireStore()
const userStore = useUserStore()

const loading = ref(false)
const submitting = ref(false)
const questionnaire = ref(null)
const questions = ref([])

// 用户答案
const answers = reactive({})

// 计算属性：判断当前用户是否有权限填写问卷
const canFill = computed(() => {
  if (!questionnaire.value) return false
  
  // 只有已发布的问卷可以填写
  return questionnaire.value.is_published
})

// 获取问卷详情
const loadQuestionnaire = async () => {
  const id = route.params.id
  if (!id) {
    Toast('问卷ID不能为空')
    router.push('/questionnaire/list')
    return
  }
  
  // 检查用户是否登录
  if (!userStore.isLoggedIn) {
    Dialog.confirm({
      title: '未登录',
      message: '登录后才能填写问卷，是否前往登录？',
      confirmButtonText: '去登录',
      cancelButtonText: '返回列表'
    }).then(() => {
      // 保存当前页面URL，登录后可以返回
      localStorage.setItem('redirectUrl', router.currentRoute.value.fullPath)
      router.push('/login')
    }).catch(() => {
      router.push('/questionnaire/list')
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
    // 检查用户是否已提交过该问卷
    try {
      // 使用store中的方法检查用户是否已提交过该问卷
      const userInfo = userStore.userInfo
      if (userInfo && userInfo.id) {
        const checkResponse = await questionnaireStore.checkSubmission(id, userInfo.id)
        if (checkResponse && checkResponse.has_submitted) {
          Toast.clear()
          Dialog.alert({
            title: '已提交',
            message: '您已经提交过该问卷，不能重复提交',
          }).then(() => {
            router.push('/questionnaire/list')
          })
          return
        }
      }
    } catch (error) {
      // 如果API不存在或发生错误，继续执行
      console.error('检查提交状态失败:', error)
    }
    
    const response = await questionnaireStore.getQuestionnaireDetail(id)
    questionnaire.value = response.data.questionnaire
    questions.value = response.data.questions
    
    // 初始化答案对象
    questions.value.forEach(q => {
      // 根据问题类型初始化不同的默认值
      if (q.type === '多选题') {
        answers[q.id] = []
      } else if (q.type === '评分题') {
        answers[q.id] = 0
      } else {
        answers[q.id] = ''
      }
      console.log(`初始化问题 ID=${q.id}, 类型=${q.type}, 默认值=`, answers[q.id])
    })
    
    // 检查问卷是否已发布
    if (!questionnaire.value.is_published) {
      Toast.fail('该问卷尚未发布')
      router.push('/questionnaire/list')
      return
    }
    
    // 检查问卷是否在有效期内
    const now = new Date()
    const startTime = new Date(questionnaire.value.start_time)
    const endTime = new Date(questionnaire.value.end_time)
    
    if (now < startTime) {
      Toast.fail('该问卷尚未开始')
      router.push('/questionnaire/list')
      return
    }
    
    if (now > endTime) {
      Toast.fail('该问卷已结束')
      router.push('/questionnaire/list')
      return
    }
  } catch (error) {
    console.error('加载问卷失败:', error)
    Toast.fail('加载失败: ' + (error.message || '未知错误'))
    router.push('/questionnaire/list')
  } finally {
    loading.value = false
    Toast.clear()
  }
}

// 处理选项变化
const handleOptionChange = (questionId, value, type) => {
  if (type === '多选题') {
    if (!answers[questionId]) {
      answers[questionId] = []
    }
    
    const index = answers[questionId].indexOf(value)
    if (index > -1) {
      answers[questionId].splice(index, 1)
    } else {
      answers[questionId].push(value)
    }
  } else {
    answers[questionId] = value
  }
}

// 验证必填项
const validateRequired = () => {
  const requiredQuestions = questions.value.filter(q => q.required)
  for (const question of requiredQuestions) {
    const answer = answers[question.id]
    if (!answer || (Array.isArray(answer) && answer.length === 0)) {
      Toast(`请回答必填问题: ${question.title}`)
      return false
    }
  }
  return true
}

// 提交问卷
const submitQuestionnaire = async () => {
  try {
    // 验证必填问题是否已回答
    const unansweredRequired = questions.value.filter(q => q.required && !answers[q.id])
    if (unansweredRequired.length > 0) {
      Toast.fail(`请回答所有必填问题 (${unansweredRequired.length}个未回答)`)
      return
    }
    
    // 准备提交数据
    const answerData = {
      questionnaire_id: parseInt(route.params.id),
      user_id: userStore.userInfo.id,
      answers: []
    }
    
    // 转换答案格式
    for (const [questionId, answer] of Object.entries(answers)) {
      if (answer !== '' && answer !== null && answer !== undefined) {
        answerData.answers.push({
          questionnaire_id: parseInt(route.params.id),
          question_id: parseInt(questionId),
          content: Array.isArray(answer) ? JSON.stringify(answer) : String(answer)
        })
      }
    }
    
    submitting.value = true
    Toast.loading({
      message: '提交中...',
      forbidClick: true,
      duration: 0
    })
    
    console.log('提交数据:', answerData)
    await questionnaireStore.submitAnswer(answerData)
    
    Toast.success('提交成功')
    
    // 提交成功后显示成功页面
    Dialog.alert({
      title: '提交成功',
      message: '感谢您的参与！',
    }).then(() => {
      router.push('/questionnaire/list')
    })
  } catch (error) {
    console.error('提交问卷失败:', error)
    
    // 如果是重复提交错误，显示特定提示
    if (error.message && error.message.includes('已经提交过')) {
      Dialog.alert({
        title: '重复提交',
        message: '您已经提交过该问卷，不能重复提交',
      }).then(() => {
        router.push('/questionnaire/list')
      })
    } else {
      Toast.fail('提交失败: ' + error.message)
    }
  } finally {
    submitting.value = false
    Toast.clear()
  }
}

// 解析问题选项
const parseOptions = (optionsString) => {
  try {
    return JSON.parse(optionsString)
  } catch (e) {
    return optionsString.split(',').map(item => item.trim())
  }
}

onMounted(() => {
  loadQuestionnaire()
})
</script>

<template>
  <div class="fill-container">
    <van-nav-bar
      title="填写问卷"
      left-text="返回"
      left-arrow
      @click-left="router.push(`/questionnaire/detail/${route.params.id}`)"
    />
    
    <div class="content" v-if="questionnaire">
      <div class="questionnaire-header">
        <h1 class="title">{{ questionnaire.title }}</h1>
        <p class="description">{{ questionnaire.description }}</p>
      </div>
      
      <div class="questions-list">
        <template v-if="questions.length > 0">
          <div 
            v-for="(question, index) in questions" 
            :key="question.id"
            class="question-item"
          >
            <div class="question-title">
              <span class="question-index">{{ index + 1 }}.</span>
              {{ question.title }}
              <span class="required" v-if="question.required">*</span>
            </div>
            
            <!-- 单选题 -->
            <template v-if="question.type === '单选题'">
              <van-radio-group v-model="answers[question.id]">
                <van-cell-group inset>
                  <van-cell
                    v-for="(option, optIndex) in parseOptions(question.options)"
                    :key="optIndex"
                    :title="option"
                    clickable
                    @click="answers[question.id] = option"
                  >
                    <template #right-icon>
                      <van-radio :name="option" />
                    </template>
                  </van-cell>
                </van-cell-group>
              </van-radio-group>
            </template>
            
            <!-- 多选题 -->
            <template v-else-if="question.type === '多选题'">
              <van-checkbox-group v-model="answers[question.id]">
                <van-cell-group inset>
                  <van-cell
                    v-for="(option, optIndex) in parseOptions(question.options)"
                    :key="optIndex"
                    :title="option"
                    clickable
                    @click="handleOptionChange(question.id, option, question.type)"
                  >
                    <template #right-icon>
                      <van-checkbox :name="option" />
                    </template>
                  </van-cell>
                </van-cell-group>
              </van-checkbox-group>
            </template>
            
            <!-- 填空题 -->
            <template v-else-if="question.type === '填空题'">
              <van-field
                v-model="answers[question.id]"
                type="textarea"
                rows="3"
                placeholder="请输入您的回答"
                :rules="question.required ? [{ required: true, message: '请填写回答' }] : []"
              />
            </template>
            
            <!-- 评分题 -->
            <template v-else-if="question.type === '评分题'">
              <van-rate v-model="answers[question.id]" :count="5" />
            </template>
            
            <!-- 其他类型 -->
            <template v-else>
              <van-field
                v-model="answers[question.id]"
                placeholder="请输入您的回答"
                :rules="question.required ? [{ required: true, message: '请填写回答' }] : []"
              />
            </template>
          </div>
        </template>
        <template v-else>
          <van-empty description="暂无问题" />
        </template>
      </div>
      
      <div class="submit-area">
        <van-button 
          type="primary" 
          block 
          round 
          @click="submitQuestionnaire"
          :loading="submitting"
          :disabled="questions.length === 0"
        >
          提交问卷
        </van-button>
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
.fill-container {
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
}

.questions-list {
  margin-bottom: 24px;
}

.question-item {
  background-color: #fff;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 16px;
}

.question-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 16px;
}

.question-index {
  margin-right: 4px;
}

.required {
  color: #ee0a24;
  margin-left: 4px;
}

.submit-area {
  margin-top: 24px;
  padding: 0 16px 32px;
}

.loading, .error {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}
</style> 