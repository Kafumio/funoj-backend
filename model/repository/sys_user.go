package repository

import (
	"gorm.io/gorm"
	"time"
)

type SysUser struct {
	gorm.Model
	Avatar       string `gorm:"column:avatar" json:"avatar"`
	UserName     string `gorm:"column:user_name" json:"username"`
	LoginName    string `gorm:"column:login_name" json:"loginName"`
	Password     string `gorm:"column:password" json:""`
	Email        string `gorm:"column:email" json:"email"`
	Phone        string `gorm:"column:phone" json:"phone"`
	Introduction string `gorm:"column:introduction" json:"introduction"`
	// 1表示男 2表示女
	Gender   int       `gorm:"column:gender" json:"gender"`
	BirthDay time.Time `gorm:"column:birth_day" json:"birthDay"`
	// 所属角色
	Roles []*SysRole `gorm:"many2many:user_role;" json:"roles"`
}

func (m *SysUser) TableName() string {
	return "user"
}
