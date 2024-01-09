package services

import (
	"errors"
	conf "funoj-backend/config"
	e "funoj-backend/consts/error"
	"funoj-backend/dao"
	"funoj-backend/db"
	"funoj-backend/model/dto"
	"funoj-backend/model/form/request"
	"funoj-backend/model/form/response"
	"funoj-backend/model/repository"
	"funoj-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type ProblemService interface {
	// CheckProblemNumber 检测题目编码
	CheckProblemNumber(problemCode string) (bool, *e.Error)
	// InsertProblem 添加题目
	InsertProblem(ctx *gin.Context, problem *repository.Problem) (uint, *e.Error)
	// UpdateProblem 更新题目
	UpdateProblem(Problem *repository.Problem) *e.Error
	// DeleteProblem 删除题目
	DeleteProblem(id uint) *e.Error
	// GetProblemList 获取题目列表
	GetProblemList(query *request.PageQuery) (*response.PageInfo, *e.Error)
	// GetUserProblemList 用户获取题目列表
	GetUserProblemList(ctx *gin.Context, query *request.PageQuery) (*response.PageInfo, *e.Error)
	// GetProblemByID 获取题目信息
	GetProblemByID(id uint) (*dto.ProblemDtoForGet, *e.Error)
	// GetProblemByNumber 根据题目编号获取题目信息
	GetProblemByNumber(number string) (*dto.ProblemDtoForGet, *e.Error)
	// GetProblemTemplateCode 获取题目的模板代码
	GetProblemTemplateCode(problemID uint, language string) (string, *e.Error)
	// UpdateProblemEnable 设置题目可用
	UpdateProblemEnable(id uint, enable int) *e.Error
}

type ProblemServiceImpl struct {
	config            *conf.AppConfig
	problemDao        dao.ProblemDao
	problemCaseDao    dao.ProblemCaseDao
	problemAttemptDao dao.ProblemAttemptDao
}

func NewProblemService(config *conf.AppConfig, problemDao dao.ProblemDao, problemCaseDao dao.ProblemCaseDao, problemAttempt dao.ProblemAttemptDao) ProblemService {
	return &ProblemServiceImpl{
		config:            config,
		problemDao:        problemDao,
		problemCaseDao:    problemCaseDao,
		problemAttemptDao: problemAttempt,
	}
}

func (svc *ProblemServiceImpl) CheckProblemNumber(problemCode string) (bool, *e.Error) {
	b, err := svc.problemDao.CheckProblemNumberExists(db.Mysql, problemCode)
	if err != nil {
		return !b, e.ErrProblemCodeCheckFailed
	}
	return !b, nil
}

func (svc *ProblemServiceImpl) InsertProblem(ctx *gin.Context, problem *repository.Problem) (uint, *e.Error) {
	problem.CreatorID = ctx.Keys["user"].(*dto.UserInfo).ID
	// 对设置值的数据设置默认值
	if problem.Name == "" {
		problem.Name = "未命名题目"
	}
	if problem.Title == "" {
		problem.Title = "标题信息"
	}
	if problem.Description == "" {
		problemDescription, err := os.ReadFile(svc.config.FilePathConfig.ProblemDescriptionTemplate)
		if err != nil {
			return 0, e.ErrProblemInsertFailed
		}
		problem.Description = string(problemDescription)
	}
	if problem.Number == "" {
		problem.Number = "未命名编号" + utils.GetGenerateUniqueCode()
	}
	// 检测编号是否重复
	if problem.Number != "" {
		b, checkError := svc.problemDao.CheckProblemNumberExists(db.Mysql, problem.Number)
		if checkError != nil {
			return 0, e.ErrMysql
		}
		if b {
			return 0, e.ErrProblemCodeIsExist
		}
	}
	// 题目难度不在范围，那么都设置为1
	if problem.Difficulty > 5 || problem.Difficulty < 1 {
		problem.Difficulty = 1
	}
	problem.Enable = -1
	// 添加
	err := svc.problemDao.InsertProblem(db.Mysql, problem)
	if err != nil {
		return 0, e.ErrMysql
	}
	return problem.ID, nil
}

func (svc *ProblemServiceImpl) UpdateProblem(problem *repository.Problem) *e.Error {
	problem.UpdatedAt = time.Now()
	if err := svc.problemDao.UpdateProblem(db.Mysql, problem); err != nil {
		log.Println(err)
		return e.ErrProblemUpdateFailed
	}
	return nil
}

func (svc *ProblemServiceImpl) DeleteProblem(id uint) *e.Error {
	// 读取Problem
	problem, err := svc.problemDao.GetProblemByID(db.Mysql, id)
	if err != nil {
		return e.ErrMysql
	}
	if problem == nil || problem.Number == "" {
		return e.ErrProblemNotExist
	}
	// 删除用例
	if err = svc.problemCaseDao.DeleteProblemCaseByProblemID(db.Mysql, id); err != nil {
		return e.ErrMysql
	}
	// 删除题目
	if err = svc.problemDao.DeleteProblemByID(db.Mysql, id); err != nil {
		return e.ErrMysql
	}
	return nil
}

func (svc *ProblemServiceImpl) GetProblemList(query *request.PageQuery) (*response.PageInfo, *e.Error) {
	var problemQuery *request.ProblemForList
	if query.Query != nil {
		problemQuery = query.Query.(*request.ProblemForList)
	}
	// 获取题目列表
	problems, err := svc.problemDao.GetProblemList(db.Mysql, query)
	if err != nil {
		return nil, e.ErrMysql
	}
	newProblems := make([]*dto.ProblemDtoForList, len(problems))
	for i := 0; i < len(problems); i++ {
		newProblems[i] = dto.NewProblemDtoForList(problems[i])
	}
	// 获取所有题目总数目
	var count int64
	count, err = svc.problemDao.GetProblemCount(db.Mysql, problemQuery)
	if err != nil {
		return nil, e.ErrMysql
	}
	pageInfo := &response.PageInfo{
		Total: count,
		Size:  int64(len(newProblems)),
		List:  newProblems,
	}
	return pageInfo, nil
}

func (svc *ProblemServiceImpl) GetUserProblemList(ctx *gin.Context, query *request.PageQuery) (*response.PageInfo, *e.Error) {
	userId := ctx.Keys["user"].(*dto.UserInfo).ID
	if query.Query != nil {
		query.Query.(*request.ProblemForList).Enable = 1
	} else {
		query.Query = &request.ProblemForList{
			Enable: 1,
		}
	}
	// 获取题目列表
	problems, err := svc.problemDao.GetProblemList(db.Mysql, query)
	if err != nil {
		return nil, e.ErrMysql
	}
	newProblems := make([]*dto.ProblemDtoForUserList, len(problems))
	for i := 0; i < len(problems); i++ {
		newProblems[i] = dto.NewProblemDtoForUserList(problems[i])
		// 读取题目完成情况
		var status int
		status, err = svc.problemAttemptDao.GetProblemAttemptStatus(db.Mysql, userId, problems[i].ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrProblemListFailed
		}
		newProblems[i].Status = status
	}
	// 获取所有题目总数目
	var count int64
	count, err = svc.problemDao.GetProblemCount(db.Mysql, query.Query.(*request.ProblemForList))
	if err != nil {
		return nil, e.ErrMysql
	}
	pageInfo := &response.PageInfo{
		Total: count,
		Size:  int64(len(newProblems)),
		List:  newProblems,
	}
	return pageInfo, nil
}

func (svc *ProblemServiceImpl) GetProblemByID(id uint) (*dto.ProblemDtoForGet, *e.Error) {
	problem, err := svc.problemDao.GetProblemByID(db.Mysql, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, e.ErrProblemNotExist
	}
	if err != nil {
		return nil, e.ErrMysql
	}
	return dto.NewProblemDtoForGet(problem), nil
}

func (svc *ProblemServiceImpl) GetProblemByNumber(number string) (*dto.ProblemDtoForGet, *e.Error) {
	problem, err := svc.problemDao.GetProblemByNumber(db.Mysql, number)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, e.ErrProblemNotExist
	}
	if err != nil {
		return nil, e.ErrMysql
	}
	return dto.NewProblemDtoForGet(problem), nil
}

func (svc *ProblemServiceImpl) GetProblemTemplateCode(problemID uint, language string) (string, *e.Error) {
	// 读取acm模板
	code, err := utils.GetAcmCodeTemplate(language)
	if err != nil {
		return "", e.ErrProblemGetFailed
	}
	return code, nil
}

// todo: 是否要加事务
func (svc *ProblemServiceImpl) UpdateProblemEnable(id uint, enable int) *e.Error {
	if err := svc.problemDao.SetProblemEnable(db.Mysql, id, enable); err != nil {
		return e.ErrMysql
	}
	return nil
}
