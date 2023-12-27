package dao

import (
	"funoj-backend/model/dto"
	"funoj-backend/model/form/request"
	"funoj-backend/model/repository"
	"gorm.io/gorm"
)

type ProblemDao interface {
	// GetProblemByNumber 根据题目编码获取题目
	GetProblemByNumber(db *gorm.DB, problemCode string) (*repository.Problem, error)
	// GetProblemIDByNumber 根据题目number获取题目id
	GetProblemIDByNumber(db *gorm.DB, problemNumber string) (uint, error)
	// GetProblemNameByID 根据题目id获取题目名称
	GetProblemNameByID(db *gorm.DB, problemID uint) (string, error)
	// GetProblemByID 根据题目id获取题目
	GetProblemByID(db *gorm.DB, problemID uint) (*repository.Problem, error)
	// InsertProblem 添加题库
	InsertProblem(db *gorm.DB, problem *repository.Problem) error
	// UpdateProblem 更新题目
	// 不修改path
	UpdateProblem(db *gorm.DB, problem *repository.Problem) error
	// UpdateProblemField 根据字段进行更新
	UpdateProblemField(db *gorm.DB, id uint, field string, value string) error
	// CheckProblemNumberExists 检测用户ID是否存在
	CheckProblemNumberExists(db *gorm.DB, problemCode string) (bool, error)
	// SetProblemEnable 让一个题目可用
	SetProblemEnable(db *gorm.DB, id uint, enable int) error
	// DeleteProblemByID 删除题目
	DeleteProblemByID(db *gorm.DB, id uint) error
	GetProblemList(db *gorm.DB, pageQuery *dto.PageQuery) ([]*repository.Problem, error)
	GetProblemCount(db *gorm.DB, problem *request.Problem) (int64, error)
}

type ProblemDaoImpl struct {
}

func NewProblemDao() ProblemDao {
	return &ProblemDaoImpl{}
}

func (dao *ProblemDaoImpl) GetProblemByNumber(db *gorm.DB, problemNumber string) (*repository.Problem, error) {
	problem := &repository.Problem{}
	err := db.Where("number = ?", problemNumber).Find(&problem).Error
	return problem, err
}

func (dao *ProblemDaoImpl) GetProblemIDByNumber(db *gorm.DB, problemNumber string) (uint, error) {
	problem, err := dao.GetProblemByNumber(db, problemNumber)
	return problem.ID, err
}

func (dao *ProblemDaoImpl) GetProblemByID(db *gorm.DB, problemID uint) (*repository.Problem, error) {
	problem := &repository.Problem{}
	err := db.First(&problem, problemID).Error
	return problem, err
}

func (dao *ProblemDaoImpl) GetProblemNameByID(db *gorm.DB, problemID uint) (string, error) {
	problem, err := dao.GetProblemByID(db, problemID)
	return problem.Name, err
}

func (dao *ProblemDaoImpl) GetProblemList(db *gorm.DB, pageQuery *dto.PageQuery) ([]*repository.Problem, error) {
	var problem *request.Problem
	if pageQuery.Query != nil {
		problem = pageQuery.Query.(*request.Problem)
	}
	if problem != nil && problem.MenuID != nil {
		db = db.Preload("Menus", func(db *gorm.DB) *gorm.DB {
			return db.Where("id = ?", problem.MenuID)
		})
	}
	if problem != nil && problem.Number != "" {
		db = db.Where("number like ?", "%"+problem.Number+"%")
	}
	if problem != nil && problem.Name != "" {
		db = db.Where("name like ?", "%"+problem.Name+"%")
	}
	if problem != nil && problem.Difficulty != 0 {
		db = db.Where("difficulty = ?", problem.Difficulty)
	}
	if problem != nil && problem.Enable != 0 {
		db = db.Where("enable = ?", problem.Enable)
	}
	offset := (pageQuery.Page - 1) * pageQuery.PageSize
	var problems []*repository.Problem
	db = db.Offset(offset).Limit(pageQuery.PageSize)
	if pageQuery.SortProperty != "" && pageQuery.SortRule != "" {
		order := pageQuery.SortProperty + " " + pageQuery.SortRule
		db = db.Order(order)
	}
	err := db.Find(&problems).Error
	return problems, err
}

func (dao *ProblemDaoImpl) GetProblemCount(db *gorm.DB, problem *request.Problem) (int64, error) {
	var count int64
	if problem != nil && problem.MenuID != nil {
		db = db.Preload("Menus", func(db *gorm.DB) *gorm.DB {
			return db.Where("id = ?", problem.MenuID)
		})
	}
	if problem != nil && problem.Name != "" {
		db = db.Where("name like ?", "%"+problem.Name+"%")
	}
	if problem != nil && problem.Number != "" {
		db = db.Where("number = ?", problem.Number)
	}
	if problem != nil && problem.Difficulty != 0 {
		db = db.Where("difficulty = ?", problem.Difficulty)
	}
	if problem != nil && problem.Enable != 0 {
		db = db.Where("enable = ?", problem.Enable)
	}
	err := db.Model(&repository.Problem{}).Count(&count).Error
	return count, err
}

func (dao *ProblemDaoImpl) InsertProblem(db *gorm.DB, problem *repository.Problem) error {
	return db.Create(problem).Error
}

func (dao *ProblemDaoImpl) UpdateProblem(db *gorm.DB, problem *repository.Problem) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&problem).Where("id = ?", problem.ID).Updates(map[string]interface{}{
			"updated_at":  problem.UpdatedAt,
			"number":      problem.Number,
			"name":        problem.Name,
			"description": problem.Description,
			"difficulty":  problem.Difficulty,
			"title":       problem.Title,
			"languages":   problem.Languages,
			"enable":      problem.Enable,
		}).Error; err != nil {
			return err
		}
		if err := tx.Model(&problem).Association("Menus").Replace(&problem.Menus); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (dao *ProblemDaoImpl) UpdateProblemField(db *gorm.DB, id uint, field string, value string) error {
	updateData := map[string]interface{}{
		field: value,
	}
	if err := db.Model(&repository.Problem{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}

func (p *ProblemDaoImpl) CheckProblemNumberExists(db *gorm.DB, problemNumber string) (bool, error) {
	//执行
	row := db.Model(&repository.Problem{}).Select("number").Where("number = ?", problemNumber)
	if row.Error != nil {
		return false, row.Error
	}
	problem := &repository.Problem{}
	row.Scan(&problem)
	return problem.Number != "", nil
}

func (p *ProblemDaoImpl) SetProblemEnable(db *gorm.DB, id uint, enable int) error {
	return db.Model(&repository.Problem{}).Where("id = ?", id).Update("enable", enable).Error
}

func (p *ProblemDaoImpl) DeleteProblemByID(db *gorm.DB, id uint) error {
	return db.Delete(&repository.Problem{}, id).Error
}
