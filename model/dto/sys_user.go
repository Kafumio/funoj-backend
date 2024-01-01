package dto

import (
	"funoj-backend/model/repository"
	"funoj-backend/utils"
)

// SysUserDto 用户dto
type SysUserDto struct {
	ID        uint       `json:"id"`
	LoginName string     `json:"loginName"`
	UserName  string     `json:"username"`
	Gender    int        `json:"gender"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	UpdateAt  utils.Time `json:"updatedAt"`
	Roles     []string   `json:"roles"`
}

func NewSysUserDto(user *repository.SysUser) *SysUserDto {
	response := &SysUserDto{
		ID:        user.ID,
		LoginName: user.LoginName,
		UserName:  user.UserName,
		Email:     user.Email,
		Phone:     user.Phone,
		UpdateAt:  utils.Time(user.UpdatedAt),
	}
	if user.Roles != nil {
		response.Roles = make([]string, len(user.Roles))
		for i, role := range user.Roles {
			response.Roles[i] = role.Name
		}
	}
	return response
}
