package dao

import (
	"funoj-backend/model/repository"
	"gorm.io/gorm"
)

type SysPermissionDao interface {
	// InsertPermission 增加系统权限
	InsertPermission(db *gorm.DB, permission *repository.SysPermission) error
	// UpdatePermission 修改系统权限
	UpdatePermission(db *gorm.DB, permission *repository.SysPermission) error
	// DeletePermissionByID 根据权限的id进行删除
	DeletePermissionByID(db *gorm.DB, id uint) error
	// GetPermissionByID 通过权限的id获取权限
	GetPermissionByID(db *gorm.DB, id uint) (*repository.SysPermission, error)
	// GetPermissionCount 获取权限总数
	GetPermissionCount(db *gorm.DB) (int64, error)
	// GetPermissionListByPathKeyword 模糊查询权限
	GetPermissionListByPathKeyword(db *gorm.DB, keyword string, page int, pageSize int) ([]*repository.SysPermission, error)
	// GetChildPermissionsByParentID 根据父权限的ID获取所有子权限
	GetChildPermissionsByParentID(db *gorm.DB, parentID uint) ([]*repository.SysPermission, error)
	// GetAllPermission 获取所有权限
	GetAllPermission(db *gorm.DB) ([]*repository.SysPermission, error)
}

type SysPermissionDaoImpl struct {
}

func NewSysPermissionDao() SysPermissionDao {
	return &SysPermissionDaoImpl{}
}

func (dao *SysPermissionDaoImpl) InsertPermission(db *gorm.DB, permission *repository.SysPermission) error {
	return db.Create(permission).Error
}

func (dao *SysPermissionDaoImpl) UpdatePermission(db *gorm.DB, permission *repository.SysPermission) error {
	return db.Model(permission).Updates(permission).Error
}

func (dao *SysPermissionDaoImpl) DeletePermissionByID(db *gorm.DB, id uint) error {
	return db.Delete(&repository.SysPermission{}, id).Error
}

func (dao *SysPermissionDaoImpl) GetAllPermission(db *gorm.DB) ([]*repository.SysPermission, error) {
	var permissionList []*repository.SysPermission
	err := db.Find(&permissionList).Error
	return permissionList, err
}

func (dao *SysPermissionDaoImpl) GetPermissionByID(db *gorm.DB, id uint) (*repository.SysPermission, error) {
	var permission *repository.SysPermission
	err := db.Find(&permission).Error
	return permission, err
}

func (dao *SysPermissionDaoImpl) GetPermissionCount(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&repository.SysPermission{}).Count(&count).Error
	return count, err
}

func (dao *SysPermissionDaoImpl) GetPermissionListByPathKeyword(db *gorm.DB, keyword string, page int, pageSize int) ([]*repository.SysPermission, error) {
	var permissions []*repository.SysPermission
	err := db.Where("path LIKE ?", "%"+keyword+"%").Offset((page - 1) * pageSize).Limit(pageSize).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (dao *SysPermissionDaoImpl) GetChildPermissionsByParentID(db *gorm.DB, parentID uint) ([]*repository.SysPermission, error) {
	var childPermissions []*repository.SysPermission
	if err := db.Where("parent_id = ?", parentID).Find(&childPermissions).Error; err != nil {
		return nil, err
	}
	return childPermissions, nil
}
