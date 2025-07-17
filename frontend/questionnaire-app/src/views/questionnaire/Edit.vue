<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuestionnaireStore } from '../../stores/questionnaire'
import { useUserStore } from '../../stores/user'
import { Toast, Dialog } from 'vant'

const route = useRoute()
const router = useRouter()
const questionnaireStore = useQuestionnaireStore()
const userStore = useUserStore()

// 检查用户是否已登录
if (!userStore.isLoggedIn) {
  Toast('请先登录')
  router.push('/login')
}

const loading = ref(false)
const showAddQuestionPopup = ref(false)
const addingQuestion = ref(false)
const isLoading = ref(true) // 初始加载状态

// 问卷信息
const questionnaireForm = reactive({
  id: 0,
  title: '',
  description: '',
  start_time: new Date(),
  end_time: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
  is_published: false,
  created_by: userStore.userInfo.id || 0
})

// 问题列表
const questions = ref([])

// 当前编辑的问题
const currentQuestion = reactive({
  title: '',
  type: '单选题',
  required: false,
  options: '',
  sort: 0
})

// 问题类型选项
const questionTypes = [
  { text: '单选题', value: '单选题' },
  { text: '多选题', value: '多选题' },
  { text: '填空题', value: '填空题' },
  { text: '评分题', value: '评分题' }
]

// 加载问卷数据
const loadQuestionnaireData = async () => {
  const id = route.params.id
  if (!id) {
    Toast('问卷ID不能为空')
    router.push('/questionnaire/list')
    return
  }
  
  isLoading.value = true
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
    
    const questionnaireData = response.data.questionnaire || {}
    const questionsData = response.data.questions || []
    
    // 如果问卷为空，显示错误
    if (!questionnaireData.id) {
      throw new Error('问卷不存在或已被删除')
    }
    
    // 检查是否已发布
    if (questionnaireData.is_published) {
      Dialog.alert({
        title: '无法编辑',
        message: '已发布的问卷不能再编辑。',
      }).then(() => {
        router.push(`/questionnaire/detail/${id}`)
      })
      return
    }
    
    // 检查编辑权限
    if (questionnaireData.created_by !== userStore.userInfo.id && !userStore.isAdmin) {
      Dialog.alert({
        title: '无权编辑',
        message: '您没有权限编辑此问卷。',
      }).then(() => {
        router.push(`/questionnaire/detail/${id}`)
      })
      return
    }
    
    // 填充表单数据
    Object.assign(questionnaireForm, {
      id: questionnaireData.id,
      title: questionnaireData.title || '',
      description: questionnaireData.description || '',
      start_time: questionnaireData.start_time ? new Date(questionnaireData.start_time) : new Date(),
      end_time: questionnaireData.end_time ? new Date(questionnaireData.end_time) : new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
      is_published: questionnaireData.is_published || false,
      created_by: questionnaireData.created_by || userStore.userInfo.id
    })
    
    // 填充问题数据
    questions.value = questionsData.map(q => ({
      id: q.id,
      title: q.title || '',
      type: q.type || '单选题',
      required: q.required || false,
      options: q.options || '',
      sort: q.sort || 0
    }))
    
  } catch (error) {
    console.error('加载问卷数据失败:', error)
    Toast.clear()
    Toast.fail('加载失败: ' + (error.message || '未知错误'))
    router.push('/questionnaire/list')
  } finally {
    isLoading.value = false
    Toast.clear()
  }
}

// 添加问题
const addQuestion = () => {
  try {
    // 防止重复提交
    if (addingQuestion.value) return;
    addingQuestion.value = true;
    
    // 验证问题
    if (!currentQuestion.title || currentQuestion.title.trim() === '') {
      Toast('请输入问题标题');
      addingQuestion.value = false;
      return;
    }
    
    // 验证选项
    if ((currentQuestion.type === '单选题' || currentQuestion.type === '多选题')) {
      if (!currentQuestion.options || currentQuestion.options.trim() === '') {
        Toast('请输入选项');
        addingQuestion.value = false;
        return;
      }
      
      // 检查选项格式
      const options = currentQuestion.options.split(',').filter(opt => opt.trim() !== '');
      if (options.length < 2) {
        Toast('请至少输入两个选项，并用逗号分隔');
        addingQuestion.value = false;
        return;
      }
    }
    
    // 添加到问题列表
    const newQuestion = {
      ...JSON.parse(JSON.stringify(currentQuestion)), // 深拷贝，避免引用问题
      id: Date.now(), // 临时ID
      sort: questions.value.length
    };
    
    questions.value.push(newQuestion);
    
    // 重置当前问题
    Object.assign(currentQuestion, {
      title: '',
      type: '单选题',
      required: false,
      options: '',
      sort: questions.value.length
    });
    
    // 关闭弹窗
    showAddQuestionPopup.value = false;
    
    Toast.success('问题添加成功');
  } catch (error) {
    console.error('添加问题失败:', error);
    Toast.fail('添加问题失败: ' + (error.message || '未知错误'));
  } finally {
    addingQuestion.value = false;
  }
}

// 删除问题
const deleteQuestion = (index) => {
  try {
    questions.value.splice(index, 1)
    // 更新排序
    questions.value.forEach((q, i) => {
      q.sort = i
    })
    
    Toast.success('问题删除成功')
  } catch (error) {
    console.error('删除问题失败:', error);
    Toast.fail('删除问题失败: ' + (error.message || '未知错误'));
  }
}

// 验证表单
const validateForm = () => {
  try {
    if (!questionnaireForm.title || questionnaireForm.title.trim() === '') {
      Toast('请输入问卷标题');
      return false;
    }
    
    if (!questionnaireForm.description || questionnaireForm.description.trim() === '') {
      Toast('请输入问卷描述');
      return false;
    }
    
    if (questions.value.length === 0) {
      Toast('请至少添加一个问题');
      return false;
    }
    
    // 检查每个问题的格式是否正确
    for (let i = 0; i < questions.value.length; i++) {
      const q = questions.value[i];
      
      if (!q.title || q.title.trim() === '') {
        Toast(`第${i+1}个问题标题不能为空`);
        return false;
      }
      
      if ((q.type === '单选题' || q.type === '多选题')) {
        if (!q.options || q.options.trim() === '') {
          Toast(`第${i+1}个问题的选项不能为空`);
          return false;
        }
        
        const options = q.options.split(',').filter(opt => opt.trim() !== '');
        if (options.length < 2) {
          Toast(`第${i+1}个问题需要至少两个选项`);
          return false;
        }
      }
    }
    
    return true;
  } catch (error) {
    console.error('表单验证错误:', error);
    Toast('表单验证失败: ' + (error.message || '未知错误'));
    return false;
  }
}

// 更新问卷
const updateQuestionnaire = async () => {
  try {
    if (!validateForm()) {
      return;
    }
    
    loading.value = true;
    Toast.loading({
      message: '保存中...',
      forbidClick: true,
      duration: 0
    });
    
    // 准备问卷数据
    const questionnaireData = {
      id: questionnaireForm.id,
      title: questionnaireForm.title.trim(),
      description: questionnaireForm.description.trim(),
      start_time: questionnaireForm.start_time,
      end_time: questionnaireForm.end_time,
      is_published: questionnaireForm.is_published,
      created_by: questionnaireForm.created_by,
      questions: questions.value.map((q, index) => ({
        id: q.id,
        title: q.title.trim(),
        type: q.type,
        required: q.required,
        options: q.options,
        sort: index
      }))
    };
    
    console.log('准备发送问卷数据:', JSON.stringify(questionnaireData));
    
    await questionnaireStore.updateQuestionnaire(questionnaireData);
    
    Toast.clear();
    Toast.success('保存成功');
    
    // 保存成功后显示对话框
    Dialog.confirm({
      title: '保存成功',
      message: '问卷已成功保存！',
      confirmButtonText: '查看详情',
      confirmButtonColor: '#1989fa',
    }).then(() => {
      router.push(`/questionnaire/detail/${questionnaireForm.id}`)
    }).catch(() => {
      // 用户取消对话框，留在当前页面
    });
  } catch (error) {
    console.error('更新问卷失败:', error);
    Toast.clear();
    
    // 提供更具体的错误信息
    let errorMessage = '未知错误';
    if (error.message === '请求超时，后端服务可能不可用') {
      errorMessage = '后端服务连接失败，请确认服务是否启动';
    } else if (error.code === 'ECONNABORTED') {
      errorMessage = '请求超时，请检查网络或后端服务';
    } else if (error.code === 'ECONNREFUSED' || (error.message && error.message.includes('ECONNREFUSED'))) {
      errorMessage = '无法连接到后端服务，请检查服务是否运行';
    } else if (error.response) {
      errorMessage = error.response.data?.message || `服务器错误 (${error.response.status})`;
    } else if (error.request) {
      errorMessage = '服务器无响应';
    } else {
      errorMessage = error.message || '保存失败';
    }
    
    Toast.fail(errorMessage);
  } finally {
    loading.value = false;
  }
}

// 打开添加问题弹窗
const openAddQuestionPopup = () => {
  // 重置当前问题表单
  Object.assign(currentQuestion, {
    title: '',
    type: '单选题',
    required: false,
    options: '',
    sort: questions.value.length
  });
  
  // 打开弹窗
  showAddQuestionPopup.value = true;
}

// 取消编辑
const cancelEdit = () => {
  Dialog.confirm({
    title: '取消编辑',
    message: '确定要取消编辑吗？未保存的修改将丢失。',
  }).then(() => {
    router.push(`/questionnaire/detail/${route.params.id}`)
  }).catch(() => {
    // 用户取消操作
  })
}

onMounted(() => {
  loadQuestionnaireData()
})
</script>

<template>
  <div class="edit-container">
    <van-nav-bar
      title="编辑问卷"
      left-text="返回"
      left-arrow
      @click-left="cancelEdit"
      class="custom-nav"
    />
    
    <div v-if="!isLoading" class="content">
      <div class="questionnaire-form">
        <!-- 问卷基本信息 -->
        <van-cell-group inset title="基本信息">
          <van-field
            v-model="questionnaireForm.title"
            name="title"
            label="问卷标题"
            placeholder="请输入问卷标题"
            :rules="[{ required: true, message: '请输入问卷标题' }]"
          />
          <van-field
            v-model="questionnaireForm.description"
            name="description"
            label="问卷描述"
            type="textarea"
            rows="3"
            placeholder="请输入问卷描述"
            :rules="[{ required: true, message: '请输入问卷描述' }]"
          />
        </van-cell-group>
        
        <!-- 问题列表 -->
        <van-cell-group inset title="问题列表" class="question-group">
          <template v-if="questions.length > 0">
            <div
              v-for="(question, index) in questions"
              :key="question.id"
              class="question-item"
            >
              <div class="question-header">
                <span class="question-index">{{ index + 1 }}.</span>
                <span class="question-title">{{ question.title }}</span>
                <span class="question-type">[{{ question.type }}]</span>
                <span class="question-required" v-if="question.required">*必填</span>
              </div>
              
              <div class="question-options" v-if="question.options">
                <p>选项：{{ question.options }}</p>
              </div>
              
              <div class="question-actions">
                <van-button size="small" type="danger" @click="deleteQuestion(index)">删除</van-button>
              </div>
            </div>
          </template>
          <template v-else>
            <van-empty description="暂无问题" />
          </template>
          
          <div class="add-question">
            <van-button type="primary" block @click="openAddQuestionPopup">添加问题</van-button>
          </div>
        </van-cell-group>
        
        <!-- 提交按钮 -->
        <div class="submit-area">
          <van-button 
            type="primary" 
            block 
            round 
            @click="updateQuestionnaire"
            :loading="loading"
            :disabled="loading"
          >
            {{ loading ? '保存中...' : '保存问卷' }}
          </van-button>
        </div>
      </div>
    </div>
    
    <div class="loading" v-else>
      <van-loading type="spinner" />
    </div>
    
    <!-- 添加问题弹窗 -->
    <van-popup
      v-model:show="showAddQuestionPopup"
      position="bottom"
      round
      closeable
      :style="{ height: '70%' }"
    >
      <div class="popup-title">添加问题</div>
      <div class="popup-content">
        <van-cell-group inset>
          <van-field
            v-model="currentQuestion.title"
            name="title"
            label="问题标题"
            placeholder="请输入问题标题"
            :rules="[{ required: true, message: '请输入问题标题' }]"
          />
          <van-field name="type" label="问题类型">
            <template #input>
              <van-radio-group v-model="currentQuestion.type">
                <van-radio v-for="item in questionTypes" :key="item.value" :name="item.value">
                  {{ item.text }}
                </van-radio>
              </van-radio-group>
            </template>
          </van-field>
          <van-field
            v-if="currentQuestion.type === '单选题' || currentQuestion.type === '多选题'"
            v-model="currentQuestion.options"
            name="options"
            label="选项"
            type="textarea"
            rows="3"
            placeholder="请输入选项，用逗号分隔"
            :rules="[{ required: true, message: '请输入选项' }]"
          />
          <van-field name="required" label="是否必填">
            <template #input>
              <van-switch v-model="currentQuestion.required" />
            </template>
          </van-field>
        </van-cell-group>
        
        <div style="margin: 16px;">
          <van-button 
            round 
            block 
            type="primary" 
            @click="addQuestion"
            :loading="addingQuestion"
            :disabled="addingQuestion"
          >
            {{ addingQuestion ? '添加中...' : '添加问题' }}
          </van-button>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<style scoped>
.edit-container {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.custom-nav {
  background: linear-gradient(135deg, #4481eb 0%, #04befe 100%);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

:deep(.custom-nav .van-nav-bar__title) {
  color: white;
  font-weight: bold;
}

:deep(.custom-nav .van-icon),
:deep(.custom-nav .van-nav-bar__text) {
  color: white !important;
}

.content {
  padding: 16px;
}

.questionnaire-form {
  width: 100%;
}

.question-group {
  margin-top: 16px;
}

.question-item {
  background-color: #fff;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 12px;
  border-left: 3px solid #1989fa;
}

.question-header {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.question-index {
  margin-right: 4px;
  font-weight: bold;
}

.question-title {
  flex: 1;
  font-weight: bold;
}

.question-type {
  margin: 0 8px;
  color: #1989fa;
  font-size: 12px;
}

.question-required {
  color: #ee0a24;
  font-size: 12px;
}

.question-options {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.question-actions {
  display: flex;
  justify-content: flex-end;
}

.add-question {
  margin-top: 16px;
}

.submit-area {
  margin-top: 24px;
  padding: 0 16px 32px;
}

.popup-title {
  text-align: center;
  padding: 16px;
  font-size: 16px;
  font-weight: bold;
}

.popup-content {
  padding: 16px;
  max-height: calc(70vh - 60px);
  overflow-y: auto;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

/* 添加问题类型单选按钮样式 */
:deep(.van-radio-group) {
  display: flex;
  flex-direction: column;
  padding: 8px 0;
}

:deep(.van-radio) {
  margin-bottom: 8px;
}
</style> 