package dto

import (
	"funoj-backend/model/repository"
	"funoj-backend/utils"
)

// ProblemMenuDtoForList 获取题目列表
type ProblemMenuDtoForList struct {
	ID           uint       `json:"id"`
	Icon         string     `json:"icon"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	CreatedAt    utils.Time `json:"createdAt"`
	UpdatedAt    utils.Time `json:"updatedAt"`
	CreatorName  string     `json:"creatorName"`
	ProblemCount int64      `json:"problemCount"`
}

func NewProblemMenuDtoForList(menu *repository.ProblemMenu) *ProblemMenuDtoForList {
	return &ProblemMenuDtoForList{
		ID:          menu.ID,
		Icon:        menu.Icon,
		Name:        menu.Name,
		Description: menu.Description,
		CreatedAt:   utils.Time(menu.CreatedAt),
		UpdatedAt:   utils.Time(menu.UpdatedAt),
	}
}

// ProblemMenuDtoForSimpleList 获取简单题目列表
type ProblemMenuDtoForSimpleList struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewProblemMenuDtoForSimpleList(menu *repository.ProblemMenu) *ProblemMenuDtoForSimpleList {
	response := &ProblemMenuDtoForSimpleList{
		ID:   menu.ID,
		Name: menu.Name,
	}
	return response
}
