package repository

import (
	"gorm.io/gorm"
	"time"
)

// 用户的提交及运行结果
type Submission struct {
	gorm.Model
	// 用户id
	UserID uint `gorm:"column:user_id" json:"userID"`
	// 题目id
	ProblemID uint `gorm:"column:problem_id" json:"problemID"`
	// 使用的编程语言
	Language string `gorm:"column:language" json:"language"`
	// 用户代码
	Code string `gorm:"column:code" json:"code"`
	// 状态
	Status int `gorm:"column:status" json:"status"`
	// 异常信息
	ErrorMessage string `gorm:"column:error_message" json:"errorMessage"`
	// 用例名称
	CaseName string `gorm:"column:case_name" json:"caseName"`
	// 用例数据
	CaseData string `gorm:"column:case_data" json:"caseData"`
	// 期望输出
	ExpectedOutput string `gorm:"column:expected_output" json:"expectedOutput"`
	// 用户输出
	UserOutput string        `gorm:"user_output" json:"userOutput"`
	TimeUsed   time.Duration // 判题使用时间
	MemoryUsed int64         // 内存使用量（以字节为单位）
}

func (s *Submission) TableName() string {
	return "submission"
}
