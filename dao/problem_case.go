package dao

import (
	"funoj-backend/model/form/request"
	"funoj-backend/model/repository"
	"gorm.io/gorm"
)

type ProblemCaseDao interface {
	// GetProblemCaseList 获取用例列表
	GetProblemCaseList(db *gorm.DB, query *request.PageQuery) ([]*repository.ProblemCase, error)
	GetAllProblemCaseByID(db *gorm.DB, problemID uint) ([]*repository.ProblemCase, error)
	// GetProblemCaseCount 获取用例数量
	GetProblemCaseCount(db *gorm.DB, problemCase *request.ProblemCaseForList) (int64, error)
	// GetProblemCaseByID 通过id获取题目用例
	GetProblemCaseByID(db *gorm.DB, id uint) (*repository.ProblemCase, error)
	// DeleteProblemCaseByID 通过id删除题目用例
	DeleteProblemCaseByID(db *gorm.DB, id uint) error
	// DeleteProblemCaseByProblemID 通过题目id删除题目用例
	DeleteProblemCaseByProblemID(db *gorm.DB, problemID uint) error
	// InsertProblemCase 添加题目用例
	InsertProblemCase(db *gorm.DB, problemCase *repository.ProblemCase) error
	// UpdateProblemCase 更新题目用例
	UpdateProblemCase(db *gorm.DB, problemCase *repository.ProblemCase) error
}

type ProblemCaseDaoImpl struct {
}

func NewProblemCaseDao() ProblemCaseDao {
	return &ProblemCaseDaoImpl{}
}

func (dao *ProblemCaseDaoImpl) GetProblemCaseList(db *gorm.DB, query *request.PageQuery) ([]*repository.ProblemCase, error) {
	var problemCase *request.ProblemCaseForList
	if query.Query != nil {
		problemCase = query.Query.(*request.ProblemCaseForList)
	}
	if problemCase != nil && problemCase.ProblemID != 0 {
		db = db.Where("problem_id = ?", problemCase.ProblemID)
	}
	if problemCase != nil && problemCase.CaseName != "" {
		db = db.Where("case_name like ?", "%"+problemCase.CaseName+"%")
	}
	offset := (query.Page - 1) * query.PageSize
	var cases []*repository.ProblemCase
	db = db.Offset(offset).Limit(query.PageSize)
	if query.SortProperty != "" && query.SortRule != "" {
		order := query.SortProperty + " " + query.SortRule
		db = db.Order(order)
	}
	err := db.Find(&cases).Error
	return cases, err
}

func (dao *ProblemCaseDaoImpl) GetAllProblemCaseByID(db *gorm.DB, problemID uint) ([]*repository.ProblemCase, error) {
	db = db.Where("problem_id = ?", problemID)
	var cases []*repository.ProblemCase
	err := db.Find(&cases).Error
	return cases, err
}

func (dao *ProblemCaseDaoImpl) GetProblemCaseCount(db *gorm.DB, problemCase *request.ProblemCaseForList) (int64, error) {
	var count int64
	if problemCase != nil && problemCase.ProblemID != 0 {
		db = db.Where("problem_id = ?", problemCase.ProblemID)
	}
	if problemCase != nil && problemCase.CaseName != "" {
		db = db.Where("case_name like ?", "%"+problemCase.CaseName+"%")
	}
	err := db.Model(&repository.ProblemCase{}).Count(&count).Error
	return count, err
}

func (dao *ProblemCaseDaoImpl) GetProblemCaseByID(db *gorm.DB, id uint) (*repository.ProblemCase, error) {
	problemCase := &repository.ProblemCase{}
	err := db.Where("id = ?", id).Find(&problemCase).Error
	return problemCase, err
}

func (dao *ProblemCaseDaoImpl) DeleteProblemCaseByID(db *gorm.DB, id uint) error {
	return db.Delete(&repository.ProblemCase{}, id).Error
}

func (dao *ProblemCaseDaoImpl) DeleteProblemCaseByProblemID(db *gorm.DB, problemID uint) error {
	return db.Where("problem_id = ?", problemID).Delete(&repository.ProblemCase{}).Error
}

func (dao *ProblemCaseDaoImpl) InsertProblemCase(db *gorm.DB, problemCase *repository.ProblemCase) error {
	return db.Create(problemCase).Error
}

func (dao *ProblemCaseDaoImpl) UpdateProblemCase(db *gorm.DB, problemCase *repository.ProblemCase) error {
	return db.Model(problemCase).Updates(problemCase).Error
}
