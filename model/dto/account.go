package dto

import (
	"funoj-backend/model/repository"
	"time"
)

// AccountInfo
// 和userInfo类似，但是比userInfo的数据多一些
type AccountInfo struct {
	Avatar       string `json:"avatar"`
	LoginName    string `json:"loginName"`
	UserName     string `json:"userName"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Introduction string `json:"introduction"`
	Gender       int    `json:"gender"`
	BirthDay     string `json:"birthDay"`
	CodingAge    int    `json:"codingAge"`
}

func NewAccountInfo(user *repository.SysUser) *AccountInfo {
	return &AccountInfo{
		Avatar:       user.Avatar,
		LoginName:    user.LoginName,
		UserName:     user.UserName,
		Email:        user.Email,
		Phone:        user.Phone,
		Introduction: user.Introduction,
		BirthDay:     user.BirthDay.Format("2006-01-02"),
		Gender:       user.Gender,
		CodingAge:    time.Now().Year() - user.CreatedAt.Year(),
	}
}
