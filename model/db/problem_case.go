package db

import "gorm.io/gorm"

// 题目用例
type ProblemCase struct {
	gorm.Model
	ProblemID uint   `gorm:"column:problem_id" json:"problemID"`
	CaseName  string `gorm:"column:case_name" json:"caseName"`
	Input     string `gorm:"column:input" json:"input"`
	Output    string `gorm:"column:output" json:"output"`
}

func (p *ProblemCase) TableName() string {
	return "problem_case"
}
