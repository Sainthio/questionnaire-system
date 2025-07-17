package handlers

import (
	"log"
	"questionnaire-system/backend/database"
	"questionnaire-system/backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler 处理用户相关请求
type UserHandler struct {
	DB *database.Database
}

// NewUserHandler 创建用户处理器
func NewUserHandler(db *database.Database) *UserHandler {
	return &UserHandler{DB: db}
}

// Register 注册用户
func (h *UserHandler) Register(c *gin.Context) {
	log.Println("收到注册请求")

	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Printf("解析请求数据失败: %v", err)
		c.JSON(400, gin.H{
			"success": false,
			"error":   "无效的请求数据",
		})
		return
	}

	log.Printf("注册信息: 用户名=%s, 邮箱=%s", user.Username, user.Email)

	// 检查用户名是否已存在
	var existingUser models.User
	result := h.DB.Where("username = ?", user.Username).First(&existingUser)
	if result.RowsAffected > 0 {
		log.Printf("用户名已存在: %s", user.Username)
		c.JSON(400, gin.H{
			"success": false,
			"error":   "用户名已存在",
		})
		return
	}

	// 检查邮箱是否已存在
	result = h.DB.Where("email = ?", user.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		log.Printf("邮箱已存在: %s", user.Email)
		c.JSON(400, gin.H{
			"success": false,
			"error":   "邮箱已存在",
		})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"error":   "密码加密失败",
		})
		return
	}
	user.Password = string(hashedPassword)
	log.Printf("密码加密结果: %s", user.Password)

	// 设置创建时间
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// 保存用户
	result = h.DB.Create(&user)
	if result.Error != nil {
		log.Printf("创建用户失败: %v", result.Error)
		c.JSON(500, gin.H{
			"success": false,
			"error":   "创建用户失败: " + result.Error.Error(),
		})
		return
	}

	log.Printf("用户注册成功: ID=%d, 用户名=%s", user.ID, user.Username)

	// 返回用户信息（不包含密码）
	user.Password = ""
	c.JSON(201, gin.H{
		"success": true,
		"message": "注册成功",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	log.Println("收到登录请求")

	// 解析登录信息
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := c.ShouldBindJSON(&loginRequest)
	if err != nil {
		log.Printf("解析请求数据失败: %v", err)
		c.JSON(400, gin.H{
			"success": false,
			"error":   "无效的请求数据",
		})
		return
	}

	log.Printf("登录尝试: 用户名=%s, 密码长度=%d", loginRequest.Username, len(loginRequest.Password))

	// 查询用户
	var user models.User
	result := h.DB.Where("username = ?", loginRequest.Username).First(&user)
	if result.Error != nil {
		log.Printf("用户不存在: %s", loginRequest.Username)
		c.JSON(401, gin.H{
			"success": false,
			"error":   "用户名或密码错误",
		})
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		log.Printf("密码错误: 用户=%s, 错误=%v", loginRequest.Username, err)
		c.JSON(401, gin.H{
			"success": false,
			"error":   "用户名或密码错误",
		})
		return
	}

	log.Printf("用户登录成功: ID=%d, 用户名=%s", user.ID, user.Username)

	// 生成简单的会话令牌（实际项目中应使用JWT或其他认证机制）
	token := "token_" + user.Username + "_" + time.Now().Format("20060102150405")

	c.JSON(200, gin.H{
		"success":  true,
		"message":  "登录成功",
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"token":    token,
	})
}

// ResetPassword 重置密码（用于测试）
func (h *UserHandler) ResetPassword(c *gin.Context) {
	log.Println("收到重置密码请求")

	// 解析请求
	var resetRequest struct {
		Username    string `json:"username"`
		NewPassword string `json:"new_password"`
	}

	err := c.ShouldBindJSON(&resetRequest)
	if err != nil {
		log.Printf("解析请求数据失败: %v", err)
		c.JSON(400, gin.H{
			"success": false,
			"error":   "无效的请求数据",
		})
		return
	}

	log.Printf("重置密码: 用户名=%s, 新密码长度=%d", resetRequest.Username, len(resetRequest.NewPassword))

	// 查询用户
	var user models.User
	result := h.DB.Where("username = ?", resetRequest.Username).First(&user)
	if result.Error != nil {
		log.Printf("用户不存在: %s", resetRequest.Username)

		// 如果用户不存在，创建一个新用户（特殊情况，仅用于测试）
		if resetRequest.Username == "testuser" {
			// 创建测试用户
			newUser := models.User{
				Username:  resetRequest.Username,
				Email:     "test@example.com",
				IsAdmin:   false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			// 加密密码
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetRequest.NewPassword), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("密码加密失败: %v", err)
				c.JSON(500, gin.H{
					"success": false,
					"error":   "密码加密失败",
				})
				return
			}
			newUser.Password = string(hashedPassword)

			// 保存用户
			result = h.DB.Create(&newUser)
			if result.Error != nil {
				log.Printf("创建测试用户失败: %v", result.Error)
				c.JSON(500, gin.H{
					"success": false,
					"error":   "创建测试用户失败",
				})
				return
			}

			log.Printf("创建测试用户成功: ID=%d, 用户名=%s", newUser.ID, newUser.Username)

			c.JSON(200, gin.H{
				"success": true,
				"message": "测试用户创建成功",
			})
			return
		}

		c.JSON(404, gin.H{
			"success": false,
			"error":   "用户不存在",
		})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetRequest.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		c.JSON(500, gin.H{
			"success": false,
			"error":   "密码加密失败",
		})
		return
	}

	// 更新密码
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()
	result = h.DB.Save(&user)
	if result.Error != nil {
		log.Printf("更新密码失败: %v", result.Error)
		c.JSON(500, gin.H{
			"success": false,
			"error":   "更新密码失败",
		})
		return
	}

	log.Printf("密码重置成功: 用户=%s, 新密码哈希=%s", user.Username, user.Password)

	// 测试密码验证
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(resetRequest.NewPassword))
	if err != nil {
		log.Printf("警告：密码验证测试失败: %v", err)
	} else {
		log.Printf("密码验证测试成功")
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "密码重置成功",
	})
}
