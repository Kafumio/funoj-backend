package services

import (
	conf "funoj-backend/config"
	"funoj-backend/dao"
	"funoj-backend/model/repository"
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

func (s *SysUserServiceImpl) GetUserByID(userID uint) (*repository.SysUser, error) {

}
