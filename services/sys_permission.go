package services

import (
	e "funoj-backend/consts/error"
	"funoj-backend/dao"
	"funoj-backend/db"
	"funoj-backend/model/dto"
	"funoj-backend/model/repository"
	"gorm.io/gorm"
	"time"
)

type SysPermissionService interface {
	// GetPermissionCount 获取权限数目
	GetPermissionCount() (int64, *e.Error)
	// DeletePermissionByID 删除权限
	DeletePermissionByID(id uint) *e.Error
	// UpdatePermission 更新权限
	UpdatePermission(permission *repository.SysPermission) *e.Error
	// InsertPermission 添加权限
	InsertPermission(permission *repository.SysPermission) (uint, *e.Error)
	// GetPermissionByID 根据id获取权限
	GetPermissionByID(id uint) (*repository.SysPermission, *e.Error)
	// GetPermissionTree 获取权限树
	GetPermissionTree() ([]*dto.SysPermissionTreeDto, *e.Error)
}

type SysPermissionServiceImpl struct {
	sysPermissionDao dao.SysPermissionDao
}

func NewSysPermissionService(permissionDao dao.SysPermissionDao) SysPermissionService {
	return &SysPermissionServiceImpl{
		sysPermissionDao: permissionDao,
	}
}

func (svc *SysPermissionServiceImpl) GetPermissionCount() (int64, *e.Error) {
	count, err := svc.sysPermissionDao.GetPermissionCount(db.Mysql)
	if err != nil {
		return 0, e.ErrMysql
	}
	return count, nil
}

func (svc *SysPermissionServiceImpl) DeletePermissionByID(id uint) *e.Error {
	err := db.Mysql.Transaction(func(tx *gorm.DB) error {
		// 递归删除权限
		err := svc.deletePermissionsRecursive(tx, id)
		return err
	})
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

// deletePermissionsRecursive 递归删除权限
func (svc *SysPermissionServiceImpl) deletePermissionsRecursive(db *gorm.DB, parentID uint) error {
	childPermissions, err := svc.sysPermissionDao.GetChildPermissionsByParentID(db, parentID)
	if err != nil {
		return err
	}
	for _, childPermission := range childPermissions {
		// 删除子api的子api
		if err = svc.deletePermissionsRecursive(db, childPermission.ID); err != nil {
			return err
		}
	}
	// 当前api
	if err = svc.sysPermissionDao.DeletePermissionByID(db, parentID); err != nil {
		return err
	}
	return nil
}

func (svc *SysPermissionServiceImpl) UpdatePermission(permission *repository.SysPermission) *e.Error {
	permission.UpdatedAt = time.Now()
	err := svc.sysPermissionDao.UpdatePermission(db.Mysql, permission)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (svc *SysPermissionServiceImpl) InsertPermission(permission *repository.SysPermission) (uint, *e.Error) {
	err := svc.sysPermissionDao.InsertPermission(db.Mysql, permission)
	if err != nil {
		return 0, e.ErrMysql
	}
	return permission.ID, nil
}

func (s *SysPermissionServiceImpl) GetPermissionByID(id uint) (*repository.SysPermission, *e.Error) {
	permission, err := s.sysPermissionDao.GetPermissionByID(db.Mysql, id)
	if err == gorm.ErrRecordNotFound {
		return nil, e.ErrPermissionNotExist
	}
	if err != nil {
		return nil, e.ErrMysql
	}
	return permission, nil
}

func (svc *SysPermissionServiceImpl) GetPermissionTree() ([]*dto.SysPermissionTreeDto, *e.Error) {
	var permissionList []*repository.SysPermission
	var err error
	if permissionList, err = svc.sysPermissionDao.GetAllPermission(db.Mysql); err != nil {
		return nil, e.ErrMysql
	}

	permissionMap := make(map[uint]*dto.SysPermissionTreeDto)
	var rootPermissions []*dto.SysPermissionTreeDto

	// 添加到map中保存
	for _, permission := range permissionList {
		permissionMap[permission.ID] = dto.NewSysApiTreeDto(permission)
	}

	// 遍历并添加到父节点中
	for _, permission := range permissionList {
		if permission.ParentID == 0 {
			rootPermissions = append(rootPermissions, permissionMap[permission.ID])
		} else {
			parentPermission, exists := permissionMap[permission.ParentID]
			if !exists {
				return nil, e.ErrPermissionUnknownError
			}
			parentPermission.Children = append(parentPermission.Children, permissionMap[permission.ID])
		}
	}

	return rootPermissions, nil
}
