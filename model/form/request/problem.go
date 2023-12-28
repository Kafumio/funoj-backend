package request

type ProblemForList struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Number     string `json:"number"`
	MenuID     *uint  `json:"menuID"`
	Difficulty int    `json:"difficulty"`
	Enable     int    `json:"enable"`
}
