package repository

import "gorm.io/gorm"

// 题目标签
type ProblemTag struct {
	gorm.Model
	Name     string     `gorm:"column:name" json:"name"`
	Problems []*Problem `gorm:"many2many:problem_menu_association" json:"problems"`
}

func (m *ProblemTag) TableName() string {
	return "problem_tag"
}
