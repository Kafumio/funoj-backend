package repository

import "gorm.io/gorm"

// 题单
type ProblemMenu struct {
	gorm.Model
	Name        string     `gorm:"column:name" json:"name"`
	Icon        string     `gorm:"column:icon" json:"icon"`
	Description string     `gorm:"column:description" json:"description"`
	CreatorID   uint       `gorm:"column:creator_id" json:"creatorID"`
	Problems    []*Problem `gorm:"many2many:problem_menu_association" json:"problems"`
}

func (p *ProblemMenu) TableName() string {
	return "problem_menu"
}
