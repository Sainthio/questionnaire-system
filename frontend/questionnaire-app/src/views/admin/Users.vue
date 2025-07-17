<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'
import { userApi } from '../../api/admin'
import { Toast, Dialog } from 'vant'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(true)
const users = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 用户详情相关
const showUserDetail = ref(false)
const currentUser = ref(null)
const userDetailLoading = ref(false)

// 用户编辑相关
const showEditDialog = ref(false)
const editForm = ref({
  id: null,
  email: '',
  phone: '',
  is_admin: false
})

// 修改密码相关
const showPasswordDialog = ref(false)
const passwordForm = ref({
  id: null,
  username: '',
  new_password: '',
  confirm_password: ''
})

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
  
  loadUsers()
})

// 加载用户列表
const loadUsers = async () => {
  loading.value = true
  try {
    const response = await userApi.getAllUsers(currentPage.value, pageSize.value)
    users.value = response.data.data.users
    total.value = response.data.data.total
    console.log('用户列表:', users.value)
  } catch (error) {
    console.error('获取用户列表失败:', error)
    Toast('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 查看用户详情
const viewUserDetail = async (userId) => {
  userDetailLoading.value = true
  showUserDetail.value = true
  
  try {
    const response = await userApi.getUserDetail(userId)
    currentUser.value = response.data.data
    
    // 确保用户数据存在
    if (!currentUser.value) {
      currentUser.value = {
        user: {
          id: userId,
          username: '未知用户',
          email: '',
          phone: '',
          is_admin: false,
          created_at: new Date(),
          updated_at: new Date()
        },
        questionnaire_count: 0,
        submission_count: 0
      }
    }
    
    // 确保user对象存在
    if (!currentUser.value.user) {
      currentUser.value.user = {
        id: userId,
        username: '未知用户',
        email: '',
        phone: '',
        is_admin: false,
        created_at: new Date(),
        updated_at: new Date()
      }
    }
    
    // 确保统计数据存在
    if (currentUser.value.questionnaire_count === undefined) {
      currentUser.value.questionnaire_count = 0
    }
    if (currentUser.value.submission_count === undefined) {
      currentUser.value.submission_count = 0
    }
    
    console.log('用户详情:', currentUser.value)
  } catch (error) {
    console.error('获取用户详情失败:', error)
    Toast('获取用户详情失败')
    showUserDetail.value = false
  } finally {
    userDetailLoading.value = false
  }
}

// 编辑用户
const editUser = (user) => {
  editForm.value = {
    id: user.id,
    email: user.email || '',
    phone: user.phone || '',
    is_admin: user.is_admin
  }
  showEditDialog.value = true
}

// 保存用户编辑
const saveUserEdit = async () => {
  try {
    await userApi.updateUser(editForm.value)
    Toast('用户更新成功')
    showEditDialog.value = false
    loadUsers() // 重新加载用户列表
  } catch (error) {
    console.error('更新用户失败:', error)
    Toast('更新用户失败')
  }
}

// 修改用户密码
const changePassword = (user) => {
  passwordForm.value = {
    id: user.id,
    username: user.username,
    new_password: '',
    confirm_password: ''
  }
  showPasswordDialog.value = true
}

// 保存新密码
const saveNewPassword = async () => {
  // 验证密码
  if (!passwordForm.value.new_password) {
    Toast('请输入新密码')
    return
  }
  
  if (passwordForm.value.new_password.length < 6) {
    Toast('密码长度不能少于6位')
    return
  }
  
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    Toast('两次输入的密码不一致')
    return
  }
  
  try {
    // 调用后端API重置密码
    await fetch('/api/user/reset-password', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        username: passwordForm.value.username,
        new_password: passwordForm.value.new_password
      })
    })
    
    Toast('密码修改成功')
    showPasswordDialog.value = false
  } catch (error) {
    console.error('修改密码失败:', error)
    Toast('修改密码失败')
  }
}

// 删除用户
const confirmDeleteUser = (userId, username) => {
  Dialog.confirm({
    title: '确认删除',
    message: `确定要删除用户 "${username}" 吗？此操作不可撤销，将同时删除该用户创建的问卷和提交的答案。`,
    showCancelButton: true
  }).then(async () => {
    try {
      await userApi.deleteUser(userId)
      Toast('用户删除成功')
      loadUsers() // 重新加载用户列表
    } catch (error) {
      console.error('删除用户失败:', error)
      Toast('删除用户失败')
    }
  })
}

// 页面变化
const onPageChange = (page) => {
  currentPage.value = page
  loadUsers()
}
</script>

<template>
  <div class="admin-users">
    <van-nav-bar
      title="用户管理"
      left-text="返回"
      left-arrow
      @click-left="router.push('/admin')"
    />
    
    <div class="users-container">
      <van-loading v-if="loading" size="24px" vertical>加载中...</van-loading>
      
      <template v-else>
        <div class="table-header">
          <h2>用户列表 ({{ total }})</h2>
        </div>
        
        <div class="users-table">
          <van-cell-group inset>
            <!-- 表头 -->
            <van-cell title-class="table-title" value-class="table-title">
              <template #title>
                <div class="cell-title">用户名</div>
              </template>
              <template #value>
                <div class="cell-value">
                  <span class="email-column">邮箱</span>
                  <span class="role-column">角色</span>
                  <span class="action-column">操作</span>
                </div>
              </template>
            </van-cell>
            
            <!-- 表格内容 -->
            <van-cell v-for="user in users" :key="user.id" clickable @click="viewUserDetail(user.id)">
              <template #title>
                <div class="cell-title">
                  {{ user.username }}
                </div>
              </template>
              <template #value>
                <div class="cell-value">
                  <span class="email-column">{{ user.email || '-' }}</span>
                  <span class="role-column">
                    <van-tag type="primary" v-if="user.is_admin">管理员</van-tag>
                    <van-tag type="success" v-else>普通用户</van-tag>
                  </span>
                  <span class="action-column">
                    <van-button size="mini" type="primary" @click.stop="editUser(user)">编辑</van-button>
                    <van-button size="mini" type="warning" @click.stop="changePassword(user)">修改密码</van-button>
                    <van-button 
                      size="mini" 
                      type="danger" 
                      @click.stop="confirmDeleteUser(user.id, user.username)"
                      :disabled="user.is_admin"
                    >
                      删除
                    </van-button>
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
    
    <!-- 用户详情弹窗 -->
    <van-popup v-model:show="showUserDetail" round position="bottom" :style="{ height: '70%' }">
      <div class="user-detail">
        <van-nav-bar
          title="用户详情"
          left-text="返回"
          left-arrow
          @click-left="showUserDetail = false"
        />
        
        <van-loading v-if="userDetailLoading" size="24px" vertical>加载中...</van-loading>
        
        <template v-else-if="currentUser && currentUser.user">
          <div class="detail-content">
            <div class="user-info">
              <div class="avatar">
                <van-icon name="manager" size="48" />
              </div>
              <div class="basic-info">
                <h3>{{ currentUser.user.username }}</h3>
                <p>
                  <van-tag type="primary" v-if="currentUser.user.is_admin">管理员</van-tag>
                  <van-tag type="success" v-else>普通用户</van-tag>
                </p>
              </div>
            </div>
            
            <van-cell-group inset title="基本信息">
              <van-cell title="用户ID" :value="currentUser.user.id" />
              <van-cell title="邮箱" :value="currentUser.user.email || '-'" />
              <van-cell title="手机" :value="currentUser.user.phone || '-'" />
              <van-cell title="注册时间" :value="new Date(currentUser.user.created_at).toLocaleString()" />
              <van-cell title="最后更新" :value="new Date(currentUser.user.updated_at).toLocaleString()" />
            </van-cell-group>
            
            <van-cell-group inset title="活动统计">
              <van-cell title="创建的问卷" :value="currentUser.questionnaire_count || 0" />
              <van-cell title="提交的答卷" :value="currentUser.submission_count || 0" />
            </van-cell-group>
            
            <div class="detail-actions">
              <van-button type="primary" block @click="editUser(currentUser.user)">编辑用户</van-button>
              <van-button type="warning" block @click="changePassword(currentUser.user)">修改密码</van-button>
              <van-button 
                type="danger" 
                block 
                @click="confirmDeleteUser(currentUser.user.id, currentUser.user.username)"
                :disabled="currentUser.user.is_admin"
              >
                删除用户
              </van-button>
            </div>
          </div>
        </template>
        
        <div v-else class="error-message">
          <van-empty description="无法加载用户数据">
            <template #bottom>
              <van-button round type="primary" @click="showUserDetail = false">返回</van-button>
            </template>
          </van-empty>
        </div>
      </div>
    </van-popup>
    
    <!-- 编辑用户弹窗 -->
    <van-popup v-model:show="showEditDialog" round position="bottom" :style="{ height: '50%' }">
      <div class="edit-form">
        <van-nav-bar
          title="编辑用户"
          left-text="取消"
          right-text="保存"
          left-arrow
          @click-left="showEditDialog = false"
          @click-right="saveUserEdit"
        />
        
        <van-form>
          <van-cell-group inset>
            <van-field
              v-model="editForm.email"
              label="邮箱"
              placeholder="请输入邮箱"
              type="email"
            />
            
            <van-field
              v-model="editForm.phone"
              label="手机号"
              placeholder="请输入手机号"
              type="tel"
            />
            
            <van-cell title="管理员权限">
              <template #right-icon>
                <van-switch v-model="editForm.is_admin" size="24" />
              </template>
            </van-cell>
          </van-cell-group>
        </van-form>
      </div>
    </van-popup>
    
    <!-- 修改密码弹窗 -->
    <van-popup v-model:show="showPasswordDialog" round position="bottom" :style="{ height: '50%' }">
      <div class="edit-form">
        <van-nav-bar
          title="修改密码"
          left-text="取消"
          right-text="保存"
          left-arrow
          @click-left="showPasswordDialog = false"
          @click-right="saveNewPassword"
        />
        
        <van-form>
          <van-cell-group inset>
            <van-field
              readonly
              label="用户名"
              :value="passwordForm.username"
            />
            
            <van-field
              v-model="passwordForm.new_password"
              label="新密码"
              placeholder="请输入新密码"
              type="password"
            />
            
            <van-field
              v-model="passwordForm.confirm_password"
              label="确认密码"
              placeholder="请再次输入新密码"
              type="password"
            />
          </van-cell-group>
          
          <div class="password-tips">
            <p>密码要求：</p>
            <ul>
              <li>长度不少于6位</li>
              <li>建议包含字母、数字和特殊字符</li>
            </ul>
          </div>
        </van-form>
      </div>
    </van-popup>
  </div>
</template>

<style scoped>
.admin-users {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.users-container {
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

.users-table {
  margin-bottom: 20px;
}

.table-title {
  font-weight: bold;
  color: #323233;
}

.cell-title {
  font-weight: bold;
}

.cell-value {
  display: flex;
  align-items: center;
}

.email-column {
  flex: 1;
  margin-right: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.role-column {
  width: 70px;
  text-align: center;
  margin-right: 10px;
}

.action-column {
  width: 180px;
  display: flex;
  justify-content: space-between;
}

.action-column .van-button {
  margin-left: 4px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.user-detail {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.detail-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px;
  background-color: #fff;
  border-radius: 8px;
}

.avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background-color: #f2f3f5;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-right: 16px;
}

.basic-info h3 {
  margin: 0 0 8px;
  font-size: 18px;
}

.basic-info p {
  margin: 0;
}

.detail-actions {
  margin-top: 20px;
  display: grid;
  grid-gap: 12px;
}

.edit-form {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.password-tips {
  padding: 16px;
  color: #646566;
  font-size: 14px;
}

.password-tips p {
  margin: 0 0 8px;
}

.password-tips ul {
  margin: 0;
  padding-left: 20px;
}

.error-message {
  padding: 40px 16px;
}
</style> 