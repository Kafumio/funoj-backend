package services

import (
	"errors"
	"funoj-backend/config"
	"funoj-backend/consts"
	e "funoj-backend/consts/error"
	"funoj-backend/dao"
	"funoj-backend/db"
	"funoj-backend/model/dto"
	"funoj-backend/model/repository"
	"funoj-backend/utils"
	"github.com/Chain-Zhang/pinyin"
	"gorm.io/gorm"
	"log"
	"time"
)

const (
	RegisterEmailProKey = "emailcode-register-"
	LoginEmailProKey    = "emailcode-login-"
)

type AuthService interface {
	// LoginByPassword 密码登录 account可能是邮箱可能是用户id
	LoginByPassword(account string, password string) (string, *e.Error)
	// LoginByEmail 邮箱验证登录
	LoginByEmail(email string, code string) (string, *e.Error)
	// SendAuthCode 获取邮件的验证码
	SendAuthCode(email string, kind string) (string, *e.Error)
	// UserRegister 用户注册
	UserRegister(user *repository.SysUser, code string) *e.Error
}

type AuthServiceImpl struct {
	config     *config.AppConfig
	sysUserDao dao.SysUserDao
	sysRoleDao dao.SysRoleDao
}

func NewAuthService(config *config.AppConfig, userDao dao.SysUserDao, roleDao dao.SysRoleDao) AuthService {
	return &AuthServiceImpl{
		config:     config,
		sysUserDao: userDao,
		sysRoleDao: roleDao,
	}
}

func (svc *AuthServiceImpl) LoginByPassword(account string, password string) (string, *e.Error) {
	var user *repository.SysUser
	var userErr error
	if utils.VerifyEmailFormat(account) {
		user, userErr = svc.sysUserDao.GetUserByEmail(db.Mysql, account)
	} else {
		user, userErr = svc.sysUserDao.GetUserByLoginName(db.Mysql, account)
	}
	if user == nil || errors.Is(userErr, gorm.ErrRecordNotFound) {
		return "", e.ErrUserNotExist
	}
	if userErr != nil {
		log.Println(userErr)
		return "", e.ErrUserUnknownError
	}
	// 比较密码
	if !utils.ComparePwd(user.Password, password) {
		return "", e.ErrUserNameOrPasswordWrong
	}
	// 读取权限
	for i := 0; i < len(user.Roles); i++ {
		var err error
		user.Roles[i].Permissions, err = svc.sysRoleDao.GetPermissionsByRoleID(db.Mysql, user.Roles[i].ID)
		if err != nil {
			return "", e.ErrUserUnknownError
		}
	}
	userInfo := dto.NewUserInfo(user)
	token, err := utils.GenerateToken(utils.Claims{
		ID:          userInfo.ID,
		Avatar:      userInfo.Avatar,
		UserName:    userInfo.UserName,
		LoginName:   userInfo.LoginName,
		Phone:       userInfo.Phone,
		Email:       userInfo.Email,
		Roles:       userInfo.Roles,
		Permissions: userInfo.Permissions,
	})
	if err != nil {
		log.Println(err)
		return "", e.ErrUserUnknownError
	}
	return token, nil
}

func (svc *AuthServiceImpl) LoginByEmail(email string, code string) (string, *e.Error) {
	if !utils.VerifyEmailFormat(email) {
		return "", e.ErrUserUnknownError
	}
	// 获取用户
	user, err := svc.sysUserDao.GetUserByEmail(db.Mysql, email)
	if err != nil {
		log.Println(err)
		return "", e.ErrUserUnknownError
	}
	// 检测验证码
	key := LoginEmailProKey + email
	result, err2 := db.Redis.Get(key).Result()
	if err2 != nil {
		log.Println(err)
		return "", e.ErrUserUnknownError
	}
	if result != code {
		return "", e.ErrLoginCodeWrong
	}
	// 读取权限
	for i := 0; i < len(user.Roles); i++ {
		var err error
		user.Roles[i].Permissions, err = svc.sysRoleDao.GetPermissionsByRoleID(db.Mysql, user.Roles[i].ID)
		if err != nil {
			return "", e.ErrUserUnknownError
		}
	}
	userInfo := dto.NewUserInfo(user)
	token, err := utils.GenerateToken(utils.Claims{
		ID:          userInfo.ID,
		Avatar:      userInfo.Avatar,
		UserName:    userInfo.UserName,
		LoginName:   userInfo.LoginName,
		Phone:       userInfo.Phone,
		Email:       userInfo.Email,
		Roles:       userInfo.Roles,
		Permissions: userInfo.Permissions,
	})
	if err != nil {
		log.Println(err)
		return "", e.ErrUserUnknownError
	}
	return token, nil
}

func (svc *AuthServiceImpl) SendAuthCode(email string, kind string) (string, *e.Error) {
	if kind == "register" {
		f, err := svc.sysUserDao.CheckEmail(db.Mysql, email)
		if err != nil {
			return "", e.ErrUserUnknownError
		}
		if f {
			return "", e.ErrUserEmailIsExist
		}
	}

	var subject string
	if kind == "register" {
		subject = "funoj 注册验证码"
	} else if kind == "login" {
		subject = "funoj 登录验证码"
	}
	// 发送code
	code := utils.GetCheckNumber(6)
	message := utils.SysEmailMessage{
		To:      []string{email},
		Subject: subject,
		Body:    "验证码：" + code,
	}
	err := utils.SendSysEmail(svc.config.EmailConfig, message)
	if err != nil {
		log.Println(err)
		return "", e.ErrUserUnknownError
	}
	// 存储到redis
	var key string
	if kind == "register" {
		key = RegisterEmailProKey + email
	} else {
		key = LoginEmailProKey + email
	}
	_, err2 := db.Redis.Set(key, code, 10*time.Minute).Result()
	if err2 != nil {
		log.Println(err2)
		return "", e.ErrUserUnknownError
	}
	return code, nil
}

func (svc *AuthServiceImpl) UserRegister(user *repository.SysUser, code string) *e.Error {
	// 检测是否已注册过
	f, _ := svc.sysUserDao.CheckEmail(db.Mysql, user.Email)
	if f {
		return e.ErrUserEmailIsExist
	}
	// 检测code
	result := db.Redis.Get(RegisterEmailProKey + user.Email)
	if result.Err() != nil {
		return e.ErrUserUnknownError
	}
	if result.Val() != code {
		return e.ErrRoleUnknownError
	}
	// 设置用户名
	if user.UserName == "" {
		user.UserName = "fun-coder"
		return nil
	}
	// 生成用户名称，唯一
	loginName, err := pinyin.New(user.UserName).Split("").Convert()
	if err != nil {
		return e.ErrUserUnknownError
	}
	loginName = loginName + utils.GetCheckNumber(3)
	for i := 0; i < 5; i++ {
		b, err := svc.sysUserDao.CheckLoginName(db.Mysql, user.LoginName)
		if err != nil {
			log.Println(err)
			return e.ErrUserUnknownError
		}
		if b {
			loginName = loginName + utils.GetCheckNumber(1)
		} else {
			break
		}
	}
	user.LoginName = loginName
	if len(user.Password) < 6 {
		return e.ErrUserPasswordNotEnoughAccuracy
	}
	//进行注册操作
	newPassword, err := utils.GetPwd(user.Password)
	if err != nil {
		return e.ErrPasswordEncodeFailed
	}
	user.Password = string(newPassword)

	err = db.Mysql.Transaction(func(tx *gorm.DB) error {
		err2 := svc.sysUserDao.InsertUser(tx, user)
		if err2 != nil {
			return err
		}
		err2 = svc.sysUserDao.InsertRolesToUser(tx, user.ID, []uint{consts.UserID})
		return err2
	})
	if err != nil {
		return e.ErrMysql
	}
	return nil
}
