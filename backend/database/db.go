package database

import (
	"fmt"
	"log"
	"questionnaire-system/backend/config"
	"questionnaire-system/backend/models"

	"crypto/md5"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Database 封装数据库连接
type Database struct {
	*gorm.DB
}

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("数据库连接失败: %v", err)
		return nil, err
	}

	// 自动迁移数据表结构
	err = db.AutoMigrate(
		&models.User{},
		&models.Questionnaire{},
		&models.Question{},
		&models.Answer{},
		&models.Submission{},
	)
	if err != nil {
		log.Printf("数据库迁移失败: %v", err)
		return nil, err
	}

	log.Println("数据库连接成功")

	// 创建测试账号
	createTestAccounts(db)

	return &Database{db}, nil
}

// createTestAccounts 创建测试账号
func createTestAccounts(db *gorm.DB) {
	// 检查是否已有管理员账号
	var adminCount int64
	db.Model(&models.User{}).Where("is_admin = ?", true).Count(&adminCount)

	// 如果没有管理员账号，创建一个
	if adminCount == 0 {
		// 简单的密码加密（实际应用中应使用bcrypt等更安全的方式）
		adminPassword := fmt.Sprintf("%x", md5.Sum([]byte("admin123")))

		admin := models.User{
			Username: "admin",
			Password: adminPassword,
			Email:    "admin@example.com",
			IsAdmin:  true,
		}

		result := db.Create(&admin)
		if result.Error != nil {
			log.Printf("创建管理员账号失败: %v", result.Error)
		} else {
			log.Printf("已创建管理员账号: admin/admin123")
		}
	}

	// 检查是否已有普通用户账号
	var userCount int64
	db.Model(&models.User{}).Where("username = ?", "test").Count(&userCount)

	// 如果没有测试用户账号，创建一个
	if userCount == 0 {
		// 简单的密码加密
		testPassword := fmt.Sprintf("%x", md5.Sum([]byte("test123")))

		user := models.User{
			Username: "test",
			Password: testPassword,
			Email:    "test@example.com",
			IsAdmin:  false,
		}

		result := db.Create(&user)
		if result.Error != nil {
			log.Printf("创建测试用户账号失败: %v", result.Error)
		} else {
			log.Printf("已创建测试用户账号: test/test123")
		}
	}
}
