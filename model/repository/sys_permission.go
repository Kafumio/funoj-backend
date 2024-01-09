package repository

import "gorm.io/gorm"

type SysPermission struct {
	gorm.Model
	ParentID    uint   `gorm:"column:parent_id" json:"parentID"`      // 父权限id
	Code        string `gorm:"column:path" json:"path"`               // 权限代码
	Name        string `gorm:"column:name" json:"name"`               // 权限名
	Description string `gorm:"column:description" json:"description"` // 描述
	CateGory    string `gorm:"column:category" json:"cateGory"`       // 权限类别
	Method      string `gorm:"column:method" json:"method"`           // 权限请求方法
}

func (m *SysPermission) TableName() string {
	return "permission"
}
