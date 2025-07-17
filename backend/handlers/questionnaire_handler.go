package handlers

import (
	"log"
	"net/http"
	"questionnaire-system/backend/database"
	"questionnaire-system/backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// QuestionnaireHandler 处理问卷相关请求
type QuestionnaireHandler struct {
	DB *database.Database
}

// NewQuestionnaireHandler 创建问卷处理器
func NewQuestionnaireHandler(db *database.Database) *QuestionnaireHandler {
	return &QuestionnaireHandler{DB: db}
}

// CreateQuestionnaire 创建问卷
func (h *QuestionnaireHandler) CreateQuestionnaire(c *gin.Context) {
	log.Println("收到创建问卷请求")

	// 解析请求数据
	type QuestionRequest struct {
		Title    string `json:"title"`
		Type     string `json:"type"`
		Required bool   `json:"required"`
		Options  string `json:"options"`
		Sort     int    `json:"sort"`
	}

	type QuestionnaireRequest struct {
		Title       string            `json:"title"`
		Description string            `json:"description"`
		CreatedBy   uint              `json:"created_by"`
		StartTime   time.Time         `json:"start_time"`
		EndTime     time.Time         `json:"end_time"`
		IsPublished bool              `json:"is_published"`
		Questions   []QuestionRequest `json:"questions"`
	}

	var request QuestionnaireRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Printf("解析请求数据失败: %v", err)
		c.JSON(400, gin.H{
			"success": false,
			"message": "无效的请求数据",
		})
		return
	}

	log.Printf("问卷数据: 标题=%s, 描述=%s, 创建者ID=%d, 问题数=%d, 是否发布=%v",
		request.Title, request.Description, request.CreatedBy, len(request.Questions), request.IsPublished)

	// 验证创建者ID
	if request.CreatedBy <= 0 {
		log.Printf("无效的创建者ID: %d", request.CreatedBy)
		c.JSON(400, gin.H{
			"success": false,
			"message": "无效的创建者ID",
		})
		return
	}

	// 验证问题数量
	if len(request.Questions) == 0 {
		log.Printf("问卷没有问题")
		c.JSON(400, gin.H{
			"success": false,
			"message": "问卷必须包含至少一个问题",
		})
		return
	}

	// 创建问卷对象
	questionnaire := models.Questionnaire{
		Title:       request.Title,
		Description: request.Description,
		CreatedBy:   request.CreatedBy,
		StartTime:   request.StartTime,
		EndTime:     request.EndTime,
		IsPublished: request.IsPublished,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 开始事务
	tx := h.DB.Begin()

	// 保存问卷
	if err := tx.Create(&questionnaire).Error; err != nil {
		tx.Rollback()
		log.Printf("创建问卷失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "创建问卷失败: " + err.Error(),
		})
		return
	}

	log.Printf("问卷创建成功: ID=%d, 标题=%s", questionnaire.ID, questionnaire.Title)

	// 保存问题
	var questions []models.Question
	for i, q := range request.Questions {
		question := models.Question{
			QuestionnaireID: questionnaire.ID,
			Title:           q.Title,
			Type:            q.Type,
			Required:        q.Required,
			Options:         q.Options,
			Sort:            i,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err := tx.Create(&question).Error; err != nil {
			tx.Rollback()
			log.Printf("创建问题失败: %v", err)
			c.JSON(500, gin.H{
				"success": false,
				"message": "创建问题失败: " + err.Error(),
			})
			return
		}

		questions = append(questions, question)
		log.Printf("创建问题: ID=%d, 标题=%s, 类型=%s, 选项=%s",
			question.ID, question.Title, question.Type, question.Options)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "创建问卷失败: " + err.Error(),
		})
		return
	}

	log.Printf("问卷创建完成: ID=%d, 标题=%s, 问题数=%d", questionnaire.ID, questionnaire.Title, len(questions))

	// 返回创建的问卷和问题
	c.JSON(201, gin.H{
		"success": true,
		"message": "问卷创建成功",
		"data": map[string]interface{}{
			"questionnaire": questionnaire,
			"questions":     questions,
		},
	})
}

// GetQuestionnaireDetail 获取问卷详情
func (h *QuestionnaireHandler) GetQuestionnaireDetail(c *gin.Context) {
	log.Printf("获取问卷详情: ID=%s", c.Query("id"))

	// 从URL获取问卷ID
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

	log.Printf("获取问卷详情: ID=%d", id)

	// 查询问卷
	var questionnaire models.Questionnaire
	result := h.DB.First(&questionnaire, id)
	if result.Error != nil {
		log.Printf("问卷不存在: ID=%d", id)
		c.JSON(404, gin.H{
			"success": false,
			"message": "问卷不存在",
		})
		return
	}

	// 查询问卷的问题
	var questions []models.Question
	h.DB.Where("questionnaire_id = ?", id).Find(&questions)

	// 构造响应数据
	type Response struct {
		Questionnaire models.Questionnaire `json:"questionnaire"`
		Questions     []models.Question    `json:"questions"`
	}

	response := Response{
		Questionnaire: questionnaire,
		Questions:     questions,
	}

	// 返回问卷详情
	c.JSON(200, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetQuestionnaires 获取问卷列表
func (h *QuestionnaireHandler) GetQuestionnaires(c *gin.Context) {
	log.Println("获取问卷列表")

	// 获取查询参数
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	userIDStr := c.Query("user_id")

	// 设置默认值
	page := 1
	pageSize := 10
	var userID uint = 0

	// 解析参数
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	if userIDStr != "" {
		if uid, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			userID = uint(uid)
		}
	}

	log.Printf("查询参数: 页码=%d, 每页数量=%d, 用户ID=%d", page, pageSize, userID)

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	query := h.DB.Model(&models.Questionnaire{})

	// 根据用户ID过滤
	if userID > 0 {
		// 如果指定了用户ID，则查询该用户创建的问卷或已发布的问卷
		query = query.Where("is_published = ? OR created_by = ?", true, userID)
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询分页数据
	var questionnaires []models.Questionnaire
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&questionnaires)

	log.Printf("查询到 %d 个问卷", len(questionnaires))
	for i, q := range questionnaires {
		log.Printf("问卷 %d: ID=%d, 标题=%s, 创建者=%d, 已发布=%v", i+1, q.ID, q.Title, q.CreatedBy, q.IsPublished)
	}

	// 构造响应
	type QuestionnaireWithInfo struct {
		Questionnaire models.Questionnaire `json:"questionnaire"`
		CreatorName   string               `json:"creator_name"`
	}

	var questionnaireWithInfo []QuestionnaireWithInfo
	for _, q := range questionnaires {
		var creator models.User
		h.DB.Select("username").First(&creator, q.CreatedBy)

		questionnaireWithInfo = append(questionnaireWithInfo, QuestionnaireWithInfo{
			Questionnaire: q,
			CreatorName:   creator.Username,
		})
	}

	type Response struct {
		Total          int64                   `json:"total"`
		Page           int                     `json:"page"`
		PageSize       int                     `json:"page_size"`
		Questionnaires []QuestionnaireWithInfo `json:"questionnaires"`
	}

	response := Response{
		Total:          total,
		Page:           page,
		PageSize:       pageSize,
		Questionnaires: questionnaireWithInfo,
	}

	// 返回JSON响应
	c.JSON(200, gin.H{
		"success": true,
		"data":    response,
	})
}

// SubmitQuestionnaire 提交问卷答案
func (h *QuestionnaireHandler) SubmitQuestionnaire(c *gin.Context) {
	log.Println("收到提交问卷请求")

	// 解析请求数据
	type AnswerRequest struct {
		QuestionnaireID uint            `json:"questionnaire_id"`
		UserID          uint            `json:"user_id"`
		Answers         []models.Answer `json:"answers"`
	}

	var request AnswerRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Printf("解析请求数据失败: %v", err)
		c.JSON(400, gin.H{
			"success": false,
			"message": "无效的请求数据",
		})
		return
	}

	log.Printf("提交数据: 问卷ID=%d, 用户ID=%d, 答案数量=%d",
		request.QuestionnaireID, request.UserID, len(request.Answers))

	// 验证问卷ID和用户ID
	if request.QuestionnaireID <= 0 || request.UserID <= 0 {
		log.Printf("无效的问卷ID或用户ID: 问卷ID=%d, 用户ID=%d", request.QuestionnaireID, request.UserID)
		c.JSON(400, gin.H{
			"success": false,
			"message": "无效的问卷ID或用户ID",
		})
		return
	}

	// 检查问卷是否存在
	var questionnaire models.Questionnaire
	if err := h.DB.First(&questionnaire, request.QuestionnaireID).Error; err != nil {
		log.Printf("问卷不存在: ID=%d", request.QuestionnaireID)
		c.JSON(404, gin.H{
			"success": false,
			"message": "问卷不存在",
		})
		return
	}

	// 检查用户是否已提交过该问卷
	var existingSubmission models.Submission
	result := h.DB.Where("questionnaire_id = ? AND user_id = ?", request.QuestionnaireID, request.UserID).First(&existingSubmission)
	if result.Error == nil {
		log.Printf("用户已提交过该问卷: 问卷ID=%d, 用户ID=%d", request.QuestionnaireID, request.UserID)
		c.JSON(409, gin.H{
			"success": false,
			"message": "您已经提交过该问卷，不能重复提交",
		})
		return
	}

	// 开始事务
	tx := h.DB.Begin()

	// 创建提交记录
	submission := models.Submission{
		QuestionnaireID: request.QuestionnaireID,
		UserID:          request.UserID,
		SubmittedAt:     time.Now(),
		IPAddress:       c.ClientIP(),
	}

	if err := tx.Create(&submission).Error; err != nil {
		tx.Rollback()
		log.Printf("创建提交记录失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "提交问卷失败: " + err.Error(),
		})
		return
	}

	// 保存答案
	for _, answer := range request.Answers {
		answer.UserID = request.UserID
		answer.CreatedAt = time.Now()

		if err := tx.Create(&answer).Error; err != nil {
			tx.Rollback()
			log.Printf("保存答案失败: %v", err)
			c.JSON(500, gin.H{
				"success": false,
				"message": "提交问卷失败: " + err.Error(),
			})
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "提交问卷失败: " + err.Error(),
		})
		return
	}

	log.Printf("问卷提交成功: 问卷ID=%d, 用户ID=%d", request.QuestionnaireID, request.UserID)

	// 返回成功响应
	c.JSON(201, gin.H{
		"success": true,
		"message": "问卷提交成功",
	})
}

// UpdateQuestionnaireStatus 更新问卷状态（发布/取消发布）
func (h *QuestionnaireHandler) UpdateQuestionnaireStatus(c *gin.Context) {
	log.Println("收到更新问卷状态请求")

	// 解析请求数据
	var request struct {
		ID          uint `json:"id"`
		IsPublished bool `json:"is_published"`
	}

	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Printf("解析请求数据失败: %v", err)
		c.JSON(400, gin.H{
			"success": false,
			"message": "无效的请求数据",
		})
		return
	}

	log.Printf("更新问卷状态: ID=%d, 发布状态=%v", request.ID, request.IsPublished)

	// 查询问卷
	var questionnaire models.Questionnaire
	result := h.DB.First(&questionnaire, request.ID)
	if result.Error != nil {
		log.Printf("问卷不存在: ID=%d", request.ID)
		c.JSON(404, gin.H{
			"success": false,
			"message": "问卷不存在",
		})
		return
	}

	// 更新问卷状态
	questionnaire.IsPublished = request.IsPublished
	questionnaire.UpdatedAt = time.Now()

	if err := h.DB.Save(&questionnaire).Error; err != nil {
		log.Printf("更新问卷状态失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "更新问卷状态失败: " + err.Error(),
		})
		return
	}

	log.Printf("问卷状态更新成功: ID=%d, 发布状态=%v", questionnaire.ID, questionnaire.IsPublished)

	// 返回更新后的问卷
	c.JSON(200, gin.H{
		"success": true,
		"message": "问卷状态更新成功",
		"data":    questionnaire,
	})
}

// UpdateQuestionnaire 更新问卷
func (h *QuestionnaireHandler) UpdateQuestionnaire(c *gin.Context) {
	log.Println("收到更新问卷请求")

	// 解析请求数据
	type QuestionRequest struct {
		ID       uint   `json:"id"`
		Title    string `json:"title"`
		Type     string `json:"type"`
		Required bool   `json:"required"`
		Options  string `json:"options"`
		Sort     int    `json:"sort"`
	}

	type QuestionnaireRequest struct {
		ID          uint              `json:"id"`
		Title       string            `json:"title"`
		Description string            `json:"description"`
		CreatedBy   uint              `json:"created_by"`
		StartTime   time.Time         `json:"start_time"`
		EndTime     time.Time         `json:"end_time"`
		IsPublished bool              `json:"is_published"`
		Questions   []QuestionRequest `json:"questions"`
	}

	var request QuestionnaireRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Printf("解析请求数据失败: %v", err)
		c.JSON(400, gin.H{
			"success": false,
			"message": "无效的请求数据",
		})
		return
	}

	log.Printf("更新问卷: ID=%d, 标题=%s, 问题数=%d", request.ID, request.Title, len(request.Questions))

	// 查询问卷
	var questionnaire models.Questionnaire
	result := h.DB.First(&questionnaire, request.ID)
	if result.Error != nil {
		log.Printf("问卷不存在: ID=%d", request.ID)
		c.JSON(404, gin.H{
			"success": false,
			"message": "问卷不存在",
		})
		return
	}

	// 验证权限（只有创建者可以编辑）
	if questionnaire.CreatedBy != request.CreatedBy {
		log.Printf("权限不足: 用户ID=%d, 问卷创建者ID=%d", request.CreatedBy, questionnaire.CreatedBy)
		c.JSON(403, gin.H{
			"success": false,
			"message": "您没有权限编辑此问卷",
		})
		return
	}

	// 验证问卷状态（已发布的问卷不能编辑）
	if questionnaire.IsPublished {
		log.Printf("问卷已发布，不能编辑: ID=%d", request.ID)
		c.JSON(400, gin.H{
			"success": false,
			"message": "已发布的问卷不能编辑",
		})
		return
	}

	// 开始事务
	tx := h.DB.Begin()

	// 更新问卷信息
	questionnaire.Title = request.Title
	questionnaire.Description = request.Description
	questionnaire.StartTime = request.StartTime
	questionnaire.EndTime = request.EndTime
	questionnaire.UpdatedAt = time.Now()

	if err := tx.Save(&questionnaire).Error; err != nil {
		tx.Rollback()
		log.Printf("更新问卷失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "更新问卷失败: " + err.Error(),
		})
		return
	}

	// 删除原有问题
	if err := tx.Where("questionnaire_id = ?", request.ID).Delete(&models.Question{}).Error; err != nil {
		tx.Rollback()
		log.Printf("删除原有问题失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "更新问卷失败: " + err.Error(),
		})
		return
	}

	// 创建新问题
	var questions []models.Question
	for i, q := range request.Questions {
		question := models.Question{
			QuestionnaireID: questionnaire.ID,
			Title:           q.Title,
			Type:            q.Type,
			Required:        q.Required,
			Options:         q.Options,
			Sort:            i,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err := tx.Create(&question).Error; err != nil {
			tx.Rollback()
			log.Printf("创建问题失败: %v", err)
			c.JSON(500, gin.H{
				"success": false,
				"message": "更新问卷失败: " + err.Error(),
			})
			return
		}

		questions = append(questions, question)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "更新问卷失败: " + err.Error(),
		})
		return
	}

	log.Printf("问卷更新成功: ID=%d, 标题=%s, 问题数=%d", questionnaire.ID, questionnaire.Title, len(questions))

	// 返回更新后的问卷和问题
	c.JSON(200, gin.H{
		"success": true,
		"message": "问卷更新成功",
		"data": map[string]interface{}{
			"questionnaire": questionnaire,
			"questions":     questions,
		},
	})
}

// DeleteQuestionnaire 删除问卷
func (h *QuestionnaireHandler) DeleteQuestionnaire(c *gin.Context) {
	log.Println("收到删除问卷请求")

	// 从URL获取问卷ID
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

	log.Printf("删除问卷: ID=%d", id)

	// 查询问卷
	var questionnaire models.Questionnaire
	result := h.DB.First(&questionnaire, id)
	if result.Error != nil {
		log.Printf("问卷不存在: ID=%d", id)
		c.JSON(404, gin.H{
			"success": false,
			"message": "问卷不存在",
		})
		return
	}

	// 开始事务
	tx := h.DB.Begin()

	// 删除相关问题
	if err := tx.Where("questionnaire_id = ?", id).Delete(&models.Question{}).Error; err != nil {
		tx.Rollback()
		log.Printf("删除问题失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "删除问卷失败",
		})
		return
	}

	// 删除相关答案
	if err := tx.Exec("DELETE FROM answers WHERE question_id IN (SELECT id FROM questions WHERE questionnaire_id = ?)", id).Error; err != nil {
		tx.Rollback()
		log.Printf("删除答案失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "删除问卷失败",
		})
		return
	}

	// 删除相关提交记录
	if err := tx.Where("questionnaire_id = ?", id).Delete(&models.Submission{}).Error; err != nil {
		tx.Rollback()
		log.Printf("删除提交记录失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "删除问卷失败",
		})
		return
	}

	// 删除问卷
	if err := tx.Delete(&questionnaire).Error; err != nil {
		tx.Rollback()
		log.Printf("删除问卷失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "删除问卷失败",
		})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "删除问卷失败",
		})
		return
	}

	log.Printf("问卷删除成功: ID=%d", id)

	// 返回成功响应
	c.JSON(200, gin.H{
		"success": true,
		"message": "问卷删除成功",
	})
}

// GetQuestionnaireResults 获取问卷填写结果
func (h *QuestionnaireHandler) GetQuestionnaireResults(c *gin.Context) {
	log.Println("收到获取问卷结果请求")

	// 从URL获取问卷ID
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

	// 验证认证信息
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Println("缺少认证信息")
		c.JSON(401, gin.H{
			"success": false,
			"message": "未授权访问",
		})
		return
	}

	// 简单验证token格式
	if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		log.Println("认证格式错误")
		c.JSON(401, gin.H{
			"success": false,
			"message": "认证格式错误",
		})
		return
	}

	// 获取用户ID
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "缺少用户ID",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "无效的用户ID",
		})
		return
	}

	log.Printf("获取问卷结果: ID=%d, 用户ID=%d", id, userID)

	// 查询问卷
	var questionnaire models.Questionnaire
	result := h.DB.First(&questionnaire, id)
	if result.Error != nil {
		log.Printf("问卷不存在: ID=%d", id)
		c.JSON(404, gin.H{
			"success": false,
			"message": "问卷不存在",
		})
		return
	}

	// 验证权限（只有创建者可以查看结果）
	if questionnaire.CreatedBy != uint(userID) {
		// 检查用户是否是管理员
		var user models.User
		userResult := h.DB.First(&user, userID)
		if userResult.Error != nil || !user.IsAdmin {
			log.Printf("权限不足: 用户ID=%d, 问卷创建者ID=%d", userID, questionnaire.CreatedBy)
			c.JSON(403, gin.H{
				"success": false,
				"message": "您没有权限查看此问卷的结果",
			})
			return
		}
	}

	// 查询问卷的问题
	var questions []models.Question
	h.DB.Where("questionnaire_id = ?", id).Find(&questions)

	// 查询提交记录
	var submissions []models.Submission
	h.DB.Where("questionnaire_id = ?", id).Find(&submissions)

	// 查询每个提交的答案
	type SubmissionWithAnswers struct {
		Submission models.Submission `json:"submission"`
		Answers    []struct {
			ID         uint      `json:"id"`
			QuestionID uint      `json:"question_id"`
			Content    string    `json:"content"`
			CreatedAt  time.Time `json:"created_at"`
		} `json:"answers"`
		UserInfo struct {
			Username string `json:"username"`
		} `json:"user_info"`
	}

	var submissionsWithAnswers []SubmissionWithAnswers
	for _, submission := range submissions {
		var answers []models.Answer
		log.Printf("查询提交答案: 用户ID=%d, 问卷ID=%d", submission.UserID, id)

		h.DB.Where("user_id = ? AND question_id IN (SELECT id FROM questions WHERE questionnaire_id = ?)",
			submission.UserID, id).Find(&answers)

		log.Printf("找到答案数量: %d", len(answers))

		// 查询用户信息
		var user struct {
			Username string
		}
		h.DB.Table("users").Select("username").Where("id = ?", submission.UserID).Scan(&user)

		// 转换答案格式，确保字段名一致
		var formattedAnswers []struct {
			ID         uint      `json:"id"`
			QuestionID uint      `json:"question_id"`
			Content    string    `json:"content"`
			CreatedAt  time.Time `json:"created_at"`
		}

		for _, answer := range answers {
			formattedAnswers = append(formattedAnswers, struct {
				ID         uint      `json:"id"`
				QuestionID uint      `json:"question_id"`
				Content    string    `json:"content"`
				CreatedAt  time.Time `json:"created_at"`
			}{
				ID:         answer.ID,
				QuestionID: answer.QuestionID,
				Content:    answer.Content,
				CreatedAt:  answer.CreatedAt,
			})
		}

		submissionWithAnswers := SubmissionWithAnswers{
			Submission: submission,
			Answers:    formattedAnswers,
		}
		submissionWithAnswers.UserInfo.Username = user.Username

		submissionsWithAnswers = append(submissionsWithAnswers, submissionWithAnswers)
	}

	// 构造响应数据
	type Response struct {
		Questionnaire    models.Questionnaire    `json:"questionnaire"`
		Questions        []models.Question       `json:"questions"`
		Submissions      []SubmissionWithAnswers `json:"submissions"`
		TotalSubmissions int                     `json:"total_submissions"`
	}

	response := Response{
		Questionnaire:    questionnaire,
		Questions:        questions,
		Submissions:      submissionsWithAnswers,
		TotalSubmissions: len(submissions),
	}

	// 返回问卷结果
	c.JSON(200, gin.H{
		"success": true,
		"data":    response,
	})
}

// CheckSubmission 检查用户是否已提交过问卷
func (h *QuestionnaireHandler) CheckSubmission(c *gin.Context) {
	log.Println("收到检查提交状态请求")

	// 获取问卷ID和用户ID
	questionnaireIDStr := c.Query("questionnaire_id")
	userIDStr := c.Query("user_id")

	if questionnaireIDStr == "" || userIDStr == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "缺少必要参数",
		})
		return
	}

	questionnaireID, err := strconv.ParseUint(questionnaireIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "无效的问卷ID",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "无效的用户ID",
		})
		return
	}

	log.Printf("检查提交状态: 问卷ID=%d, 用户ID=%d", questionnaireID, userID)

	// 查询提交记录
	var submission models.Submission
	result := h.DB.Where("questionnaire_id = ? AND user_id = ?", questionnaireID, userID).First(&submission)

	// 构造响应
	hasSubmitted := result.Error == nil
	log.Printf("查询结果: 是否已提交=%v", hasSubmitted)

	response := gin.H{
		"success":       true,
		"has_submitted": hasSubmitted,
	}

	// 如果已提交，添加提交信息
	if hasSubmitted {
		response["submission"] = submission
	} else {
		response["submission"] = nil
	}

	c.JSON(200, response)
}

// GetSystemStats 获取系统统计数据
func (h *QuestionnaireHandler) GetSystemStats(c *gin.Context) {
	// 查询注册用户数
	var userCount int64
	result := h.DB.Model(&models.User{}).Count(&userCount)
	if result.Error != nil {
		log.Printf("获取用户数量失败: %v", result.Error)
		userCount = 0
	}

	// 查询问卷总数
	var questionnaireCount int64
	result = h.DB.Model(&models.Questionnaire{}).Count(&questionnaireCount)
	if result.Error != nil {
		log.Printf("获取问卷数量失败: %v", result.Error)
		questionnaireCount = 0
	}

	// 查询提交总数
	var submissionCount int64
	result = h.DB.Model(&models.Submission{}).Count(&submissionCount)
	if result.Error != nil {
		log.Printf("获取提交数量失败: %v", result.Error)
		submissionCount = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"code":                0,
		"message":             "获取统计数据成功",
		"user_count":          userCount,
		"questionnaire_count": questionnaireCount,
		"submission_count":    submissionCount,
	})
}
