package dto

import "funoj-backend/utils"

// SysUserDtoForList 获取用户列表
type SysUserDtoForList struct {
	ID        uint       `json:"id"`
	LoginName string     `json:"loginName"`
	UserName  string     `json:"username"`
	Gender    int        `json:"gender"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	UpdateAt  utils.Time `json:"updatedAt"`
	Roles     []string   `json:"roles"`
}
