package handlers

import (
	"log"
	"questionnaire-system/backend/database"
	"questionnaire-system/backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminHandler 处理管理员相关请求
type AdminHandler struct {
	DB *database.Database
}

// NewAdminHandler 创建管理员处理器
func NewAdminHandler(db *database.Database) *AdminHandler {
	return &AdminHandler{DB: db}
}

// GetAllUsers 获取所有用户信息
func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	log.Println("管理员请求: 获取所有用户")

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 查询用户总数
	var total int64
	h.DB.Model(&models.User{}).Count(&total)

	// 查询用户列表（不返回密码字段）
	var users []models.User
	h.DB.Select("id, username, email, phone, is_admin, created_at, updated_at").
		Offset(offset).Limit(pageSize).
		Order("id desc").
		Find(&users)

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"users":     users,
		},
	})
}

// GetUserDetail 获取用户详情
func (h *AdminHandler) GetUserDetail(c *gin.Context) {
	log.Println("管理员请求: 获取用户详情")

	// 获取用户ID
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "缺少用户ID",
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "无效的用户ID",
		})
		return
	}

	// 查询用户
	var user models.User
	result := h.DB.Select("id, username, email, phone, is_admin, created_at, updated_at").
		First(&user, id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	// 查询用户创建的问卷数量
	var questionnaireCount int64
	h.DB.Model(&models.Questionnaire{}).Where("created_by = ?", id).Count(&questionnaireCount)

	// 查询用户提交的答卷数量
	var submissionCount int64
	h.DB.Model(&models.Submission{}).Where("user_id = ?", id).Count(&submissionCount)

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"user":                user,
			"questionnaire_count": questionnaireCount,
			"submission_count":    submissionCount,
		},
	})
}

// UpdateUser 更新用户信息
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	log.Println("管理员请求: 更新用户信息")

	// 解析请求
	var request struct {
		ID      uint   `json:"id"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		IsAdmin bool   `json:"is_admin"`
	}

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "无效的请求数据",
		})
		return
	}

	// 查询用户
	var user models.User
	result := h.DB.First(&user, request.ID)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	// 更新用户信息
	user.Email = request.Email
	user.Phone = request.Phone
	user.IsAdmin = request.IsAdmin
	user.UpdatedAt = time.Now()

	result = h.DB.Save(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "更新用户失败: " + result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "用户更新成功",
	})
}

// DeleteUser 删除用户
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	log.Println("管理员请求: 删除用户")

	// 获取用户ID
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "缺少用户ID",
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "无效的用户ID",
		})
		return
	}

	// 查询用户
	var user models.User
	result := h.DB.First(&user, id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	// 不允许删除管理员
	if user.IsAdmin {
		c.JSON(403, gin.H{
			"success": false,
			"message": "不允许删除管理员账户",
		})
		return
	}

	// 开始事务
	tx := h.DB.Begin()

	// 删除用户创建的问卷
	if err := tx.Where("created_by = ?", id).Delete(&models.Questionnaire{}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"success": false,
			"message": "删除用户问卷失败: " + err.Error(),
		})
		return
	}

	// 删除用户提交的答案
	if err := tx.Where("user_id = ?", id).Delete(&models.Answer{}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"success": false,
			"message": "删除用户答案失败: " + err.Error(),
		})
		return
	}

	// 删除用户提交记录
	if err := tx.Where("user_id = ?", id).Delete(&models.Submission{}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"success": false,
			"message": "删除用户提交记录失败: " + err.Error(),
		})
		return
	}

	// 删除用户
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"success": false,
			"message": "删除用户失败: " + err.Error(),
		})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "删除用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "用户删除成功",
	})
}

// GetAllQuestionnaires 获取所有问卷
func (h *AdminHandler) GetAllQuestionnaires(c *gin.Context) {
	log.Println("管理员请求: 获取所有问卷")

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 查询问卷总数
	var total int64
	h.DB.Model(&models.Questionnaire{}).Count(&total)

	// 查询问卷列表
	var questionnaires []models.Questionnaire
	h.DB.Offset(offset).Limit(pageSize).
		Order("id desc").
		Find(&questionnaires)

	// 获取每个问卷的创建者信息和提交数量
	type QuestionnaireWithInfo struct {
		Questionnaire   models.Questionnaire `json:"questionnaire"`
		CreatorName     string               `json:"creator_name"`
		SubmissionCount int64                `json:"submission_count"`
		QuestionCount   int64                `json:"question_count"`
	}

	var result []QuestionnaireWithInfo
	for _, q := range questionnaires {
		var creator models.User
		h.DB.Select("username").First(&creator, q.CreatedBy)

		var submissionCount int64
		h.DB.Model(&models.Submission{}).Where("questionnaire_id = ?", q.ID).Count(&submissionCount)

		var questionCount int64
		h.DB.Model(&models.Question{}).Where("questionnaire_id = ?", q.ID).Count(&questionCount)

		result = append(result, QuestionnaireWithInfo{
			Questionnaire:   q,
			CreatorName:     creator.Username,
			SubmissionCount: submissionCount,
			QuestionCount:   questionCount,
		})
	}

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"total":          total,
			"page":           page,
			"page_size":      pageSize,
			"questionnaires": result,
		},
	})
}

// GetSystemStatistics 获取系统统计信息
func (h *AdminHandler) GetSystemStatistics(c *gin.Context) {
	log.Println("管理员请求: 获取系统统计信息")

	// 用户统计
	var userCount int64
	h.DB.Model(&models.User{}).Count(&userCount)

	var adminCount int64
	h.DB.Model(&models.User{}).Where("is_admin = ?", true).Count(&adminCount)

	// 问卷统计
	var questionnaireCount int64
	h.DB.Model(&models.Questionnaire{}).Count(&questionnaireCount)

	var publishedQuestionnaireCount int64
	h.DB.Model(&models.Questionnaire{}).Where("is_published = ?", true).Count(&publishedQuestionnaireCount)

	// 问题统计
	var questionCount int64
	h.DB.Model(&models.Question{}).Count(&questionCount)

	// 提交统计
	var submissionCount int64
	h.DB.Model(&models.Submission{}).Count(&submissionCount)

	var answerCount int64
	h.DB.Model(&models.Answer{}).Count(&answerCount)

	// 活跃度统计 - 最近7天的提交数量
	var recentSubmissionCount int64
	h.DB.Model(&models.Submission{}).
		Where("submitted_at > ?", time.Now().AddDate(0, 0, -7)).
		Count(&recentSubmissionCount)

	// 返回统计数据
	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"user_statistics": gin.H{
				"total_users":  userCount,
				"admin_users":  adminCount,
				"normal_users": userCount - adminCount,
			},
			"questionnaire_statistics": gin.H{
				"total_questionnaires":       questionnaireCount,
				"published_questionnaires":   publishedQuestionnaireCount,
				"unpublished_questionnaires": questionnaireCount - publishedQuestionnaireCount,
				"total_questions":            questionCount,
			},
			"submission_statistics": gin.H{
				"total_submissions":              submissionCount,
				"total_answers":                  answerCount,
				"recent_submissions":             recentSubmissionCount,
				"average_answers_per_submission": float64(answerCount) / float64(submissionCount),
			},
		},
	})
}

// GetQuestionnaireSubmissions 获取问卷提交详情
func (h *AdminHandler) GetQuestionnaireSubmissions(c *gin.Context) {
	log.Println("管理员请求: 获取问卷提交详情")

	// 获取问卷ID
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "缺少问卷ID",
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "无效的问卷ID",
		})
		return
	}

	// 查询问卷
	var questionnaire models.Questionnaire
	result := h.DB.First(&questionnaire, id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "问卷不存在",
		})
		return
	}

	// 查询问题
	var questions []models.Question
	h.DB.Where("questionnaire_id = ?", id).Order("sort").Find(&questions)

	// 查询提交记录
	var submissions []models.Submission
	h.DB.Where("questionnaire_id = ?", id).Order("submitted_at desc").Find(&submissions)

	// 构建提交详情
	type SubmissionDetail struct {
		Submission models.Submission `json:"submission"`
		User       struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
		Answers []struct {
			QuestionID uint   `json:"question_id"`
			Content    string `json:"content"`
		} `json:"answers"`
	}

	var submissionDetails []SubmissionDetail
	for _, submission := range submissions {
		var detail SubmissionDetail
		detail.Submission = submission

		// 获取用户信息
		var user models.User
		h.DB.Select("id, username").First(&user, submission.UserID)
		detail.User.ID = user.ID
		detail.User.Username = user.Username

		// 获取答案
		var answers []models.Answer
		h.DB.Where("user_id = ? AND question_id IN (?)",
			submission.UserID,
			h.DB.Model(&models.Question{}).Select("id").Where("questionnaire_id = ?", id),
		).Find(&answers)

		for _, answer := range answers {
			detail.Answers = append(detail.Answers, struct {
				QuestionID uint   `json:"question_id"`
				Content    string `json:"content"`
			}{
				QuestionID: answer.QuestionID,
				Content:    answer.Content,
			})
		}

		submissionDetails = append(submissionDetails, detail)
	}

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"questionnaire":      questionnaire,
			"questions":          questions,
			"submission_details": submissionDetails,
		},
	})
}
