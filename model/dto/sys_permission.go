package dto

import (
	"funoj-backend/model/repository"
	"funoj-backend/utils"
)

// 权限树dto
type SysPermissionTreeDto struct {
	ID          uint                    `json:"id"`
	ParentID    uint                    `json:"parentApiID"` // 父权限id
	Method      string                  `json:"method"`      // 请求方法
	Name        string                  `json:"name"`        // 权限名
	CateGory    string                  `json:"cateGory"`    // 权限类别
	Description string                  `json:"description"` // 描述
	UpdatedAt   utils.Time              `json:"updatedAt"`   // 更新时间
	Children    []*SysPermissionTreeDto `json:"children"`    // 子权限组
}

func NewSysApiTreeDto(sysPermission *repository.SysPermission) *SysPermissionTreeDto {
	return &SysPermissionTreeDto{
		ID:          sysPermission.ID,
		ParentID:    sysPermission.ParentID,
		Method:      sysPermission.Method,
		Name:        sysPermission.Name,
		CateGory:    sysPermission.CateGory,
		Description: sysPermission.Description,
		UpdatedAt:   utils.Time(sysPermission.UpdatedAt),
	}
}
