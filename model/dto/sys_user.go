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

type UserInfo struct {
	ID          uint     `json:"id"`
	Avatar      string   `json:"avatar"`
	LoginName   string   `json:"loginName"`
	UserName    string   `json:"userName"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	Roles       []uint   `json:"roles"`
	Permissions []string `json:"permissions"`
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

func NewUserInfo(user *repository.SysUser) *UserInfo {
	userInfo := &UserInfo{
		ID:        user.ID,
		Avatar:    user.Avatar,
		LoginName: user.LoginName,
		UserName:  user.UserName,
		Email:     user.Email,
		Phone:     user.Phone,
	}
	userInfo.Roles = make([]uint, len(user.Roles))
	for i := 0; i < len(user.Roles); i++ {
		userInfo.Roles[i] = user.Roles[i].ID
		for j := 0; j < len(user.Roles[i].Permissions); j++ {
			userInfo.Permissions = append(userInfo.Permissions, user.Roles[i].Permissions[j].Code)
		}
	}
	return userInfo
}
