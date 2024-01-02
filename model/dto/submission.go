package dto

import (
	"funoj-backend/model/repository"
	"funoj-backend/utils"
)

type SubmissionDto struct {
	ID           uint       `json:"id"`
	ProblemName  string     `json:"problemName"`
	Status       int        `json:"status"`
	ErrorMessage string     `json:"errorMessage"`
	CreatedAt    utils.Time `json:"createdAt"`
}

func NewSubmissionDto(submission *repository.Submission) *SubmissionDto {
	return &SubmissionDto{
		ID:           submission.ID,
		Status:       submission.Status,
		ErrorMessage: submission.ErrorMessage,
		CreatedAt:    utils.Time(submission.CreatedAt),
	}
}
