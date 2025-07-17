package middleware

import (
	"log"
	"questionnaire-system/backend/database"
	"questionnaire-system/backend/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware 管理员权限验证中间件
func AdminAuthMiddleware(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("执行管理员权限验证中间件")

		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("缺少Authorization头")
			c.JSON(401, gin.H{
				"success": false,
				"message": "未授权访问",
			})
			c.Abort()
			return
		}

		// 解析token
		// 格式: Bearer token_username_timestamp
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("无效的Authorization头格式")
			c.JSON(401, gin.H{
				"success": false,
				"message": "无效的授权格式",
			})
			c.Abort()
			return
		}

		token := parts[1]
		// 从token中提取用户名 (格式: token_username_timestamp)
		tokenParts := strings.Split(token, "_")
		if len(tokenParts) < 2 {
			log.Println("无效的token格式")
			c.JSON(401, gin.H{
				"success": false,
				"message": "无效的令牌",
			})
			c.Abort()
			return
		}

		username := tokenParts[1]
		log.Printf("从token中提取的用户名: %s", username)

		// 查询用户
		var user models.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			log.Printf("用户不存在: %s", username)
			c.JSON(401, gin.H{
				"success": false,
				"message": "无效的用户",
			})
			c.Abort()
			return
		}

		// 验证管理员权限
		if !user.IsAdmin {
			log.Printf("用户不是管理员: %s", username)
			c.JSON(403, gin.H{
				"success": false,
				"message": "需要管理员权限",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", user.ID)
		c.Set("username", user.Username)
		c.Set("is_admin", user.IsAdmin)

		log.Printf("管理员验证通过: %s", username)
		c.Next()
	}
}
