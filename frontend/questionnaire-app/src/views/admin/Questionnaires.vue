<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'
import { questionnaireApi } from '../../api/admin'
import { Toast, Dialog } from 'vant'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(true)
const questionnaires = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 问卷详情相关
const showSubmissionDetail = ref(false)
const currentQuestionnaire = ref(null)
const submissionLoading = ref(false)
const activeSubmissions = ref([]) // 添加折叠面板控制变量

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
  
  loadQuestionnaires()
})

// 加载问卷列表
const loadQuestionnaires = async () => {
  loading.value = true
  try {
    const response = await questionnaireApi.getAllQuestionnaires(currentPage.value, pageSize.value)
    questionnaires.value = response.data.data.questionnaires
    total.value = response.data.data.total
    console.log('问卷列表:', questionnaires.value)
  } catch (error) {
    console.error('获取问卷列表失败:', error)
    Toast('获取问卷列表失败')
  } finally {
    loading.value = false
  }
}

// 查看问卷提交详情
const viewSubmissionDetail = async (questionnaireId) => {
  submissionLoading.value = true
  showSubmissionDetail.value = true
  
  try {
    const response = await questionnaireApi.getQuestionnaireSubmissions(questionnaireId)
    currentQuestionnaire.value = response.data.data
    
    // 确保submission_details存在
    if (!currentQuestionnaire.value.submission_details) {
      currentQuestionnaire.value.submission_details = []
    }
    
    console.log('问卷提交详情:', currentQuestionnaire.value)
  } catch (error) {
    console.error('获取问卷提交详情失败:', error)
    Toast('获取问卷提交详情失败')
    showSubmissionDetail.value = false
  } finally {
    submissionLoading.value = false
  }
}

// 查看问卷详情
const viewQuestionnaireDetail = (id) => {
  // 跳转到问卷详情页
  router.push(`/questionnaire/detail/${id}`)
}

// 管理问卷（返回问卷列表）
const manageQuestionnaires = () => {
  // 从问卷详情返回问卷列表
  router.push('/admin/questionnaires')
}

// 页面变化
const onPageChange = (page) => {
  currentPage.value = page
  loadQuestionnaires()
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString()
}

// 获取问题的选项
const getQuestionOptions = (optionsJson) => {
  try {
    return JSON.parse(optionsJson)
  } catch (e) {
    return []
  }
}

// 获取问题类型中文名称
const getQuestionTypeName = (type) => {
  const typeMap = {
    'single': '单选题',
    'multiple': '多选题',
    'text': '填空题',
    'rating': '评分题'
  }
  return typeMap[type] || type
}

// 格式化答案内容
const formatAnswerContent = (content, questionType, optionsJson) => {
  if (!content) return '-'
  
  if (questionType === 'single') {
    // 单选题，显示选项文本
    try {
      const options = JSON.parse(optionsJson)
      const selectedOption = options.find(opt => opt.value === content)
      return selectedOption ? selectedOption.text : content
    } catch (e) {
      return content
    }
  } else if (questionType === 'multiple') {
    // 多选题，显示多个选项文本
    try {
      const options = JSON.parse(optionsJson)
      const selectedValues = content.split(',')
      const selectedTexts = selectedValues.map(value => {
        const option = options.find(opt => opt.value === value)
        return option ? option.text : value
      })
      return selectedTexts.join(', ')
    } catch (e) {
      return content
    }
  }
  
  return content
}
</script>

<template>
  <div class="admin-questionnaires">
    <van-nav-bar
      title="问卷列表管理"
      left-text="返回"
      left-arrow
      @click-left="router.push('/admin')"
    />
    
    <div class="questionnaires-container">
      <van-loading v-if="loading" size="24px" vertical>加载中...</van-loading>
      
      <template v-else>
        <div class="table-header">
          <h2>问卷列表 ({{ total }})</h2>
        </div>
        
        <div class="questionnaires-table">
          <van-cell-group inset>
            <!-- 表头 -->
            <van-cell title-class="table-title" value-class="table-title">
              <template #title>
                <div class="cell-title">问卷标题</div>
              </template>
              <template #value>
                <div class="cell-value">
                  <span class="creator-column">创建者</span>
                  <span class="status-column">状态</span>
                  <span class="count-column">问题/提交</span>
                  <span class="action-column">操作</span>
                </div>
              </template>
            </van-cell>
            
            <!-- 表格内容 -->
            <van-cell v-for="item in questionnaires" :key="item.questionnaire.id" clickable @click="viewQuestionnaireDetail(item.questionnaire.id)">
              <template #title>
                <div class="cell-title">
                  {{ item.questionnaire.title }}
                </div>
              </template>
              <template #value>
                <div class="cell-value">
                  <span class="creator-column">{{ item.creator_name || '未知' }}</span>
                  <span class="status-column">
                    <van-tag type="primary" v-if="item.questionnaire.is_published">已发布</van-tag>
                    <van-tag type="warning" v-else>未发布</van-tag>
                  </span>
                  <span class="count-column">
                    {{ item.question_count }}/{{ item.submission_count }}
                  </span>
                  <span class="action-column">
                    <van-button size="mini" type="primary" @click.stop="viewSubmissionDetail(item.questionnaire.id)">查看提交</van-button>
                  </span>
                </div>
              </template>
            </van-cell>
          </van-cell-group>
          
          <!-- 分页 -->
          <div class="pagination">
            <van-pagination
              v-model="currentPage"
              :total-items="total"
              :items-per-page="pageSize"
              @change="onPageChange"
            />
          </div>
        </div>
      </template>
    </div>
    
    <!-- 问卷提交详情弹窗 -->
    <van-popup v-model:show="showSubmissionDetail" round position="bottom" :style="{ height: '90%' }">
      <div class="submission-detail">
        <van-nav-bar
          title="问卷提交统计详情"
          left-text="返回"
          left-arrow
          @click-left="showSubmissionDetail = false"
        />
        
        <van-loading v-if="submissionLoading" size="24px" vertical>加载中...</van-loading>
        
        <template v-else-if="currentQuestionnaire">
          <div class="detail-content">
            <!-- 问卷基本信息 -->
            <div class="questionnaire-info">
              <h3>{{ currentQuestionnaire.questionnaire?.title || '无标题' }}</h3>
              <p class="description">{{ currentQuestionnaire.questionnaire?.description || '无描述' }}</p>
              <div class="meta-info">
                <van-tag type="primary" v-if="currentQuestionnaire.questionnaire?.is_published">已发布</van-tag>
                <van-tag type="warning" v-else>未发布</van-tag>
                <span class="date">开始时间: {{ formatDate(currentQuestionnaire.questionnaire?.start_time) }}</span>
                <span class="date">结束时间: {{ formatDate(currentQuestionnaire.questionnaire?.end_time) }}</span>
              </div>
              
              <!-- 添加"管理问卷"按钮 -->
              <div class="action-buttons">
                <van-button type="primary" size="small" @click="manageQuestionnaires">管理问卷</van-button>
                <van-button type="info" size="small" @click="viewQuestionnaireDetail(currentQuestionnaire.questionnaire?.id)">问卷详情</van-button>
              </div>
            </div>
            
            <!-- 提交统计 -->
            <van-cell-group inset title="提交统计">
              <van-cell title="总提交数" :value="currentQuestionnaire.submission_details?.length || 0" />
            </van-cell-group>
            
            <!-- 问题列表 -->
            <van-cell-group inset title="问题列表" v-if="currentQuestionnaire.questions?.length">
              <van-cell v-for="question in currentQuestionnaire.questions" :key="question.id">
                <template #title>
                  <div class="question-title">
                    <span class="question-index">Q{{ question.sort + 1 }}</span>
                    <span>{{ question.title }}</span>
                    <van-tag size="small" :type="question.required ? 'danger' : 'default'">
                      {{ question.required ? '必填' : '选填' }}
                    </van-tag>
                    <van-tag size="small" type="primary">{{ getQuestionTypeName(question.type) }}</van-tag>
                  </div>
                  
                  <div v-if="question.type !== 'text'" class="question-options">
                    <div v-for="(option, index) in getQuestionOptions(question.options)" :key="index" class="option-item">
                      {{ option.text }}
                    </div>
                  </div>
                </template>
              </van-cell>
            </van-cell-group>
            
            <!-- 提交详情 -->
            <div v-if="currentQuestionnaire.submission_details && currentQuestionnaire.submission_details.length > 0">
              <h4 class="submission-title">提交详情 ({{ currentQuestionnaire.submission_details.length }})</h4>
              
              <van-collapse v-model="activeSubmissions">
                <van-collapse-item 
                  v-for="(submission, index) in currentQuestionnaire.submission_details" 
                  :key="submission.submission.id"
                  :title="`提交 #${index + 1} - ${submission.user?.username || '未知用户'}`"
                  :name="submission.submission.id"
                >
                  <div class="submission-info">
                    <p>提交时间: {{ formatDate(submission.submission?.submitted_at) }}</p>
                    <p>IP地址: {{ submission.submission?.ip_address || '未知' }}</p>
                  </div>
                  
                  <van-cell-group title="回答详情" v-if="submission.answers && submission.answers.length > 0">
                    <van-cell v-for="question in currentQuestionnaire.questions || []" :key="question.id">
                      <template #title>
                        <div class="answer-question">
                          <span class="question-index">Q{{ question.sort + 1 }}</span>
                          <span>{{ question.title }}</span>
                        </div>
                      </template>
                      <template #value>
                        <div class="answer-content">
                          {{ 
                            formatAnswerContent(
                              submission.answers.find(a => a && a.question_id === question.id)?.content, 
                              question.type,
                              question.options
                            ) 
                          }}
                        </div>
                      </template>
                    </van-cell>
                  </van-cell-group>
                  <div v-else class="empty-answers">
                    <p>暂无回答记录</p>
                  </div>
                </van-collapse-item>
              </van-collapse>
            </div>
            
            <div v-else class="empty-submissions">
              <van-empty description="暂无提交记录" />
            </div>
          </div>
        </template>
        
        <div v-else class="error-message">
          <van-empty description="无法加载问卷数据">
            <template #bottom>
              <van-button round type="primary" @click="showSubmissionDetail = false">返回</van-button>
            </template>
          </van-empty>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<style scoped>
.admin-questionnaires {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.questionnaires-container {
  padding: 16px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.table-header h2 {
  font-size: 18px;
  margin: 0;
}

.questionnaires-table {
  margin-bottom: 20px;
}

.table-title {
  font-weight: bold;
  color: #323233;
}

.cell-title {
  font-weight: bold;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.cell-value {
  display: flex;
  align-items: center;
}

.creator-column {
  width: 80px;
  margin-right: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-column {
  width: 60px;
  text-align: center;
  margin-right: 10px;
}

.count-column {
  width: 70px;
  text-align: center;
  margin-right: 10px;
}

.action-column {
  width: 90px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.submission-detail {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.detail-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.questionnaire-info {
  background-color: #fff;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 16px;
}

.questionnaire-info h3 {
  margin: 0 0 8px;
  font-size: 18px;
}

.description {
  color: #646566;
  margin-bottom: 12px;
}

.meta-info {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  margin-bottom: 12px;
}

.date {
  color: #969799;
  font-size: 14px;
  margin-left: 8px;
}

.action-buttons {
  display: flex;
  gap: 12px;
  margin-top: 12px;
}

.question-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.question-index {
  font-weight: bold;
  color: #1989fa;
  margin-right: 4px;
}

.question-options {
  margin-top: 8px;
  padding-left: 24px;
}

.option-item {
  margin-bottom: 4px;
  color: #646566;
}

.submission-title {
  margin: 20px 0 12px;
  font-size: 16px;
  color: #323233;
}

.submission-info {
  margin-bottom: 12px;
  color: #646566;
}

.submission-info p {
  margin: 4px 0;
}

.answer-question {
  display: flex;
  align-items: center;
  gap: 8px;
}

.answer-content {
  color: #1989fa;
  font-weight: 500;
}

.empty-submissions {
  margin-top: 40px;
}

.empty-answers {
  padding: 16px;
  text-align: center;
  color: #969799;
}

.error-message {
  padding: 40px 16px;
}
</style> 