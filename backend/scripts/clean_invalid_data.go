package main

import (
	"fmt"
	"log"
	"os"
	"questionnaire-system/backend/config"
	"questionnaire-system/backend/database"
	"questionnaire-system/backend/models"
	"strings"
)

func main() {
	for {
		// 设置日志输出
		logFile, err := os.OpenFile("clean_data.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("无法打开日志文件: %v\n", err)
			return
		}
		defer logFile.Close()

		// 同时输出到控制台和文件
		log.SetOutput(logFile)
		log.SetFlags(log.LstdFlags | log.Lshortfile)

		log.Println("开始清理无效数据")

		// 加载配置
		config := config.LoadConfig()

		// 初始化数据库连接
		db, err := database.InitDB(config)
		if err != nil {
			log.Fatalf("数据库连接失败: %v", err)
		}

		// 1. 查找所有user_id为0的答案记录
		var invalidAnswers []models.Answer
		if err := db.Where("user_id = 0").Find(&invalidAnswers).Error; err != nil {
			log.Fatalf("查询无效答案失败: %v", err)
		}

		log.Printf("找到 %d 条user_id为0的答案记录", len(invalidAnswers))

		// 2. 查找所有user_id为0的提交记录
		var invalidSubmissions []models.Submission
		if err := db.Where("user_id = 0").Find(&invalidSubmissions).Error; err != nil {
			log.Fatalf("查询无效提交记录失败: %v", err)
		}

		log.Printf("找到 %d 条user_id为0的提交记录", len(invalidSubmissions))

		// 显示菜单
		fmt.Printf("找到 %d 条user_id为0的答案记录和 %d 条user_id为0的提交记录\n", len(invalidAnswers), len(invalidSubmissions))
		fmt.Println("请选择操作:")
		fmt.Println("1. 查看无效答案详情")
		fmt.Println("2. 查看无效提交记录详情")
		fmt.Println("3. 删除所有无效数据")
		fmt.Println("4. 退出")
		fmt.Print("请输入选项 (1-4): ")

		var option string
		fmt.Scanln(&option)

		switch strings.TrimSpace(option) {
		case "1":
			// 显示无效答案详情
			fmt.Println("\n无效答案详情:")
			fmt.Println("ID\t问题ID\t用户ID\t答案内容\t创建时间")
			fmt.Println("--------------------------------------------------")
			for _, answer := range invalidAnswers {
				fmt.Printf("%d\t%d\t%d\t%s\t%v\n",
					answer.ID, answer.QuestionID, answer.UserID,
					answer.Content, answer.CreatedAt)
			}
			return

		case "2":
			// 显示无效提交记录详情
			fmt.Println("\n无效提交记录详情:")
			fmt.Println("ID\t问卷ID\t提交时间\tIP地址")
			fmt.Println("--------------------------------------------------")
			for _, submission := range invalidSubmissions {
				fmt.Printf("%d\t%d\t%v\t%s\n",
					submission.ID, submission.QuestionnaireID,
					submission.SubmittedAt, submission.IPAddress)
			}
			return

		case "3":
			// 删除所有无效数据
			fmt.Print("确认删除所有无效数据? (y/n): ")
			var confirm string
			fmt.Scanln(&confirm)

			if confirm != "y" && confirm != "Y" {
				log.Println("用户取消操作，退出")
				fmt.Println("操作已取消")
				return
			}

			// 开始事务
			tx := db.Begin()

			// 删除无效答案
			if len(invalidAnswers) > 0 {
				if err := tx.Delete(&invalidAnswers).Error; err != nil {
					tx.Rollback()
					log.Fatalf("删除无效答案失败: %v", err)
				}
				log.Printf("已删除 %d 条无效答案", len(invalidAnswers))
			}

			// 删除无效提交记录
			if len(invalidSubmissions) > 0 {
				if err := tx.Delete(&invalidSubmissions).Error; err != nil {
					tx.Rollback()
					log.Fatalf("删除无效提交记录失败: %v", err)
				}
				log.Printf("已删除 %d 条无效提交记录", len(invalidSubmissions))
			}

			// 提交事务
			if err := tx.Commit().Error; err != nil {
				log.Fatalf("提交事务失败: %v", err)
			}

			log.Println("数据清理完成")
			fmt.Println("数据清理完成，详细信息请查看clean_data.log")

		case "4":
			fmt.Println("退出程序")
			return

		default:
			fmt.Println("无效选项，退出程序")
			return
		}
	}
}
