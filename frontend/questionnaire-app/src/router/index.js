import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import QuestionnaireList from '../views/questionnaire/List.vue'
import QuestionnaireDetail from '../views/questionnaire/Detail.vue'
import QuestionnaireFill from '../views/questionnaire/Fill.vue'
import QuestionnaireCreate from '../views/questionnaire/Create.vue'
import QuestionnaireEdit from '../views/questionnaire/Edit.vue'
import QuestionnaireResults from '../views/questionnaire/Results.vue'
import AdminDashboard from '../views/admin/Dashboard.vue'
import AdminUsers from '../views/admin/Users.vue'
import AdminQuestionnaires from '../views/admin/Questionnaires.vue'
import AdminStatistics from '../views/admin/Statistics.vue'

// 页面组件会在后续创建
const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  },
  {
    path: '/questionnaire/list',
    name: 'QuestionnaireList',
    component: QuestionnaireList
  },
  {
    path: '/questionnaire/detail/:id',
    name: 'QuestionnaireDetail',
    component: QuestionnaireDetail
  },
  {
    path: '/questionnaire/fill/:id',
    name: 'QuestionnaireFill',
    component: QuestionnaireFill
  },
  {
    path: '/questionnaire/create',
    name: 'QuestionnaireCreate',
    component: QuestionnaireCreate
  },
  {
    path: '/questionnaire/edit/:id',
    name: 'QuestionnaireEdit',
    component: QuestionnaireEdit
  },
  {
    path: '/questionnaire/results/:id',
    name: 'QuestionnaireResults',
    component: QuestionnaireResults
  },
  {
    path: '/questionnaire/results/:id/question/:questionId',
    name: 'QuestionnaireQuestionResults',
    component: QuestionnaireResults
  },
  // 管理员路由
  {
    path: '/admin',
    name: 'AdminDashboard',
    component: AdminDashboard,
    meta: { requiresAdmin: true }
  },
  {
    path: '/admin/users',
    name: 'AdminUsers',
    component: AdminUsers,
    meta: { requiresAdmin: true }
  },
  {
    path: '/admin/questionnaires',
    name: 'AdminQuestionnaires',
    component: AdminQuestionnaires,
    meta: { requiresAdmin: true }
  },
  {
    path: '/admin/statistics',
    name: 'AdminStatistics',
    component: AdminStatistics,
    meta: { requiresAdmin: true }
  }
]

const router = createRouter({
  // 使用Hash模式而不是History模式，避免刷新页面404
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
  const isAdmin = userInfo.is_admin || false
  
  // 需要登录的页面
  const authPages = ['/questionnaire/create', '/questionnaire/edit', '/questionnaire/results']
  
  // 需要管理员权限的页面
  const adminPages = ['/admin', '/admin/users', '/admin/questionnaires', '/admin/statistics']
  
  // 检查是否是需要登录的页面
  const needAuth = authPages.some(path => to.path.startsWith(path))
  
  // 检查是否是需要管理员权限的页面
  const needAdmin = to.meta.requiresAdmin || adminPages.some(path => to.path.startsWith(path))
  
  if (needAdmin) {
    if (!token) {
      next('/login')
    } else if (!isAdmin) {
      // 如果已登录但不是管理员
      next('/')
    } else {
      next()
    }
  } else if (needAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router 