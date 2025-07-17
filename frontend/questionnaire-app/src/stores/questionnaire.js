import { defineStore } from 'pinia'
import api from '../services/api'

// 定义问卷状态管理
export const useQuestionnaireStore = defineStore('questionnaire', {
  state: () => ({
    questionnaireList: [],
    currentQuestionnaire: null,
    totalCount: 0,
    currentPage: 1,
    pageSize: 10
  }),
  
  actions: {
    // 获取问卷列表
    async getQuestionnaireList(page = 1) {
      try {
        // 获取当前用户ID
        const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
        const userId = userInfo.id || 0
        const isAdmin = userInfo.is_admin || false
        
        // 添加用户ID参数，后端可以据此返回适合该用户的问卷列表
        const response = await api.get(`/questionnaire/list?page=${page}&user_id=${userId}`)
        
        // 确保返回的数据是有效的，如果不是则使用默认值
        let questionnaires = response?.data?.questionnaires || []
        
        // 如果后端没有进行过滤，前端再次过滤：普通用户只能看到已发布的问卷或自己创建的问卷
        if (!isAdmin && userId > 0) {
          questionnaires = questionnaires.filter(q => 
            q.is_published || q.created_by === userId
          )
        } else if (!isAdmin) {
          // 未登录用户只能看到已发布的问卷
          questionnaires = questionnaires.filter(q => q.is_published)
        }
        
        this.questionnaireList = questionnaires
        this.totalCount = response?.data?.total || questionnaires.length
        this.currentPage = response?.data?.page || page
        this.pageSize = response?.data?.page_size || 10
        
        return Promise.resolve(response)
      } catch (error) {
        console.error('获取问卷列表失败:', error)
        // 出错时设置默认值，确保UI不会崩溃
        this.questionnaireList = []
        this.totalCount = 0
        
        return Promise.reject(error)
      }
    },
    
    // 获取问卷详情
    async getQuestionnaireDetail(id) {
      try {
        const response = await api.get(`/questionnaire/detail?id=${id}`)
        
        // 检查响应数据结构
        if (!response || !response.data) {
          console.error('问卷详情数据格式错误:', response)
          throw new Error('返回数据格式错误')
        }
        
        // 保存当前问卷数据
        this.currentQuestionnaire = response.data.questionnaire || null
        
        return Promise.resolve(response)
      } catch (error) {
        console.error('获取问卷详情失败:', error)
        this.currentQuestionnaire = null
        
        // 提供更具体的错误信息
        let errorMessage = '未知错误'
        if (error.response) {
          errorMessage = error.response.data?.message || `服务器错误 (${error.response.status})`
        } else if (error.request) {
          errorMessage = '服务器无响应'
        } else {
          errorMessage = error.message || '获取问卷详情失败'
        }
        
        return Promise.reject(new Error(errorMessage))
      }
    },
    
    // 创建问卷
    async createQuestionnaire(questionnaireData) {
      try {
        // 确保日期格式正确
        if (questionnaireData.start_time && typeof questionnaireData.start_time === 'string') {
          questionnaireData.start_time = new Date(questionnaireData.start_time).toISOString();
        }
        
        if (questionnaireData.end_time && typeof questionnaireData.end_time === 'string') {
          questionnaireData.end_time = new Date(questionnaireData.end_time).toISOString();
        }
        
        // 确保问题数据格式正确
        if (questionnaireData.questions && Array.isArray(questionnaireData.questions)) {
          questionnaireData.questions = questionnaireData.questions.map((q, index) => ({
            ...q,
            sort: index
          }));
        }
        
        console.log('发送创建问卷请求:', JSON.stringify(questionnaireData));
        
        const response = await api.post('/questionnaire/create', questionnaireData, { timeout: 20000 })
        return Promise.resolve(response)
      } catch (error) {
        console.error('创建问卷失败:', error)
        return Promise.reject(error)
      }
    },
    
    // 提交问卷答案
    async submitAnswer(answerData) {
      try {
        // 确保答案字段名称正确
        if (answerData.answers && Array.isArray(answerData.answers)) {
          // 检查并修正字段名
          answerData.answers = answerData.answers.map(answer => {
            // 如果使用了content字段而不是answer_content，进行转换
            if (answer.content !== undefined && answer.answer_content === undefined) {
              return {
                ...answer,
                answer_content: answer.content
              };
            }
            return answer;
          });
        }
        
        console.log('提交问卷答案数据:', JSON.stringify(answerData));
        
        const response = await api.post('/questionnaire/submit', answerData)
        return Promise.resolve(response)
      } catch (error) {
        console.error('提交问卷答案失败:', error)
        
        // 提供更具体的错误信息
        let errorMessage = '未知错误'
        if (error.response) {
          // 检查是否是重复提交错误
          if (error.response.status === 409) {
            errorMessage = '您已经提交过该问卷，不能重复提交'
          } else {
            errorMessage = error.response.data?.message || `服务器错误 (${error.response.status})`
          }
        } else if (error.request) {
          errorMessage = '服务器无响应'
        } else {
          errorMessage = error.message || '提交问卷答案失败'
        }
        
        return Promise.reject(new Error(errorMessage))
      }
    },
    
    // 更新问卷状态（发布/取消发布）
    async updateQuestionnaireStatus(id, isPublished) {
      try {
        // 添加更详细的日志
        console.log('更新问卷状态请求参数:', { id, is_published: isPublished })
        
        const response = await api.put('/questionnaire/update-status', { id, is_published: isPublished })
        
        // 添加响应日志
        console.log('更新问卷状态响应:', response)
        
        return Promise.resolve(response)
      } catch (error) {
        console.error('更新问卷状态失败:', error)
        
        // 提供更详细的错误信息
        let errorMessage = '未知错误'
        if (error.response) {
          console.error('错误响应数据:', error.response.data)
          console.error('错误状态码:', error.response.status)
          errorMessage = error.response.data?.message || `服务器错误 (${error.response.status})`
        } else if (error.request) {
          errorMessage = '服务器无响应'
        } else {
          errorMessage = error.message || '更新问卷状态失败'
        }
        
        return Promise.reject(new Error(errorMessage))
      }
    },
    
    // 编辑问卷
    async updateQuestionnaire(questionnaireData) {
      try {
        // 确保日期格式正确
        if (questionnaireData.start_time && typeof questionnaireData.start_time === 'string') {
          questionnaireData.start_time = new Date(questionnaireData.start_time).toISOString();
        }
        
        if (questionnaireData.end_time && typeof questionnaireData.end_time === 'string') {
          questionnaireData.end_time = new Date(questionnaireData.end_time).toISOString();
        }
        
        // 确保问题数据格式正确
        if (questionnaireData.questions && Array.isArray(questionnaireData.questions)) {
          questionnaireData.questions = questionnaireData.questions.map((q, index) => ({
            ...q,
            sort: index
          }));
        }
        
        console.log('发送更新问卷请求:', JSON.stringify(questionnaireData));
        
        const response = await api.put('/questionnaire/update', questionnaireData, { timeout: 20000 })
        return Promise.resolve(response)
      } catch (error) {
        console.error('更新问卷失败:', error)
        
        // 提供更具体的错误信息
        let errorMessage = '未知错误'
        if (error.response) {
          errorMessage = error.response.data?.message || `服务器错误 (${error.response.status})`
        } else if (error.request) {
          errorMessage = '服务器无响应'
        } else {
          errorMessage = error.message || '更新问卷失败'
        }
        
        return Promise.reject(new Error(errorMessage))
      }
    },
    
    // 删除问卷
    async deleteQuestionnaire(id) {
      try {
        const response = await api.delete(`/questionnaire/delete?id=${id}`)
        return Promise.resolve(response)
      } catch (error) {
        console.error('删除问卷失败:', error)
        
        // 提供更具体的错误信息
        let errorMessage = '未知错误'
        if (error.response) {
          errorMessage = error.response.data?.message || `服务器错误 (${error.response.status})`
        } else if (error.request) {
          errorMessage = '服务器无响应'
        } else {
          errorMessage = error.message || '删除问卷失败'
        }
        
        return Promise.reject(new Error(errorMessage))
      }
    },
    
    // 获取问卷结果
    async getQuestionnaireResults(id) {
      try {
        // 获取用户信息
        const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
        const userId = userInfo.id || 0
        
        if (userId <= 0) {
          throw new Error('未登录或登录已过期')
        }
        
        const response = await api.get(`/questionnaire/results?id=${id}&user_id=${userId}`)
        return Promise.resolve(response)
      } catch (error) {
        console.error('获取问卷结果失败:', error)
        
        // 提供更具体的错误信息
        let errorMessage = '未知错误'
        if (error.response) {
          errorMessage = error.response.data?.message || `服务器错误 (${error.response.status})`
        } else if (error.request) {
          errorMessage = '服务器无响应'
        } else {
          errorMessage = error.message || '获取问卷结果失败'
        }
        
        return Promise.reject(new Error(errorMessage))
      }
    },
    
    // 检查用户是否已提交过问卷
    async checkSubmission(questionnaireId, userId) {
      try {
        const response = await api.get(`/questionnaire/check-submission?questionnaire_id=${questionnaireId}&user_id=${userId}`)
        return Promise.resolve(response)
      } catch (error) {
        console.error('检查提交状态失败:', error)
        return Promise.reject(error)
      }
    }
  }
}) 