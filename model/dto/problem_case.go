package dto

import (
	"funoj-backend/model/repository"
	"funoj-backend/utils"
)

type ProblemCaseDtoForList struct {
	ID        uint       `json:"id"`
	CaseName  string     `json:"caseName"`
	Input     string     `json:"input"`
	Output    string     `json:"output"`
	CreatedAt utils.Time `json:"createdAt"`
}

func NewProblemCaseDtoForList(problemCase *repository.ProblemCase) *ProblemCaseDtoForList {
	return &ProblemCaseDtoForList{
		ID:        problemCase.ID,
		CaseName:  problemCase.CaseName,
		Input:     problemCase.Input,
		Output:    problemCase.Output,
		CreatedAt: utils.Time(problemCase.CreatedAt),
	}
}

type ProblemCaseDto struct {
	ID       uint   `json:"id"`
	CaseName string `json:"caseName"`
	Input    string `json:"input"`
	Output   string `json:"output"`
}

func NewProblemCaseDto(problemCase *repository.ProblemCase) *ProblemCaseDto {
	return &ProblemCaseDto{
		ID:       problemCase.ID,
		CaseName: problemCase.CaseName,
		Input:    problemCase.Input,
		Output:   problemCase.Output,
	}
}
