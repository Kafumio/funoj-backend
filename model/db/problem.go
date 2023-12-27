package db

import "gorm.io/gorm"

type Problem struct {
	gorm.Model
	CreatorID   uint   `gorm:"column:creator_id" json:"creatorID"`
	Number      string `gorm:"column:number;type:varchar(255);unique_index:idx_number" json:"number"`
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description;type:text" json:"description"`
	Title       string `gorm:"column:title" json:"title"`
	Difficulty  int    `gorm:"column:difficulty" json:"difficulty"`
	// 0空值，1启用，-1停用
	Enable int `gorm:"column:enable" json:"enable"`
	// 支持的语言用,分割
	Languages string `gorm:"column:languages" json:"languages"`
	// 所属标签
	// Tags []*ProblemTag `gorm:"many2many:problem_tag_association" json:"tags"`
}

func (p *Problem) TableName() string {
	return "problem"
}
