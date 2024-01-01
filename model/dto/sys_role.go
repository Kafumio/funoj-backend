package dto

import (
	"funoj-backend/model/repository"
	"funoj-backend/utils"
)

// SysRoleDto 角色dto
type SysRoleDto struct {
	ID          uint       `json:"id"`
	UpdatedAt   utils.Time `json:"updatedAt"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
}

func NewSysRoleDto(role *repository.SysRole) *SysRoleDto {
	response := &SysRoleDto{
		ID:          role.ID,
		UpdatedAt:   utils.Time(role.UpdatedAt),
		Name:        role.Name,
		Description: role.Description,
	}
	return response
}

// SimpleRoleDto 简单的角色，用于获取只有id和名称的角色列表
type SimpleRoleDto struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewSimpleRoleDto(role *repository.SysRole) *SimpleRoleDto {
	response := &SimpleRoleDto{
		ID:   role.ID,
		Name: role.Name,
	}
	return response
}
