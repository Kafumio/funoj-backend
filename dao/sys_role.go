package dao

import (
	"funoj-backend/model/form/request"
	"funoj-backend/model/repository"
	"gorm.io/gorm"
)

type SysRoleDao interface {
	// InsertRole 创建角色
	InsertRole(db *gorm.DB, role *repository.SysRole) error
	// UpdateRole 更新角色
	UpdateRole(db *gorm.DB, role *repository.SysRole) error
	// DeleteRoleByID 删除角色
	DeleteRoleByID(db *gorm.DB, id uint) error
	// GetRoleByID 通过角色id获取角色
	GetRoleByID(db *gorm.DB, roleID uint) (*repository.SysRole, error)
	// GetRoleList 获取角色列表
	GetRoleList(db *gorm.DB, pageQuery *request.PageQuery) ([]*repository.SysRole, error)
	// GetRoleCount 获取所有角色数量
	GetRoleCount(db *gorm.DB, role *request.SysRoleForList) (int64, error)
	// GrantPermissionsToRole 赋予角色权限
	GrantPermissionsToRole(db *gorm.DB, roleID uint, permissionIDs []uint) error
	// DisGrantPermissionsByRoleID 通过RoleID回收对应角色的权限
	DisGrantPermissionsByRoleID(db *gorm.DB, roleID uint) error
	// GetPermissionIDsByRoleID 获取用户关联的所有权限的id
	GetPermissionIDsByRoleID(db *gorm.DB, roleID uint) ([]uint, error)
	// GetPermissionsByRoleID 获取用户关联的所有权限
	GetPermissionsByRoleID(db *gorm.DB, roleID uint) ([]*repository.SysPermission, error)
	// GetAllSimpleRoleList 获取所有角色列表，只含有id和name
	GetAllSimpleRoleList(db *gorm.DB) ([]*repository.SysRole, error)
}

type SysRoleDaoImpl struct {
}

func NewSysRoleDao() SysRoleDao {
	return &SysRoleDaoImpl{}
}

func (s *SysRoleDaoImpl) InsertRole(db *gorm.DB, role *repository.SysRole) error {
	return db.Create(role).Error
}

func (s *SysRoleDaoImpl) UpdateRole(db *gorm.DB, role *repository.SysRole) error {
	return db.Model(role).Updates(role).Error
}

func (s *SysRoleDaoImpl) DeleteRoleByID(db *gorm.DB, id uint) error {
	return db.Delete(&repository.SysRole{}, id).Error
}

func (s *SysRoleDaoImpl) GetRoleByID(db *gorm.DB, id uint) (*repository.SysRole, error) {
	role := &repository.SysRole{}
	err := db.Where("id = ?", id).Find(&role).Error
	return role, err
}

func (r *SysRoleDaoImpl) GetRoleList(db *gorm.DB, pageQuery *request.PageQuery) ([]*repository.SysRole, error) {
	var role *request.SysRoleForList
	if pageQuery.Query != nil {
		role = pageQuery.Query.(*request.SysRoleForList)
	}
	offset := (pageQuery.Page - 1) * pageQuery.PageSize
	var roles []*repository.SysRole
	if role != nil && role.Name != "" {
		db = db.Where("name LIKE ?", "%"+role.Name+"%")
	}
	if role != nil && role.Description != "" {
		db = db.Where("description LIKE ?", "%"+role.Description+"%")
	}
	db = db.Limit(pageQuery.PageSize).Offset(offset)
	if pageQuery.SortProperty != "" && pageQuery.SortRule != "" {
		order := pageQuery.SortProperty + " " + pageQuery.SortRule
		db = db.Order(order)
	}
	err := db.Find(&roles).Error
	return roles, err
}

func (s *SysRoleDaoImpl) GetRoleCount(db *gorm.DB, role *request.SysRoleForList) (int64, error) {
	var count int64
	if role != nil && role.Name != "" {
		db = db.Where("name LIKE ?", "%"+role.Name+"%")
	}
	if role != nil && role.Description != "" {
		db = db.Where("description LIKE ?", "%"+role.Description+"%")
	}
	err := db.Model(&repository.SysRole{}).Count(&count).Error
	return count, err
}

func (s *SysRoleDaoImpl) GrantPermissionsToRole(db *gorm.DB, roleID uint, permissionIDs []uint) error {
	role := &repository.SysRole{}
	role.ID = roleID
	var permissions []repository.SysPermission
	for _, permissionID := range permissionIDs {
		permission := repository.SysPermission{}
		permission.ID = permissionID
		permissions = append(permissions, permission)
	}
	err := db.Model(role).Association("Permissions").Append(permissions)
	return err
}

func (s *SysRoleDaoImpl) DisGrantPermissionsByRoleID(db *gorm.DB, roleID uint) error {
	role := repository.SysRole{}
	role.ID = roleID
	err := db.Model(&role).Association("Permissions").Clear()
	return err
}

func (s *SysRoleDaoImpl) GetPermissionIDsByRoleID(db *gorm.DB, roleID uint) ([]uint, error) {
	var role repository.SysRole
	role.ID = roleID
	if err := db.Model(&role).Association("Permissions").Find(&role.Permissions); err != nil {
		return nil, err
	}
	permissionIDs := make([]uint, len(role.Permissions))
	for i, permission := range role.Permissions {
		permissionIDs[i] = permission.ID
	}
	return permissionIDs, nil
}

func (s *SysRoleDaoImpl) GetPermissionsByRoleID(db *gorm.DB, roleID uint) ([]*repository.SysPermission, error) {
	role := repository.SysRole{}
	role.ID = roleID
	if err := db.Model(&role).Association("Permissions").Find(&role.Permissions); err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

func (r *SysRoleDaoImpl) GetAllSimpleRoleList(db *gorm.DB) ([]*repository.SysRole, error) {
	var roles []*repository.SysRole
	err := db.Select([]string{"id", "name"}).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
