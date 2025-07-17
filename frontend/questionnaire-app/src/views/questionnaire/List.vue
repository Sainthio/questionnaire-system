<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useQuestionnaireStore } from '../../stores/questionnaire'
import { useUserStore } from '../../stores/user'
import { Toast, Dialog } from 'vant'
import axios from 'axios'
import api from '../../services/api'

const router = useRouter()
const questionnaireStore = useQuestionnaireStore()
const userStore = useUserStore()

const loading = ref(false)
const refreshing = ref(false)
const finished = ref(false)
const showAnimation = ref(false)
const error = ref(false)

// 计算属性：确保问卷列表始终是数组
const questionnaireList = computed(() => {
  return questionnaireStore.questionnaireList || []
})

onMounted(() => {
  loadQuestionnaireList()
  
  // 添加进入动画
  setTimeout(() => {
    showAnimation.value = true
  }, 100)
})

// 获取问卷列表
const loadQuestionnaireList = async (page = 1) => {
  if (loading.value) return
  
  loading.value = true
  error.value = false
  
  // 设置超时，增加超时时间
  const timeoutPromise = new Promise((_, reject) => {
    setTimeout(() => reject(new Error('请求超时，后端服务可能不可用')), 30000);
  });
  
  Toast.loading({
    message: '加载中...',
    forbidClick: true,
    duration: 0
  })
  
  try {
    // 使用store方法获取数据
    await Promise.race([
      questionnaireStore.getQuestionnaireList(page),
      timeoutPromise
    ]);
    
    if ((questionnaireStore.questionnaireList || []).length >= questionnaireStore.totalCount) {
      finished.value = true;
    }
  } catch (err) {
    console.error('加载问卷列表失败:', err);
    error.value = true;
    
    // 提供更具体的错误信息
    let errorMessage = '未知错误';
    if (err.message === '请求超时，后端服务可能不可用') {
      errorMessage = '后端服务响应超时，请检查服务是否正常运行';
    } else if (err.code === 'ECONNABORTED') {
      errorMessage = '请求超时，请检查网络或后端服务';
    } else if (err.code === 'ECONNREFUSED' || (err.message && err.message.includes('ECONNREFUSED'))) {
      errorMessage = '无法连接到后端服务，请检查服务是否运行';
    } else if (err.response) {
      errorMessage = err.response.data?.message || `服务器错误 (${err.response.status})`;
    } else if (err.request) {
      errorMessage = '服务器无响应';
    } else {
      errorMessage = err.message || '加载失败';
    }
    
    Toast.clear();
    Toast.fail(errorMessage);
    
    // 设置一些默认数据，以防页面显示异常
    if (!questionnaireStore.questionnaireList || questionnaireStore.questionnaireList.length === 0) {
      questionnaireStore.questionnaireList = [];
      questionnaireStore.totalCount = 0;
    }
  } finally {
    loading.value = false;
    refreshing.value = false;
    Toast.clear();
  }
}

// 刷新列表
const onRefresh = () => {
  finished.value = false
  questionnaireStore.currentPage = 1
  loadQuestionnaireList(1)
}

// 加载更多
const onLoad = () => {
  const nextPage = questionnaireStore.currentPage + 1
  loadQuestionnaireList(nextPage)
}

// 查看问卷详情
const viewDetail = (id) => {
  router.push(`/questionnaire/detail/${id}`)
}

// 填写问卷
const fillQuestionnaire = (item) => {
  const id = item.questionnaire ? item.questionnaire.id : item.id
  const isPublished = item.questionnaire ? item.questionnaire.is_published : item.is_published
  
  // 检查问卷是否已发布
  if (!isPublished) {
    Dialog.alert({
      title: '无法填写',
      message: '该问卷尚未发布，暂时无法填写。',
    })
    return
  }
  
  router.push(`/questionnaire/fill/${id}`)
}

// 编辑问卷
const editQuestionnaire = (item) => {
  const id = item.questionnaire ? item.questionnaire.id : item.id
  const createdBy = item.questionnaire ? item.questionnaire.created_by : item.created_by
  const isPublished = item.questionnaire ? item.questionnaire.is_published : item.is_published
  
  // 检查是否是创建者
  if (userStore.userInfo.id !== createdBy) {
    Toast('只有问卷创建者可以编辑问卷')
    return
  }
  
  // 检查问卷是否已发布
  if (isPublished) {
    Toast('已发布的问卷无法编辑')
    return
  }
  
  router.push(`/questionnaire/edit/${id}`)
}

// 删除问卷
const deleteQuestionnaire = (item) => {
  const id = item.questionnaire ? item.questionnaire.id : item.id
  const createdBy = item.questionnaire ? item.questionnaire.created_by : item.created_by
  const title = item.questionnaire ? item.questionnaire.title : item.title
  
  // 检查是否是创建者或管理员
  if (userStore.userInfo.id !== createdBy && !userStore.isAdmin) {
    Toast('只有问卷创建者或管理员可以删除问卷')
    return
  }
  
  Dialog.confirm({
    title: '确认删除',
    message: `确定要删除问卷"${title}"吗？此操作不可恢复。`,
  })
  .then(async () => {
    try {
      Toast.loading({
        message: '删除中...',
        forbidClick: true,
        duration: 0
      })
      
      // 调用store中的删除方法
      await questionnaireStore.deleteQuestionnaire(id)
      
      Toast.success('删除成功')
      // 刷新列表
      onRefresh()
    } catch (err) {
      console.error('删除问卷失败:', err)
      Toast.fail(err.message || '删除失败，请稍后再试')
    } finally {
      Toast.clear()
    }
  })
  .catch(() => {
    // 取消删除
  })
}

// 发布或取消发布问卷
const togglePublishStatus = async (item) => {
  const id = item.questionnaire ? item.questionnaire.id : item.id
  const createdBy = item.questionnaire ? item.questionnaire.created_by : item.created_by
  const title = item.questionnaire ? item.questionnaire.title : item.title
  const isPublished = item.questionnaire ? item.questionnaire.is_published : item.is_published
  
  // 检查是否是创建者或管理员
  if (userStore.userInfo.id !== createdBy && !userStore.isAdmin) {
    Toast('只有问卷创建者或管理员可以修改发布状态')
    return
  }
  
  const newStatus = !isPublished
  const actionText = newStatus ? '发布' : '取消发布'
  
  Dialog.confirm({
    title: `确认${actionText}`,
    message: `确定要${actionText}问卷"${title}"吗？`,
  })
  .then(async () => {
    try {
      Toast.loading({
        message: `${actionText}中...`,
        forbidClick: true,
        duration: 0
      })
      
      // 调用API更新问卷状态
      await questionnaireStore.updateQuestionnaireStatus(id, newStatus)
      
      Toast.success(`${actionText}成功`)
      // 刷新列表
      onRefresh()
    } catch (err) {
      console.error(`${actionText}问卷失败:`, err)
      Toast.fail(`${actionText}失败，请稍后再试`)
    } finally {
      Toast.clear()
    }
  })
  .catch(() => {
    // 取消操作
  })
}

// 创建问卷
const createQuestionnaire = async () => {
  if (!userStore.isLoggedIn) {
    Toast('请先登录')
    router.push('/login')
    return
  }
  
  try {
    // 检查后端服务是否可用
    Toast.loading({
      message: '检查服务可用性...',
      forbidClick: true,
      duration: 0
    });
    
    const checkPromise = fetch('/api/health').then(res => res.ok);
    const timeoutPromise = new Promise(resolve => setTimeout(() => resolve(false), 3000));
    
    const isServerAvailable = await Promise.race([checkPromise, timeoutPromise]);
    
    Toast.clear();
    
    if (!isServerAvailable) {
      Toast.fail('后端服务不可用，请稍后再试');
      return;
    }
    
    // 后端可用，跳转到创建页面
    router.push('/questionnaire/create')
  } catch (err) {
    console.error('检查后端服务失败:', err);
    Toast.clear();
    Toast.fail('无法连接到后端服务，请确认服务是否启动');
  }
}

// 获取状态标签类型
const getStatusTagType = (isPublished) => {
  return isPublished ? 'primary' : 'warning'
}

// 获取状态标签文本
const getStatusText = (isPublished) => {
  return isPublished ? '已发布' : '未发布'
}

// 处理返回按钮点击
const handleBackClick = () => {
  if (userStore.isAdmin) {
    // 管理员用户返回到管理员页面
    router.push('/admin')
  } else {
    // 普通用户返回到首页
    router.push('/')
  }
}
</script>

<template>
  <div class="list-container">
    <van-nav-bar
      title="问卷列表"
      left-text="返回"
      left-arrow
      right-text="创建"
      @click-left="handleBackClick"
      @click-right="createQuestionnaire"
      class="custom-nav"
    />
    
    <div class="list-wrapper" :class="{ 'show-animation': showAnimation }">
      <div class="list-header">
        <h2 class="list-title">所有问卷</h2>
        <p class="list-subtitle">浏览、填写或创建您的问卷</p>
      </div>
      
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh" class="pull-refresh">
        <van-list
          v-model:loading="loading"
          :finished="finished"
          finished-text="没有更多了"
          @load="onLoad"
          class="questionnaire-list"
        >
          <div class="list-content">
            <template v-if="questionnaireList.length > 0">
              <div 
                v-for="item in questionnaireList" 
                :key="item.questionnaire ? item.questionnaire.id : item.id"
                class="questionnaire-card"
              >
                <div class="card-header">
                  <div class="card-title-area">
                    <h3 class="card-title">{{ item.questionnaire ? item.questionnaire.title : item.title }}</h3>
                    <van-tag 
                      :type="getStatusTagType(item.questionnaire ? item.questionnaire.is_published : item.is_published)" 
                      round
                      class="status-tag"
                    >
                      {{ getStatusText(item.questionnaire ? item.questionnaire.is_published : item.is_published) }}
                    </van-tag>
                  </div>
                  <p class="card-desc">{{ (item.questionnaire ? item.questionnaire.description : item.description) || '暂无描述' }}</p>
                </div>
                
                <div class="card-content">
                  <div class="card-info">
                    <div class="info-item">
                      <van-icon name="clock-o" class="info-icon" />
                      <span>{{ new Date(item.questionnaire ? item.questionnaire.created_at : item.created_at).toLocaleDateString() }}</span>
                    </div>
                    <div class="info-item">
                      <van-icon name="user-o" class="info-icon" />
                      <span>创建者: {{ item.creator_name || '未知' }}</span>
                    </div>
                  </div>
                </div>
                
                <div class="card-footer">
                  <!-- 详情按钮，所有人都可以看到 -->
                  <van-button 
                    size="small" 
                    type="primary" 
                    icon="description" 
                    @click="viewDetail(item.questionnaire ? item.questionnaire.id : item.id)"
                    class="action-btn"
                  >
                    详情
                  </van-button>
                  
                  <!-- 已发布的问卷显示填写按钮 -->
                  <van-button 
                    v-if="item.questionnaire ? item.questionnaire.is_published : item.is_published"
                    size="small" 
                    type="success" 
                    icon="edit" 
                    @click="fillQuestionnaire(item)"
                    class="action-btn"
                  >
                    填写
                  </van-button>
                  
                  <!-- 如果是问卷创建者，显示更多操作按钮 -->
                  <template v-if="userStore.userInfo.id === (item.questionnaire ? item.questionnaire.created_by : item.created_by) || userStore.isAdmin">
                    <!-- 未发布的问卷显示编辑按钮 -->
                    <van-button 
                      v-if="!(item.questionnaire ? item.questionnaire.is_published : item.is_published)"
                      size="small" 
                      type="warning" 
                      icon="edit" 
                      @click="editQuestionnaire(item)"
                      class="action-btn"
                    >
                      编辑
                    </van-button>
                    
                    <!-- 发布/取消发布按钮 -->
                    <van-button 
                      size="small" 
                      :type="(item.questionnaire ? item.questionnaire.is_published : item.is_published) ? 'default' : 'primary'" 
                      :icon="(item.questionnaire ? item.questionnaire.is_published : item.is_published) ? 'down' : 'ascending'" 
                      @click="togglePublishStatus(item)"
                      class="action-btn"
                    >
                      {{ (item.questionnaire ? item.questionnaire.is_published : item.is_published) ? '取消发布' : '发布' }}
                    </van-button>
                    
                    <!-- 删除按钮 -->
                    <van-button 
                      size="small" 
                      type="danger" 
                      icon="delete" 
                      @click="deleteQuestionnaire(item)"
                      class="action-btn"
                    >
                      删除
                    </van-button>
                  </template>
                </div>
              </div>
            </template>
            <template v-else-if="error">
              <div class="error-container">
                <van-empty description="加载失败，请重试">
                  <template #image>
                    <div class="error-icon-container">
                      <van-icon name="warning-o" size="48" class="error-icon" />
                    </div>
                  </template>
                  <div class="retry-button">
                    <van-button round type="danger" @click="onRefresh">
                      重新加载
                    </van-button>
                  </div>
                </van-empty>
              </div>
            </template>
            <template v-else>
              <div class="empty-list">
                <van-empty description="暂无问卷">
                  <template #image>
                    <van-icon name="description" size="64" class="empty-icon" />
                  </template>
                  <van-button round type="primary" size="small" @click="createQuestionnaire">
                    创建问卷
                  </van-button>
                </van-empty>
              </div>
            </template>
          </div>
        </van-list>
      </van-pull-refresh>
    </div>
    
    <!-- 浮动按钮 -->
    <div 
      v-if="userStore.isLoggedIn"
      class="floating-button"
      @click="createQuestionnaire"
    >
      <van-icon name="plus" size="24" />
    </div>
  </div>
</template>

<style scoped>
.list-container {
  min-height: 100vh;
  background-color: #f7f8fa;
  position: relative;
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

.list-wrapper {
  padding: 0;
  opacity: 0;
  transform: translateY(20px);
  transition: all 0.5s ease;
}

.show-animation {
  opacity: 1;
  transform: translateY(0);
}

.list-header {
  padding: 20px 16px;
  background: white;
  margin-bottom: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.list-title {
  margin: 0;
  font-size: 20px;
  font-weight: bold;
  color: #323233;
}

.list-subtitle {
  margin: 5px 0 0;
  font-size: 14px;
  color: #969799;
}

.pull-refresh {
  min-height: calc(100vh - 46px - 80px);
}

.list-content {
  padding: 12px;
}

.questionnaire-card {
  margin-bottom: 16px;
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.questionnaire-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.card-header {
  padding: 16px;
  border-bottom: 1px solid #f5f5f5;
}

.card-title-area {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.card-title {
  margin: 0;
  font-size: 16px;
  font-weight: bold;
  color: #323233;
  flex: 1;
}

.status-tag {
  flex-shrink: 0;
  margin-left: 8px;
}

.card-desc {
  margin: 0;
  font-size: 14px;
  color: #646566;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.card-content {
  padding: 12px 16px;
}

.card-info {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #969799;
}

.info-item {
  display: flex;
  align-items: center;
}

.info-icon {
  margin-right: 4px;
}

.card-footer {
  padding: 12px 16px;
  border-top: 1px solid #f5f5f5;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.action-btn {
  border-radius: 16px;
  font-size: 12px;
  transition: transform 0.2s ease;
}

.action-btn:active {
  transform: scale(0.95);
}

.unpublished-tag {
  margin-left: 8px;
  vertical-align: middle;
}

.empty-list {
  padding: 60px 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.empty-icon {
  color: #dcdee0;
}

.error-container {
  padding: 20px;
  text-align: center;
  margin: 40px 0;
}

.error-icon-container {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 8px;
}

.error-icon {
  color: #ee0a24;
}

.retry-button {
  margin-top: 16px;
}

/* 浮动按钮样式 */
.floating-button {
  position: fixed;
  right: 16px;
  bottom: 16px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #4481eb 0%, #04befe 100%);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  cursor: pointer;
  z-index: 99;
  transition: transform 0.2s ease;
}

.floating-button:active {
  transform: scale(0.95);
}

/* 桌面端样式 */
:deep(.desktop) .list-wrapper {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

:deep(.desktop) .list-header {
  border-radius: 12px;
  margin-bottom: 20px;
}

:deep(.desktop) .questionnaire-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

:deep(.desktop) .questionnaire-card {
  margin-bottom: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
}

:deep(.desktop) .card-content {
  flex-grow: 1;
}

:deep(.desktop) .floating-button {
  right: calc(50% - 400px + 16px);
}
</style> 