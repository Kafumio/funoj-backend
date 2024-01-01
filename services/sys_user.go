package services

import (
	"errors"
	conf "funoj-backend/config"
	e "funoj-backend/consts/error"
	"funoj-backend/dao"
	"funoj-backend/db"
	"funoj-backend/model/dto"
	"funoj-backend/model/form/request"
	"funoj-backend/model/form/respnose"
	"funoj-backend/model/repository"
	"funoj-backend/utils"
	"gorm.io/gorm"
	"time"
)

type SysUserService interface {
}

type SysUserServiceImpl struct {
	config     *conf.AppConfig
	sysUserDao dao.SysUserDao
	sysRoleDao dao.SysRoleDao
}

func NewSysUserService(config *conf.AppConfig, userDao dao.SysUserDao, roleDao dao.SysRoleDao) SysUserService {
	return &SysUserServiceImpl{
		config:     config,
		sysUserDao: userDao,
		sysRoleDao: roleDao,
	}
}

func (s *SysUserServiceImpl) GetUserByID(userID uint) (*repository.SysUser, *e.Error) {
	user, err := s.sysUserDao.GetUserByID(db.Mysql, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, e.ErrUserNotExist
	}
	if err != nil {
		return nil, e.ErrMysql
	}
	return user, nil
}

func (s *SysUserServiceImpl) InsertSysUser(user *repository.SysUser) (uint, *e.Error) {
	// 设置默认用户名
	if user.UserName == "" {
		user.UserName = "funcoder"
	}
	// 随机登录名称
	if user.LoginName == "" {
		user.LoginName = user.LoginName + utils.GetUUID()
	}
	// 设置默认密码
	if user.Password == "" {
		user.Password = s.config.DefaultPassword
	}
	// 设置默认出生时间
	t := time.Time{}
	if user.BirthDay == t {
		user.BirthDay = time.Now()
	}
	// 设置默认性别
	if user.Gender != 1 && user.Gender != 2 {
		user.Gender = 1
	}
	p, err := utils.GetPwd(user.Password)
	if err != nil {
		return 0, e.ErrMysql
	}
	user.Password = string(p)
	err = s.sysUserDao.InsertUser(db.Mysql, user)
	if err != nil {
		return 0, e.ErrMysql
	}
	return user.ID, nil
}

func (s *SysUserServiceImpl) UpdateSysUser(user *repository.SysUser) *e.Error {
	user.UpdatedAt = time.Now()
	err := s.sysUserDao.UpdateUser(db.Mysql, user)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (s *SysUserServiceImpl) DeleteSysUser(userID uint) *e.Error {
	err := s.sysUserDao.DeleteUserByID(db.Mysql, userID)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (s *SysUserServiceImpl) GetSysUserList(pageQuery *request.PageQuery) (*respnose.PageInfo, *e.Error) {
	var pageInfo *respnose.PageInfo
	var userQuery *request.SysUserForList
	if pageQuery.Query != nil {
		userQuery = pageQuery.Query.(*request.SysUserForList)
	}
	err := db.Mysql.Transaction(func(tx *gorm.DB) error {
		userList, err := s.sysUserDao.GetUserList(tx, pageQuery)
		if err != nil {
			return err
		}
		userDtoList := make([]*dto.SysUserDto, len(userList))
		for i, user := range userList {
			user.Roles, err = s.sysUserDao.GetRolesByUserID(tx, user.ID)
			if err != nil {
				return err
			}
			userDtoList[i] = dto.NewSysUserDto(user)
		}
		var count int64
		count, err = s.sysUserDao.GetUserCount(tx, userQuery)
		if err != nil {
			return err
		}
		pageInfo = &respnose.PageInfo{
			Total: count,
			Size:  int64(len(userDtoList)),
			List:  userDtoList,
		}
		return nil
	})
	if err != nil {
		return nil, e.ErrMysql
	}
	return pageInfo, nil
}
