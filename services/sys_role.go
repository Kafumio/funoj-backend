package services

import (
	"errors"
	e "funoj-backend/consts/error"
	"funoj-backend/dao"
	"funoj-backend/db"
	"funoj-backend/model/dto"
	"funoj-backend/model/form/request"
	"funoj-backend/model/form/response"
	"funoj-backend/model/repository"
	"gorm.io/gorm"
	"time"
)

type SysRoleService interface {
	// GetRoleByID 根角色户id获取角色信息
	GetRoleByID(roleID uint) (*repository.SysRole, *e.Error)
	// InsertRole 添加角色
	InsertRole(sysSysRole *repository.SysRole) (uint, *e.Error)
	// UpdateRole 更新角色
	UpdateRole(SysRole *repository.SysRole) *e.Error
	// DeleteRole 删除角色
	DeleteRole(id uint) *e.Error
	// GetRoleList 获取角色列表
	GetRoleList(pageQuery *request.PageQuery) (*response.PageInfo, *e.Error)
	// UpdateRolePermissions 更新角色权限
	UpdateRolePermissions(roleID uint, permissionIDs []uint) *e.Error
	// GetPermissionIDsByRoleID 通过角色id获取该角色拥有的权限ID
	GetPermissionIDsByRoleID(roleID uint) ([]uint, *e.Error)
	// GetPermissionsByRoleID 通过角色id获取该角色的所有权限
	GetPermissionsByRoleID(roleID uint) ([]*repository.SysPermission, *e.Error)
}

type SysRoleServiceImpl struct {
	sysRoleDao dao.SysRoleDao
}

func NewSysRoleService(roleDao dao.SysRoleDao) SysRoleService {
	return &SysRoleServiceImpl{
		sysRoleDao: roleDao,
	}
}

func (svc *SysRoleServiceImpl) GetRoleByID(roleID uint) (*repository.SysRole, *e.Error) {
	role, err := svc.sysRoleDao.GetRoleByID(db.Mysql, roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrRoleNotExist
		} else {
			return nil, e.ErrMysql
		}
	}
	return role, nil
}

func (svc *SysRoleServiceImpl) InsertRole(role *repository.SysRole) (uint, *e.Error) {
	err := svc.sysRoleDao.InsertRole(db.Mysql, role)
	if err != nil {
		return 0, e.ErrMysql
	}
	return role.ID, nil
}

func (svc *SysRoleServiceImpl) UpdateRole(sysRole *repository.SysRole) *e.Error {
	sysRole.UpdatedAt = time.Now()
	err := svc.sysRoleDao.UpdateRole(db.Mysql, sysRole)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (svc *SysRoleServiceImpl) DeleteRole(id uint) *e.Error {
	// 删除删除角色
	err := svc.sysRoleDao.DeleteRoleByID(db.Mysql, id)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (svc *SysRoleServiceImpl) GetRoleList(query *request.PageQuery) (*response.PageInfo, *e.Error) {
	var roleQuery *request.SysRoleForList
	if query.Query != nil {
		roleQuery = query.Query.(*request.SysRoleForList)
	}
	// 获取角色列表
	sysSysRoles, err := svc.sysRoleDao.GetRoleList(db.Mysql, query)
	if err != nil {
		return nil, e.ErrMysql
	}
	newSysRoles := make([]*dto.SysRoleDto, len(sysSysRoles))
	for i := 0; i < len(sysSysRoles); i++ {
		newSysRoles[i] = dto.NewSysRoleDto(sysSysRoles[i])
	}
	// 获取所有角色总数目
	var count int64
	count, err = svc.sysRoleDao.GetRoleCount(db.Mysql, roleQuery)
	if err != nil {
		return nil, e.ErrMysql
	}
	pageInfo := &response.PageInfo{
		Total: count,
		Size:  int64(len(newSysRoles)),
		List:  newSysRoles,
	}
	return pageInfo, nil
}

func (svc *SysRoleServiceImpl) UpdateRolePermissions(roleID uint, permissionIDs []uint) *e.Error {
	tx := db.Mysql.Begin()
	err := svc.sysRoleDao.DisGrantPermissionsByRoleID(tx, roleID)
	if err != nil {
		tx.Rollback()
		return e.ErrMysql
	}
	err = svc.sysRoleDao.GrantPermissionsToRole(tx, roleID, permissionIDs)
	if err != nil {
		tx.Rollback()
		return e.ErrMysql
	}
	tx.Commit()
	return nil
}

func (svc *SysRoleServiceImpl) GetPermissionIDsByRoleID(roleID uint) ([]uint, *e.Error) {
	permissionIDs, err := svc.sysRoleDao.GetPermissionIDsByRoleID(db.Mysql, roleID)
	if err != nil {
		return nil, e.ErrMysql
	}
	return permissionIDs, nil
}

func (svc *SysRoleServiceImpl) GetPermissionsByRoleID(roleID uint) ([]*repository.SysPermission, *e.Error) {
	permissions, err := svc.sysRoleDao.GetPermissionsByRoleID(db.Mysql, roleID)
	if err != nil {
		return nil, e.ErrMysql
	}
	return permissions, nil
}
