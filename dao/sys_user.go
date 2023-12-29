package dao

import (
	"errors"
	"funoj-backend/model/form/request"
	"funoj-backend/model/repository"
	"gorm.io/gorm"
)

type SysUserDao interface {
	// InsertUser 添加用户
	InsertUser(db *gorm.DB, user *repository.SysUser) error
	// DeleteUserByID 通过用户id删除用户
	DeleteUserByID(db *gorm.DB, userID uint) error
	// UpdateUser 更新用户信息
	UpdateUser(db *gorm.DB, user *repository.SysUser) error
	// GetUserByID 通过用户id获取用户信息
	GetUserByID(db *gorm.DB, userID uint) (*repository.SysUser, error)
	// GetUserNameByID 通过用户id获取用户名
	GetUserNameByID(db *gorm.DB, userID uint) (string, error)
	// GetUserByLoginName 通过登录名获取用户信息
	GetUserByLoginName(db *gorm.DB, loginName string) (*repository.SysUser, error)
	// GetUserByEmail 通过邮箱获取用户信息
	GetUserByEmail(db *gorm.DB, email string) (*repository.SysUser, error)
	// GetUserList 获取用户列表
	GetUserList(db *gorm.DB, pageQuery request.PageQuery) ([]*repository.SysUser, error)
	// GetUserCount 获取用户数
	GetUserCount(db *gorm.DB, user *request.SysUserForList) (int64, error)
	// CheckLoginName 校验登录名
	CheckLoginName(db *gorm.DB, loginName string) (bool, error)
	// ListLoginName 获取登录名列表
	ListLoginName(db *gorm.DB, loginName string) ([]string, error)
	// CheckEmail 校验邮箱
	CheckEmail(db *gorm.DB, email string) (bool, error)
	// DeleteUserRoleByUserID 通过用户id删除用户角色
	DeleteUserRoleByUserID(db *gorm.DB, userID uint) error
	// GetRolesByUserID 通过用户id获取用户的角色身份
	GetRolesByUserID(db *gorm.DB, userID uint) ([]*repository.SysRole, error)
	// InsertRolesToUser 归档用户到对应角色组
	InsertRolesToUser(db *gorm.DB, userID uint, roleIDs []uint) error
	// UpdateUserPassword 更新用户密码
	UpdateUserPassword(db *gorm.DB, userID uint, password string) error
}

type SysUserDaoImpl struct {
}

func NewSysUserDao() SysUserDao {
	return &SysUserDaoImpl{}
}

func (s *SysUserDaoImpl) InsertUser(db *gorm.DB, user *repository.SysUser) error {
	return db.Create(user).Error
}

func (s *SysUserDaoImpl) DeleteUserByID(db *gorm.DB, userID uint) error {
	return db.Delete(&repository.SysUser{}, "id = ?", userID).Error
}

func (s *SysUserDaoImpl) UpdateUser(db *gorm.DB, user *repository.SysUser) error {
	return db.Model(user).Updates(user).Error
}

func (s *SysUserDaoImpl) GetUserByID(db *gorm.DB, userID uint) (*repository.SysUser, error) {
	user := &repository.SysUser{}
	err := db.Where("id = ?", userID).Find(&user).Error
	return user, err
}

func (s *SysUserDaoImpl) GetUserNameByID(db *gorm.DB, userID uint) (string, error) {
	user := &repository.SysUser{}
	err := db.Where("id = ?", userID).Select("user_name").Find(&user).Error
	return user.UserName, err
}

func (s *SysUserDaoImpl) GetUserByLoginName(db *gorm.DB, loginName string) (*repository.SysUser, error) {
	user := &repository.SysUser{}
	err := db.Where("login_name = ?", loginName).Find(&user).Error
	return user, err
}

func (s *SysUserDaoImpl) GetUserByEmail(db *gorm.DB, email string) (*repository.SysUser, error) {
	user := &repository.SysUser{}
	err := db.Where("email = ?", email).Find(&user).Error
	return user, err
}

func (s *SysUserDaoImpl) GetUserList(db *gorm.DB, pageQuery request.PageQuery) ([]*repository.SysUser, error) {
	var user *request.SysUserForList
	if pageQuery.Query != nil {
		user = pageQuery.Query.(*request.SysUserForList)
	}
	offset := (pageQuery.Page - 1) * pageQuery.PageSize
	var users []*repository.SysUser
	if user != nil && user.UserName != "" {
		db = db.Where("user_name LIKE ?", "%"+user.UserName+"%")
	}
	if user != nil && user.Gender != 0 {
		db = db.Where("gender = ?", user.Gender)
	}
	if user != nil && user.LoginName != "" {
		db = db.Where("login_name LIKE ?", "%"+user.LoginName+"%")
	}
	if user != nil && user.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+user.Phone+"%")
	}
	if user != nil && user.Email != "" {
		db = db.Where("email LIKE ?", "%"+user.Email+"%")
	}
	err := db.Limit(pageQuery.PageSize).Offset(offset).Find(&users).Error
	return users, err
}

func (s *SysUserDaoImpl) GetUserCount(db *gorm.DB, user *request.SysUserForList) (int64, error) {
	var count int64
	if user != nil && user.UserName != "" {
		db = db.Where("user_name LIKE ?", "%"+user.UserName+"%")
	}
	if user != nil && user.Gender != 0 {
		db = db.Where("gender = ?", user.Gender)
	}
	if user != nil && user.LoginName != "" {
		db = db.Where("login_name LIKE ?", "%"+user.LoginName+"%")
	}
	if user != nil && user.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+user.Phone+"%")
	}
	if user != nil && user.Email != "" {
		db = db.Where("email LIKE ?", "%"+user.Email+"%")
	}
	err := db.Model(&repository.SysUser{}).Count(&count).Error
	return count, err
}

func (s *SysUserDaoImpl) CheckLoginName(db *gorm.DB, loginName string) (bool, error) {
	var user *repository.SysUser
	err := db.Where("login_name = ?", loginName).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *SysUserDaoImpl) ListLoginName(db *gorm.DB, loginName string) ([]string, error) {
	var users []*repository.SysUser
	err := db.Model(&repository.SysUser{}).Where("login_name like ?", loginName+"%").
		Select("login_name").Find(&users).Error
	if err != nil {
		return nil, err
	}
	answer := make([]string, len(users))
	for i := 0; i < len(users); i++ {
		answer[i] = users[i].LoginName
	}
	return answer, nil
}

func (s *SysUserDaoImpl) CheckEmail(db *gorm.DB, email string) (bool, error) {
	var user *repository.SysUser
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return user.ID != 0, nil
}

func (s *SysUserDaoImpl) DeleteUserRoleByUserID(db *gorm.DB, userID uint) error {
	user := repository.SysUser{}
	user.ID = userID
	if err := db.Model(&user).Association("Roles").Clear(); err != nil {
		return err
	}
	return nil
}

func (s *SysUserDaoImpl) GetRolesByUserID(db *gorm.DB, userID uint) ([]*repository.SysRole, error) {
	user := repository.SysUser{}
	user.ID = userID
	if err := db.Model(&user).Association("Roles").Find(&user.Roles); err != nil {
		return nil, err
	}
	return user.Roles, nil
}

func (s *SysUserDaoImpl) GetRoleIDsByUserID(db *gorm.DB, userID uint) ([]uint, error) {
	user := repository.SysUser{}
	user.ID = userID
	if err := db.Model(&user).Association("Roles").Find(&user.Roles); err != nil {
		return nil, err
	}
	roleIDs := make([]uint, len(user.Roles))
	for i, role := range user.Roles {
		roleIDs[i] = role.ID
	}
	return roleIDs, nil
}

func (s *SysUserDaoImpl) InsertRolesToUser(db *gorm.DB, userID uint, roleIDs []uint) error {
	user := &repository.SysUser{}
	user.ID = userID
	for _, roleID := range roleIDs {
		role := &repository.SysRole{}
		role.ID = roleID
		err := db.Model(user).Association("Roles").Append(role)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SysUserDaoImpl) UpdateUserPassword(db *gorm.DB, userID uint, password string) error {
	user := repository.SysUser{}
	user.ID = userID
	return db.Model(&user).Update("password", password).Error
}
