package services

import (
	e "funoj-backend/consts/error"
	"funoj-backend/dao"
	"funoj-backend/db"
	"funoj-backend/model/dto"
	"funoj-backend/model/form/request"
	"funoj-backend/model/form/response"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

type SubmissionService interface {
	// GetActivityMap 获取活动图
	GetActivityMap(ctx *gin.Context, year int) ([]*dto.ActivityItem, *e.Error)
	// GetActivityYear 获取用户有活动的年份
	GetActivityYear(ctx *gin.Context) ([]string, *e.Error)
	// GetUserSubmissionList 获取用户
	GetUserSubmissionList(ctx *gin.Context, pageQuery *request.PageQuery) (*response.PageInfo, *e.Error)
}

type SubmissionServiceImpl struct {
	submissionDao dao.SubmissionDao
	problemDao    dao.ProblemDao
}

func NewSubmissionDao(submissionDao dao.SubmissionDao, problemDao dao.ProblemDao) SubmissionService {
	return &SubmissionServiceImpl{
		submissionDao: submissionDao,
		problemDao:    problemDao,
	}
}

func (s *SubmissionServiceImpl) GetActivityMap(ctx *gin.Context, year int) ([]*dto.ActivityItem, *e.Error) {
	user := ctx.Keys["user"].(*dto.UserInfo)
	var startDate time.Time
	var endDate time.Time
	// 如果year == 0，获取以今天截至的一年的数据
	if year == 0 {
		endDate = time.Now()
		startDate = time.Date(endDate.Year()-1, endDate.Month()+1, endDate.Day(),
			0, 0, 0, 0, time.Local)
	} else {
		startDate, endDate = s.getYearRange(year)
	}
	submissions, err := s.submissionDao.GetUserSimpleSubmissionsByTime(db.Mysql, user.ID, startDate, endDate)
	if err != nil {
		return nil, e.ErrMysql
	}
	// 构建活动数据
	m := make(map[string]int, 366)
	for i := 0; i < len(submissions); i++ {
		date := submissions[i].CreatedAt.Format("2006-01-02")
		m[date]++
	}
	answer := make([]*dto.ActivityItem, len(m))
	i := 0
	for k, v := range m {
		answer[i] = &dto.ActivityItem{
			Date:  k,
			Count: v,
		}
		i++
	}
	return answer, nil
}

func (s *SubmissionServiceImpl) GetActivityYear(ctx *gin.Context) ([]string, *e.Error) {
	answer := []string{}
	user := ctx.Keys["user"].(*dto.UserInfo)
	beginYear := 2022
	currentYear := time.Now().Year()
	for i := beginYear; i <= currentYear; i++ {
		beginDate, endDate := s.getYearRange(i)
		b, err := s.submissionDao.CheckUserIsSubmittedByTime(db.Mysql, user.ID, beginDate, endDate)
		if err != nil {
			return nil, e.ErrMysql
		}
		if b {
			answer = append(answer, strconv.Itoa(i))
		}
	}
	return answer, nil
}

func (s *SubmissionServiceImpl) GetUserSubmissionList(ctx *gin.Context, pageQuery *request.PageQuery) (*response.PageInfo, *e.Error) {
	user := ctx.Keys["user"].(*dto.UserInfo)
	submissionReq := &request.SubmissionForList{
		UserID: user.ID,
	}
	pageQuery.Query = submissionReq
	submissions, err := s.submissionDao.GetSubmissionList(db.Mysql, pageQuery)
	if err != nil {
		log.Println(err)
		return nil, e.ErrMysql
	}
	submissionList := make([]*dto.SubmissionDto, len(submissions))
	for i := 0; i < len(submissions); i++ {
		submissionList[i] = dto.NewSubmissionDto(submissions[i])
		name, err := s.problemDao.GetProblemNameByID(db.Mysql, submissions[i].ProblemID)
		if err != nil {
			return nil, e.ErrMysql
		}
		submissionList[i].ProblemName = name
	}
	count, err2 := s.submissionDao.GetSubmissionCount(db.Mysql, submissionReq)
	if err2 != nil {
		log.Println(err2)
		return nil, e.ErrMysql
	}
	return &response.PageInfo{
		Total: count,
		Size:  int64(len(submissionList)),
		List:  submissionList,
	}, nil
}

func (s *SubmissionServiceImpl) getYearRange(year int) (time.Time, time.Time) {
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 999999999, time.Local)
	return startDate, endDate
}
