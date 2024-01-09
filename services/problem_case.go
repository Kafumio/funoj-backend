package services

type ProblemCaseService interface {
}

type ProblemCaseServiceImpl struct {
}

func NewProblemCaseService() ProblemCaseService {
	return &ProblemCaseServiceImpl{}
}

func (p *ProblemCaseServiceImpl) name() {

}
