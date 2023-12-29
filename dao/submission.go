package dao

import (
	"errors"
	"funoj-backend/model/form/request"
	"funoj-backend/model/repository"
	"gorm.io/gorm"
	"time"
)

type SubmissionDao interface {
	GetLastSubmission(db *gorm.DB, userID uint, problemID uint) (*repository.Submission, error)
	GetSubmissionList(db *gorm.DB, pageQuery *request.PageQuery) ([]*repository.Submission, error)
	GetSubmissionCount(db *gorm.DB, submission *request.SubmissionForList) (int64, error)
	GetUserSimpleSubmissionsByTime(db *gorm.DB, userID uint, begin time.Time, end time.Time) ([]*repository.Submission, error)
	CheckUserIsSubmittedByTime(db *gorm.DB, userID uint, begin time.Time, end time.Time) (bool, error)
	InsertSubmission(db *gorm.DB, submission *repository.Submission) error
}

type SubmissionDaoImpl struct {
}

func NewSubmissionDao() SubmissionDao {
	return &SubmissionDaoImpl{}
}

func (s *SubmissionDaoImpl) GetLastSubmission(db *gorm.DB, userID uint, problemID uint) (*repository.Submission, error) {
	var submission *repository.Submission
	err := db.Where("user_id = ? and problem_id = ?", userID, problemID).Last(submission).Error
	return submission, err
}

func (s *SubmissionDaoImpl) GetSubmissionList(db *gorm.DB, pageQuery *request.PageQuery) ([]*repository.Submission, error) {
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

func (s *SubmissionDaoImpl) GetSubmissionCount(db *gorm.DB, submission *request.SubmissionForList) (int64, error) {
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

func (s *SubmissionDaoImpl) GetUserSimpleSubmissionsByTime(db *gorm.DB, userID uint, begin time.Time, end time.Time) ([]*repository.Submission, error) {
	var submissions []*repository.Submission
	err := db.Where("user_id = ? and created_at >= ? and created_at <= ?", userID, begin, end).
		Select("created_at").Find(&submissions).Error
	if err != nil {
		return nil, err
	}
	return submissions, err
}

func (s *SubmissionDaoImpl) CheckUserIsSubmittedByTime(db *gorm.DB, userID uint, begin time.Time, end time.Time) (bool, error) {
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

func (s *SubmissionDaoImpl) InsertSubmission(db *gorm.DB, submission *repository.Submission) error {
	return db.Create(submission).Error
}
