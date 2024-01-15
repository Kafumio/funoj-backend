package services

import (
	conf "funoj-backend/config"
	e "funoj-backend/consts/error"
	"funoj-backend/dao"
	"funoj-backend/db"
	"funoj-backend/file_store"
	"funoj-backend/model/dto"
	"funoj-backend/model/form/response"
	"funoj-backend/model/repository"
	"funoj-backend/utils"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"path"
	"time"
)

const (
	// UserAvatarPath cos中，用户图片存储的位置
	UserAvatarPath = "/avatar/user"
)

type AccountService interface {
	// UploadAvatar 上传头像
	UploadAvatar(file *multipart.FileHeader) (string, *e.Error)
	// ReadAvatar 读取头像
	ReadAvatar(ctx *gin.Context, avatarName string)
	// GetAccountInfo 获取账号信息
	GetAccountInfo(ctx *gin.Context) (*dto.AccountInfo, *e.Error)
	// UpdateAccountInfo 更新账号信息
	UpdateAccountInfo(ctx *gin.Context, user *repository.SysUser) *e.Error
	// ChangePassword 修改密码
	ChangePassword(ctx *gin.Context, oldPassword, newPassword string) *e.Error
	// ResetPassword 重置密码
	ResetPassword(ctx *gin.Context) *e.Error
}

type AccountServiceImpl struct {
	config     *conf.AppConfig
	sysUserDao dao.SysUserDao
}

func NewAccountService(config *conf.AppConfig, userDao dao.SysUserDao) AccountService {
	return &AccountServiceImpl{
		config:     config,
		sysUserDao: userDao,
	}
}

func (svc *AccountServiceImpl) UploadAvatar(file *multipart.FileHeader) (string, *e.Error) {
	cos := file_store.NewImageCOS(svc.config.COSConfig)
	fileName := file.Filename
	fileName = utils.GetUUID() + "." + path.Base(fileName)
	file2, err := file.Open()
	if err != nil {
		return "", e.ErrBadRequest
	}
	err = cos.SaveFile(UserAvatarPath+"/"+fileName, file2)
	if err != nil {
		log.Println(err)
		return "", e.ErrServer
	}
	return svc.config.ProUrl + UserAvatarPath + "/" + fileName, nil
}

func (svc *AccountServiceImpl) ReadAvatar(ctx *gin.Context, avatarName string) {
	result := response.NewResult(ctx)
	cos := file_store.NewImageCOS(svc.config.COSConfig)
	bytes, err := cos.ReadFile(UserAvatarPath + "/" + avatarName)
	if err != nil {
		result.Error(e.ErrServer)
		return
	}
	_, _ = ctx.Writer.Write(bytes)
}

func (svc *AccountServiceImpl) GetAccountInfo(ctx *gin.Context) (*dto.AccountInfo, *e.Error) {
	user := ctx.Keys["user"].(*dto.UserInfo)
	u, err := svc.sysUserDao.GetUserByID(db.Mysql, user.ID)
	if err != nil {
		return nil, e.ErrMysql
	}
	return dto.NewAccountInfo(u), nil
}

func (svc *AccountServiceImpl) UpdateAccountInfo(ctx *gin.Context, user *repository.SysUser) *e.Error {
	userInfo := ctx.Keys["user"].(*dto.UserInfo)
	user.ID = userInfo.ID
	user.UpdatedAt = time.Now()
	// 不能更新账号名称和密码
	user.LoginName = ""
	user.Password = ""
	err := svc.sysUserDao.UpdateUser(db.Mysql, user)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (svc *AccountServiceImpl) ChangePassword(ctx *gin.Context, oldPassword, newPassword string) *e.Error {
	userInfo := ctx.Keys["user"].(*dto.UserInfo)
	// 检验用户名
	user, err := svc.sysUserDao.GetUserByID(db.Mysql, userInfo.ID)
	if err != nil {
		log.Println(err)
		return e.ErrMysql
	}
	if user == nil || user.LoginName == "" {
		return e.ErrUserNotExist
	}
	// 检验旧密码
	if !utils.ComparePwd(oldPassword, user.Password) {
		return e.ErrUserNameOrPasswordWrong
	}
	password, getPwdErr := utils.GetPwd(newPassword)
	if getPwdErr != nil {
		log.Println(getPwdErr)
		return e.ErrPasswordEncodeFailed
	}
	user.Password = string(password)
	user.UpdatedAt = time.Now()
	err = svc.sysUserDao.UpdateUser(db.Mysql, user)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (svc *AccountServiceImpl) ResetPassword(ctx *gin.Context) *e.Error {
	userInfo := ctx.Keys["user"].(*dto.UserInfo)
	passwordRoot := utils.GetRandomPassword(11)
	password, err := utils.GetPwd(passwordRoot)
	if err != nil {
		return e.ErrUserUnknownError
	}
	// 更新密码
	tx := db.Mysql.Begin()
	user := &repository.SysUser{}
	user.ID = userInfo.ID
	user.Password = string(password)
	err = svc.sysUserDao.UpdateUser(tx, user)
	if err != nil {
		tx.Rollback()
		return e.ErrMysql
	}
	// 发送密码
	message := utils.SysEmailMessage{
		To:      []string{user.Email},
		Subject: "fancode-重置密码",
		Body:    "新密码：" + passwordRoot,
	}
	err = utils.SendSysEmail(svc.config.EmailConfig, message)
	if err != nil {
		tx.Rollback()
		return e.ErrUserUnknownError
	}
	return nil
}
