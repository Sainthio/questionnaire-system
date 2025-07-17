package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"questionnaire-system/backend/config"
	"questionnaire-system/backend/database"
	"questionnaire-system/backend/handlers"
	"questionnaire-system/backend/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

// 日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 使用Gin的日志系统而不是标准输出
		gin.DefaultWriter.Write([]byte(
			"| " +
				time.Now().Format("2006/01/02 - 15:04:05") +
				" | " +
				reqMethod +
				" | " +
				reqUri +
				" | " +
				clientIP +
				" | " +
				time.Duration(latencyTime).String() +
				" | " +
				" Status: " +
				string(statusCode) +
				" |\n",
		))
	}
}

func main() {
	// 设置日志输出到文件和控制台
	logFile, err := os.OpenFile("backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	defer logFile.Close()

	// 设置日志同时输出到文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	gin.DefaultWriter = multiWriter
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 设置Gin模式 - 开发环境使用Debug模式以显示更多日志
	gin.SetMode(gin.DebugMode)

	// 输出启动信息到控制台
	fmt.Println("=== 问卷系统后端服务 ===")
	fmt.Println("正在初始化...")

	// 加载配置
	config := config.LoadConfig()

	// 初始化数据库
	db, err := database.InitDB(config)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 创建Gin路由
	router := gin.Default()

	// 应用中间件
	router.Use(CORSMiddleware())
	router.Use(LoggingMiddleware())

	// 创建处理器
	userHandler := handlers.NewUserHandler(db)
	questionnaireHandler := handlers.NewQuestionnaireHandler(db)
	adminHandler := handlers.NewAdminHandler(db)

	// 健康检查路由
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
		})
	})

	// 用户相关路由
	router.POST("/api/user/register", userHandler.Register)
	router.POST("/api/user/login", userHandler.Login)
	router.POST("/api/user/reset-password", userHandler.ResetPassword)

	// 问卷相关路由
	router.POST("/api/questionnaire/create", questionnaireHandler.CreateQuestionnaire)
	router.GET("/api/questionnaire/list", questionnaireHandler.GetQuestionnaires)
	router.GET("/api/questionnaire/detail", questionnaireHandler.GetQuestionnaireDetail)
	router.POST("/api/questionnaire/submit", questionnaireHandler.SubmitQuestionnaire)
	router.PUT("/api/questionnaire/update", questionnaireHandler.UpdateQuestionnaire)
	router.PUT("/api/questionnaire/update-status", questionnaireHandler.UpdateQuestionnaireStatus)
	router.DELETE("/api/questionnaire/delete", questionnaireHandler.DeleteQuestionnaire)
	router.GET("/api/questionnaire/results", questionnaireHandler.GetQuestionnaireResults)
	router.GET("/api/questionnaire/check-submission", questionnaireHandler.CheckSubmission)
	router.GET("/api/questionnaire/stats", questionnaireHandler.GetSystemStats)

	// 管理员路由组 - 使用管理员权限中间件
	adminGroup := router.Group("/api/admin")
	adminGroup.Use(middleware.AdminAuthMiddleware(db))
	{
		// 用户管理
		adminGroup.GET("/users", adminHandler.GetAllUsers)
		adminGroup.GET("/user/detail", adminHandler.GetUserDetail)
		adminGroup.PUT("/user/update", adminHandler.UpdateUser)
		adminGroup.DELETE("/user/delete", adminHandler.DeleteUser)

		// 问卷管理
		adminGroup.GET("/questionnaires", adminHandler.GetAllQuestionnaires)
		adminGroup.GET("/questionnaire/submissions", adminHandler.GetQuestionnaireSubmissions)

		// 系统统计
		adminGroup.GET("/statistics", adminHandler.GetSystemStatistics)
	}

	// 启动服务器
	serverAddr := ":" + config.Server.Port
	log.Printf("Gin服务器启动在 %s", serverAddr)
	log.Printf("按Ctrl+C停止服务器")

	// 使用标准的方式监听本地端口，避免Windows防火墙问题
	router.Run(serverAddr)
	// 明确指定监听所有网络接口，确保Windows防火墙可以正确识别
	//router.Run("0.0.0.0" + serverAddr) // 明确监听所有IP地址
}
