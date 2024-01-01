package services

import (
	"errors"
	conf "funoj-backend/config"
	e "funoj-backend/consts/error"
	"funoj-backend/dao"
	"funoj-backend/db"
	"funoj-backend/model/dto"
	"funoj-backend/model/form/request"
	"funoj-backend/model/form/response"
	"funoj-backend/model/repository"
	"funoj-backend/utils"
	"gorm.io/gorm"
	"time"
)

type SysUserService interface {
	// GetUserByID 根据用户id获取用户信息
	GetUserByID(userID uint) (*repository.SysUser, *e.Error)
	// InsertUser 添加用户
	InsertUser(user *repository.SysUser) (uint, *e.Error)
	// UpdateUser 更新用户，但是不更新密码
	UpdateUser(user *repository.SysUser) *e.Error
	// DeleteUser 删除用户
	DeleteUser(userID uint) *e.Error
	// GetUserList 获取用户列表
	GetUserList(pageQuery *request.PageQuery) (*response.PageInfo, *e.Error)
	// UpdateUserRoles 更新角色roleIDs
	UpdateUserRoles(userID uint, roleIDs []uint) *e.Error
	// GetRoleIDsByUserID 通过用户id获取所有角色id
	GetRoleIDsByUserID(userID uint) ([]uint, *e.Error)
	// GetAllSimpleRole 简单展示用户可归属的所有角色
	GetAllSimpleRole() ([]*dto.SimpleRoleDto, *e.Error)
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
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrUserNotExist
		} else {
			return nil, e.ErrMysql
		}
	}
	return user, nil
}

func (s *SysUserServiceImpl) InsertUser(user *repository.SysUser) (uint, *e.Error) {
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

func (s *SysUserServiceImpl) UpdateUser(user *repository.SysUser) *e.Error {
	user.UpdatedAt = time.Now()
	err := s.sysUserDao.UpdateUser(db.Mysql, user)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (s *SysUserServiceImpl) DeleteUser(userID uint) *e.Error {
	err := s.sysUserDao.DeleteUserByID(db.Mysql, userID)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (s *SysUserServiceImpl) GetUserList(pageQuery *request.PageQuery) (*response.PageInfo, *e.Error) {
	var pageInfo *response.PageInfo
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
		pageInfo = &response.PageInfo{
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

func (s *SysUserServiceImpl) UpdateUserRoles(userID uint, roleIDs []uint) *e.Error {
	tx := db.Mysql.Begin()
	err := s.sysUserDao.DeleteUserRoleByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return e.ErrMysql
	}
	err = s.sysUserDao.InsertRolesToUser(tx, userID, roleIDs)
	if err != nil {
		tx.Rollback()
		return e.ErrMysql
	}
	tx.Commit()
	return nil
}

func (s *SysUserServiceImpl) GetRoleIDsByUserID(userID uint) ([]uint, *e.Error) {
	roleIDs, err := s.sysUserDao.GetRoleIDsByUserID(db.Mysql, userID)
	if err != nil {
		return nil, e.ErrMysql
	}
	return roleIDs, nil
}

func (s *SysUserServiceImpl) GetAllSimpleRole() ([]*dto.SimpleRoleDto, *e.Error) {
	roles, err := s.sysRoleDao.GetAllSimpleRoleList(db.Mysql)
	if err != nil {
		return nil, e.ErrMysql
	}
	simpleRoles := make([]*dto.SimpleRoleDto, len(roles))
	for i, role := range roles {
		simpleRoles[i] = dto.NewSimpleRoleDto(role)
	}
	return simpleRoles, nil
}
