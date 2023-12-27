package repository

import "gorm.io/gorm"

type SysRole struct {
	gorm.Model
	Name        string `gorm:"column:name" json:"name"`               // 角色名称
	Description string `gorm:"column:description" json:"description"` // 描述
	// 角色权限组
	Permissions []*SysPermission `gorm:"many2many:role_permission" json:"permissions"`
}

func (s *SysRole) TableName() string {
	return "role"
}
