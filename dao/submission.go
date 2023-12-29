package dao

type SubmissionDao interface {
}

type SubmissionDaoImpl struct {
}

func NewSubmissionDao() SubmissionDao {
	return &SubmissionDaoImpl{}
}
