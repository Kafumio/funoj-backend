package services

import (
	e "funoj-backend/consts/error"
	"funoj-backend/dao"
)

type SysPermissionService interface {
}

type SysPermissionServiceImpl struct {
	sysPermissionDao dao.SysPermissionDao
}

func NewSysPermissionService(permissionDao dao.SysPermissionDao) SysPermissionService {
	return &SysPermissionServiceImpl{
		sysPermissionDao: permissionDao,
	}
}

func (s *SysPermissionServiceImpl) GetPermissionCount() (int64, *e.Error) {

}
