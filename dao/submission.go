package dao

import (
	"errors"
	"funoj-backend/model/form/request"
	"funoj-backend/model/repository"
	"gorm.io/gorm"
	"time"
)

type SubmissionDao interface {
	// GetLastSubmission 获取最后一次提交情况
	GetLastSubmission(db *gorm.DB, userID uint, problemID uint) (*repository.Submission, error)
	// GetSubmissionList 获取提交列表
	GetSubmissionList(db *gorm.DB, pageQuery *request.PageQuery) ([]*repository.Submission, error)
	// GetSubmissionCount 获取提交数
	GetSubmissionCount(db *gorm.DB, submission *request.SubmissionForList) (int64, error)
	// GetUserSimpleSubmissionsByTime 获取用户一段时间内的提交概况
	GetUserSimpleSubmissionsByTime(db *gorm.DB, userID uint, begin time.Time, end time.Time) ([]*repository.Submission, error)
	// CheckUserIsSubmittedByTime 检验用户是否在一段时间内进行过提交
	CheckUserIsSubmittedByTime(db *gorm.DB, userID uint, begin time.Time, end time.Time) (bool, error)
	// InsertSubmission 插入提交记录
	InsertSubmission(db *gorm.DB, submission *repository.Submission) error
}

type SubmissionDaoImpl struct {
}

func NewSubmissionDao() SubmissionDao {
	return &SubmissionDaoImpl{}
}

func (dao *SubmissionDaoImpl) GetLastSubmission(db *gorm.DB, userID uint, problemID uint) (*repository.Submission, error) {
	var submission *repository.Submission
	err := db.Where("user_id = ? and problem_id = ?", userID, problemID).Last(submission).Error
	return submission, err
}

func (dao *SubmissionDaoImpl) GetSubmissionList(db *gorm.DB, pageQuery *request.PageQuery) ([]*repository.Submission, error) {
	var submission *request.SubmissionForList
	if pageQuery.Query != nil {
		submission = pageQuery.Query.(*request.SubmissionForList)
	}
	var submissions []*repository.Submission
	if submission.UserID != 0 {
		db = db.Where("user_id = ?", submission.UserID)
	}
	if submission.ProblemID != 0 {
		db = db.Where("problem_id = ?", submission.ProblemID)
	}
	offset := (pageQuery.Page - 1) * pageQuery.PageSize
	err := db.Limit(pageQuery.PageSize).Offset(offset).Find(&submissions).Error
	return submissions, err
}

func (dao *SubmissionDaoImpl) GetSubmissionCount(db *gorm.DB, submission *request.SubmissionForList) (int64, error) {
	var count int64
	if submission != nil && submission.UserID != 0 {
		db = db.Where("user_id = ?", submission.UserID)
	}
	if submission != nil && submission.ProblemID != 0 {
		db = db.Where("problem_id = ?", submission.ProblemID)
	}
	err := db.Model(&repository.Submission{}).Count(&count).Error
	return count, err
}

func (dao *SubmissionDaoImpl) GetUserSimpleSubmissionsByTime(db *gorm.DB, userID uint, begin time.Time, end time.Time) ([]*repository.Submission, error) {
	var submissions []*repository.Submission
	err := db.Where("user_id = ? and created_at >= ? and created_at <= ?", userID, begin, end).
		Select("created_at").Find(&submissions).Error
	if err != nil {
		return nil, err
	}
	return submissions, err
}

func (dao *SubmissionDaoImpl) CheckUserIsSubmittedByTime(db *gorm.DB, userID uint, begin time.Time, end time.Time) (bool, error) {
	submission := &repository.Submission{}
	if err := db.Model(submission).Where("user_id = ?", userID).
		Where("created_at >= ? and created_at <= ?", begin, end).Take(submission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (dao *SubmissionDaoImpl) InsertSubmission(db *gorm.DB, submission *repository.Submission) error {
	return db.Create(submission).Error
}
