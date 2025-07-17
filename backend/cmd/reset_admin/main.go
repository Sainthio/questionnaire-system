package main

import (
	"fmt"
	"log"
	"questionnaire-system/backend/config"
	"questionnaire-system/backend/database"
	"questionnaire-system/backend/models"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 连接数据库
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 查找管理员账号
	var admin models.User
	result := db.Where("username = ? AND is_admin = ?", "admin", true).First(&admin)
	if result.Error != nil {
		log.Printf("管理员账号不存在，正在创建...")

		// 创建管理员账号
		admin = models.User{
			Username: "admin",
			Email:    "admin@example.com",
			IsAdmin:  true,
		}
	}

	// 设置新密码并使用bcrypt加密
	plainPassword := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	admin.Password = string(hashedPassword)

	// 保存或更新管理员账号
	if admin.ID > 0 {
		// 更新现有账号
		result = db.Save(&admin)
		if result.Error != nil {
			log.Fatalf("更新管理员密码失败: %v", result.Error)
		}
		log.Printf("管理员密码已重置: admin/%s", plainPassword)
	} else {
		// 创建新账号
		result = db.Create(&admin)
		if result.Error != nil {
			log.Fatalf("创建管理员账号失败: %v", result.Error)
		}
		log.Printf("管理员账号已创建: admin/%s", plainPassword)
	}

	fmt.Println("========================================")
	fmt.Println("管理员账号信息:")
	fmt.Println("用户名: admin")
	fmt.Println("密码: admin123")
	fmt.Println("========================================")
}
