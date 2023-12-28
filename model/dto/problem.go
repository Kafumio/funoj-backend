package dto

// ProblemDtoForGet 获取题目详细信息
type ProblemDtoForGet struct {
	ID          uint   `json:"id"`
	MenuID      *uint  `json:"menuID"`
	Name        string `json:"name"`
	Number      string `json:"number"`
	Description string `json:"description"`
	Title       string `json:"title"`
	Path        string `json:"path"`
	Difficulty  int    `json:"difficulty"`
	// 支持的语言用,分割
	Languages string `json:"languages"`
	Enable    int    `json:"enable"`
}
