package dao

import (
	"funoj-backend/model/repository"
	"gorm.io/gorm"
)

type ProblemAttemptDao interface {
	// InsertProblemAttempt 新增用户对题目提交情况
	InsertProblemAttempt(db *gorm.DB, problemAttempt *repository.ProblemAttempt) error
	// UpdateProblemAttempt 更新用户对题目提交情况
	UpdateProblemAttempt(db *gorm.DB, problemAttempt *repository.ProblemAttempt) error
	// GetProblemAttemptByID 通过用户id和题目id查询用户对题目提交情况
	GetProblemAttemptByID(db *gorm.DB, userId uint, problemId uint) (*repository.ProblemAttempt, error)
	// GetProblemAttemptStatus 通过用户id和题目id查询用户对题目提交状态
	GetProblemAttemptStatus(db *gorm.DB, userId uint, problemID uint) (int, error)
}

type ProblemAttemptDaoImpl struct {
}

func NewProblemAttemptDao() ProblemAttemptDao {
	return &ProblemAttemptDaoImpl{}
}

func (dao *ProblemAttemptDaoImpl) InsertProblemAttempt(db *gorm.DB, problemAttempt *repository.ProblemAttempt) error {
	return db.Create(problemAttempt).Error
}

func (dao *ProblemAttemptDaoImpl) UpdateProblemAttempt(db *gorm.DB, problemAttempt *repository.ProblemAttempt) error {
	return db.Model(problemAttempt).UpdateColumns(map[string]interface{}{
		"submission_count": problemAttempt.SubmissionCount,
		"success_count":    problemAttempt.SuccessCount,
		"err_count":        problemAttempt.ErrCount,
		"code":             problemAttempt.Code,
		"status":           problemAttempt.Status,
		"updated_at":       problemAttempt.UpdatedAt,
	}).Error
}

func (dao *ProblemAttemptDaoImpl) GetProblemAttemptByID(db *gorm.DB, userId uint, problemId uint) (*repository.ProblemAttempt, error) {
	problemAttempt := repository.ProblemAttempt{}
	err := db.Model(&repository.ProblemAttempt{}).Where("user_id = ? and problem_id = ?", userId, problemId).
		First(&problemAttempt).Error
	return &problemAttempt, err
}

func (dao *ProblemAttemptDaoImpl) GetProblemAttemptStatus(db *gorm.DB, userId uint, problemID uint) (int, error) {
	var problemAttempt repository.ProblemAttempt
	err := db.Model(&repository.ProblemAttempt{}).Select("status", "id").
		Where("user_id = ? and problem_id = ?", userId, problemID).First(&problemAttempt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return problemAttempt.Status, nil
}
