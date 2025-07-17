package models

import "time"

// Questionnaire 问卷模型
type Questionnaire struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedBy   uint      `json:"created_by" gorm:"not null"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	IsPublished bool      `json:"is_published" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Question 问题模型
type Question struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	QuestionnaireID uint      `json:"questionnaire_id" gorm:"not null"`
	Title           string    `json:"title" gorm:"size:255;not null"`
	Type            string    `json:"type" gorm:"size:50;not null"` // 单选、多选、填空、评分等
	Required        bool      `json:"required" gorm:"default:false"`
	Options         string    `json:"options" gorm:"type:text"` // JSON格式存储选项
	Sort            int       `json:"sort" gorm:"default:0"`    // 排序
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Answer 答案模型
type Answer struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	QuestionID uint      `json:"question_id" gorm:"not null"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	Content    string    `json:"content" gorm:"type:text"` // 答案内容
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// Submission 提交记录
type Submission struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	QuestionnaireID uint      `json:"questionnaire_id" gorm:"not null"`
	UserID          uint      `json:"user_id" gorm:"not null"`
	SubmittedAt     time.Time `json:"submitted_at" gorm:"autoCreateTime"`
	IPAddress       string    `json:"ip_address" gorm:"size:50"`
}

// User 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"size:50;not null;uniqueIndex"`
	Password  string    `json:"password,omitempty" gorm:"size:255;not null"` // 在API响应中省略密码
	Email     string    `json:"email" gorm:"size:100;uniqueIndex"`
	Phone     string    `json:"phone" gorm:"size:20"`
	IsAdmin   bool      `json:"is_admin" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
