package request

type ProblemCaseForList struct {
	ID        uint   `json:"ID"`
	ProblemID uint   `json:"problemID"`
	CaseName  string `json:"caseName"`
}
