package repository

import "gorm.io/gorm"

// ProblemAttempt 用户在一道题目中的做题情况
type ProblemAttempt struct {
	gorm.Model
	ProblemID       uint `gorm:"column:problem_id" json:"problemID"`
	UserID          uint `gorm:"column:user_id" json:"userID"`
	SubmissionCount int  `gorm:"column:submission_count" json:"submissionCount"`
	SuccessCount    int  `gorm:"column:success_count" json:"successCount"`
	ErrCount        int  `gorm:"column:err_count" json:"errCount"`
	// 最近一次的代码
	Code     string `gorm:"column:code" json:"code"`
	Language string `gorm:"column:language" json:"language"`
	// 0 未开始，1进行中 2 提交成功
	Status int `gorm:"column:status" json:"status"`
}

func (m *ProblemAttempt) TableName() string {
	return "problem_attempt"
}
