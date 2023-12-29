package dao

import (
	"funoj-backend/model/form/request"
	"funoj-backend/model/repository"
	"gorm.io/gorm"
)

type ProblemMenuDao interface {
	// InsertProblemMenu 添加题库
	InsertProblemMenu(db *gorm.DB, bank *repository.ProblemMenu) error
	// GetProblemMenuByID 根据题库id获取题库
	GetProblemMenuByID(db *gorm.DB, bankID uint) (*repository.ProblemMenu, error)
	// UpdateProblemMenu 更新题库
	UpdateProblemMenu(db *gorm.DB, bank *repository.ProblemMenu) error
	// DeleteProblemMenuByID 删除题库
	DeleteProblemMenuByID(db *gorm.DB, id uint) error
	// GetProblemMenuCount 读取题库数量
	GetProblemMenuCount(db *gorm.DB, problemBank *request.ProblemMenuForList) (int64, error)
	// GetProblemMenuList 获取题库列表
	GetProblemMenuList(db *gorm.DB, pageQuery *request.PageQuery) ([]*repository.ProblemMenu, error)
	// GetAllProblemMenu 获取所有的题目数据
	GetAllProblemMenu(db *gorm.DB) ([]*repository.ProblemMenu, error)
	// GetSimpleProblemMenuList 获取题库列表，只包含id和名称
	GetSimpleProblemMenuList(db *gorm.DB) ([]*repository.ProblemMenu, error)
}

type ProblemMenuDaoImpl struct {
}

func NewProblemMenuDao() ProblemMenuDao {
	return &ProblemMenuDaoImpl{}
}

func (p *ProblemMenuDaoImpl) InsertProblemMenu(db *gorm.DB, problemMenu *repository.ProblemMenu) error {
	return db.Create(problemMenu).Error
}

func (p *ProblemMenuDaoImpl) GetProblemMenuByID(db *gorm.DB, id uint) (*repository.ProblemMenu, error) {
	problemMenu := &repository.ProblemMenu{}
	err := db.Where("id = ?", id).Find(&problemMenu).Error
	return problemMenu, err
}

func (p *ProblemMenuDaoImpl) UpdateProblemMenu(db *gorm.DB, problemMenu *repository.ProblemMenu) error {
	return db.Model(problemMenu).Updates(problemMenu).Error
}

func (p *ProblemMenuDaoImpl) DeleteProblemMenuByID(db *gorm.DB, id uint) error {
	return db.Delete(&repository.ProblemMenu{}, id).Error
}

func (p *ProblemMenuDaoImpl) GetProblemMenuCount(db *gorm.DB, problemMenu *request.ProblemMenuForList) (int64, error) {
	var count int64
	if problemMenu != nil && problemMenu.Name != "" {
		db = db.Where("name like ?", "%"+problemMenu.Name+"%")
	}
	if problemMenu != nil && problemMenu.Description != "" {
		db = db.Where("description = ?", problemMenu.Description)
	}
	err := db.Model(&repository.ProblemMenu{}).Count(&count).Error
	return count, err
}

func (p *ProblemMenuDaoImpl) GetProblemMenuList(db *gorm.DB, pageQuery *request.PageQuery) ([]*repository.ProblemMenu, error) {
	var problemMenu *request.ProblemMenuForList
	if pageQuery.Query != nil {
		problemMenu = pageQuery.Query.(*request.ProblemMenuForList)
	}
	if problemMenu != nil && problemMenu.Name != "" {
		db = db.Where("name like ?", "%"+problemMenu.Name+"%")
	}
	if problemMenu != nil && problemMenu.Description != "" {
		db = db.Where("description like ?", "%"+problemMenu.Description+"%")
	}
	offset := (pageQuery.Page - 1) * pageQuery.PageSize
	var menus []*repository.ProblemMenu
	db = db.Offset(offset).Limit(pageQuery.PageSize)
	if pageQuery.SortProperty != "" && pageQuery.SortRule != "" {
		order := pageQuery.SortProperty + " " + pageQuery.SortRule
		db = db.Order(order)
	}
	err := db.Find(&menus).Error
	return menus, err
}

func (p *ProblemMenuDaoImpl) GetAllProblemMenu(db *gorm.DB) ([]*repository.ProblemMenu, error) {
	var menus []*repository.ProblemMenu
	err := db.Find(&menus).Error
	return menus, err
}

func (p *ProblemMenuDaoImpl) GetSimpleProblemMenuList(db *gorm.DB) ([]*repository.ProblemMenu, error) {
	var menus []*repository.ProblemMenu
	err := db.Select("id", "name").Find(&menus).Error
	return menus, err
}
