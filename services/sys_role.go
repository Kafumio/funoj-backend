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
	// InsertSysRole 添加角色
	InsertSysRole(sysSysRole *repository.SysRole) (uint, *e.Error)
	// UpdateSysRole 更新角色
	UpdateSysRole(SysRole *repository.SysRole) *e.Error
	// DeleteSysRole 删除角色
	DeleteSysRole(id uint) *e.Error
	// GetSysRoleList 获取角色列表
	GetSysRoleList(pageQuery *request.PageQuery) (*response.PageInfo, *e.Error)
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

func (s *SysRoleServiceImpl) GetRoleByID(roleID uint) (*repository.SysRole, *e.Error) {
	role, err := s.sysRoleDao.GetRoleByID(db.Mysql, roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrRoleNotExist
		} else {
			return nil, e.ErrMysql
		}
	}
	return role, nil
}

func (s *SysRoleServiceImpl) InsertSysRole(role *repository.SysRole) (uint, *e.Error) {
	err := s.sysRoleDao.InsertRole(db.Mysql, role)
	if err != nil {
		return 0, e.ErrMysql
	}
	return role.ID, nil
}

func (r *SysRoleServiceImpl) UpdateSysRole(sysRole *repository.SysRole) *e.Error {
	sysRole.UpdatedAt = time.Now()
	err := r.sysRoleDao.UpdateRole(db.Mysql, sysRole)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (r *SysRoleServiceImpl) DeleteSysRole(id uint) *e.Error {
	// 删除删除角色
	err := r.sysRoleDao.DeleteRoleByID(db.Mysql, id)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (s *SysRoleServiceImpl) GetSysRoleList(query *request.PageQuery) (*response.PageInfo, *e.Error) {
	var roleQuery *request.SysRoleForList
	if query.Query != nil {
		roleQuery = query.Query.(*request.SysRoleForList)
	}
	// 获取角色列表
	sysSysRoles, err := s.sysRoleDao.GetRoleList(db.Mysql, query)
	if err != nil {
		return nil, e.ErrMysql
	}
	newSysRoles := make([]*dto.SysRoleDto, len(sysSysRoles))
	for i := 0; i < len(sysSysRoles); i++ {
		newSysRoles[i] = dto.NewSysRoleDto(sysSysRoles[i])
	}
	// 获取所有角色总数目
	var count int64
	count, err = s.sysRoleDao.GetRoleCount(db.Mysql, roleQuery)
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

func (s *SysRoleServiceImpl) UpdateRolePermissions(roleID uint, permissionIDs []uint) *e.Error {
	tx := db.Mysql.Begin()
	err := s.sysRoleDao.DisGrantPermissionsByRoleID(tx, roleID)
	if err != nil {
		tx.Rollback()
		return e.ErrMysql
	}
	err = s.sysRoleDao.GrantPermissionsToRole(tx, roleID, permissionIDs)
	if err != nil {
		tx.Rollback()
		return e.ErrMysql
	}
	tx.Commit()
	return nil
}

func (r *SysRoleServiceImpl) GetPermissionIDsByRoleID(roleID uint) ([]uint, *e.Error) {
	permissionIDs, err := r.sysRoleDao.GetPermissionIDsByRoleID(db.Mysql, roleID)
	if err != nil {
		return nil, e.ErrMysql
	}
	return permissionIDs, nil
}

func (r *SysRoleServiceImpl) GetPermissionsByRoleID(roleID uint) ([]*repository.SysPermission, *e.Error) {
	permissions, err := r.sysRoleDao.GetPermissionsByRoleID(db.Mysql, roleID)
	if err != nil {
		return nil, e.ErrMysql
	}
	return permissions, nil
}
