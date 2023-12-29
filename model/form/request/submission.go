package request

type SubmissionForList struct {
	ID        uint `json:"id"`
	UserID    uint `json:"userID"`
	ProblemID uint `json:"problemID"`
}
