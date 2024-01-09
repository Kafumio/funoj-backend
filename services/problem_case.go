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
	"gorm.io/gorm"
	"log"
	"strconv"
	"unicode"
)

// ProblemCaseService 题目用例管理
type ProblemCaseService interface {
	// GetProblemCaseList 获取用例列表
	GetProblemCaseList(query *request.PageQuery) (*response.PageInfo, *e.Error)
	// GetProblemCaseByID 通过id获取题目用例
	GetProblemCaseByID(id uint) (*dto.ProblemCaseDto, *e.Error)
	// DeleteProblemCaseByID 通过id删除题目用例
	DeleteProblemCaseByID(id uint) *e.Error
	// InsertProblemCase 添加题目用例
	InsertProblemCase(problemCase *repository.ProblemCase) (uint, *e.Error)
	// UpdateProblemCase 更新题目用例
	UpdateProblemCase(problemCase *repository.ProblemCase) *e.Error
	// CheckProblemCaseName 检测用例名称是否重复
	CheckProblemCaseName(id uint, name string, problemID uint) (bool, *e.Error)
	// GenerateNewProblemCaseName 生成一个题目唯一用例名称，递增
	GenerateNewProblemCaseName(problemID uint) (string, *e.Error)
}

type ProblemCaseServiceImpl struct {
	config         *conf.AppConfig
	problemCaseDao dao.ProblemCaseDao
	problemDao     dao.ProblemDao
}

func NewProblemCaseService(config *conf.AppConfig, pcd dao.ProblemCaseDao, pd dao.ProblemDao) ProblemCaseService {
	return &ProblemCaseServiceImpl{
		config:         config,
		problemCaseDao: pcd,
		problemDao:     pd,
	}
}

func (svc *ProblemCaseServiceImpl) GetProblemCaseByID(id uint) (*dto.ProblemCaseDto, *e.Error) {
	problemCase, err := svc.problemCaseDao.GetProblemCaseByID(db.Mysql, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, e.ErrProblemNotExist
	}
	if err != nil {
		log.Println("Error while getting problem case name:", err)
		return nil, e.ErrMysql
	}
	return dto.NewProblemCaseDto(problemCase), nil
}

func (svc *ProblemCaseServiceImpl) GetProblemCaseList(query *request.PageQuery) (*response.PageInfo, *e.Error) {
	var problemCase *request.ProblemCaseForList
	if query.Query != nil {
		problemCase = query.Query.(*request.ProblemCaseForList)
	}
	// 获取用例列表
	cases, err := svc.problemCaseDao.GetProblemCaseList(db.Mysql, query)
	if err != nil {
		log.Println("Error while getting problem case list:", err)
		return nil, e.ErrMysql
	}
	newCases := make([]*dto.ProblemCaseDtoForList, len(cases))
	for i := 0; i < len(cases); i++ {
		newCases[i] = dto.NewProblemCaseDtoForList(cases[i])
	}
	// 获取所有用例总数目
	var count int64
	count, err = svc.problemCaseDao.GetProblemCaseCount(db.Mysql, problemCase)
	if err != nil {
		log.Println("Error while getting problem case list:", err)
		return nil, e.ErrMysql
	}
	pageInfo := &response.PageInfo{
		Total: count,
		Size:  int64(len(newCases)),
		List:  newCases,
	}
	return pageInfo, nil
}

func (svc *ProblemCaseServiceImpl) DeleteProblemCaseByID(id uint) *e.Error {
	err := svc.problemCaseDao.DeleteProblemCaseByID(db.Mysql, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return e.ErrProblemNotExist
	}
	if err != nil {
		log.Println("Error while deleting problem case:", err)
		return e.ErrMysql
	}
	return nil
}

func (svc *ProblemCaseServiceImpl) InsertProblemCase(problemCase *repository.ProblemCase) (uint, *e.Error) {
	err := svc.problemCaseDao.InsertProblemCase(db.Mysql, problemCase)
	if err != nil {
		log.Println("Error while inserting problem case:", err)
		return 0, e.ErrMysql
	}
	return problemCase.ID, nil
}

func (svc *ProblemCaseServiceImpl) UpdateProblemCase(problemCase *repository.ProblemCase) *e.Error {
	err := svc.problemCaseDao.UpdateProblemCase(db.Mysql, problemCase)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return e.ErrProblemNotExist
	}
	if err != nil {
		log.Println("Error while updating problem case:", err)
		return e.ErrMysql
	}
	return nil
}

func (svc *ProblemCaseServiceImpl) CheckProblemCaseName(id uint, caseName string, problemID uint) (bool, *e.Error) {
	l, err := svc.problemCaseDao.GetProblemCaseList(db.Mysql, &request.PageQuery{
		Page:     1,
		PageSize: 1,
		Query: &request.ProblemCaseForList{
			CaseName:  caseName,
			ProblemID: problemID,
		},
	})
	if err != nil {
		log.Println("Error while checking problem case name:", err)
		return false, e.ErrMysql
	}
	return len(l) == 0 || len(l) == 1 && l[0].ID == id, nil
}

func (svc *ProblemCaseServiceImpl) GenerateNewProblemCaseName(problemID uint) (string, *e.Error) {
	// 获取与给定问题ID相关的问题用例列表
	problemCases, err := svc.problemCaseDao.GetProblemCaseList(db.Mysql, &request.PageQuery{
		Page:         1,
		PageSize:     1,
		SortProperty: "name",
		SortRule:     "desc",
		Query: &request.ProblemCaseForList{
			ProblemID: problemID,
		},
	})
	if err != nil {
		// 记录错误并返回
		log.Println("Error while getting problem case list:", err)
		return "", e.ErrMysql
	}
	// 如果没有找到问题用例，直接返回 "1"
	if len(problemCases) == 0 {
		return "1", nil
	}

	// 获取最新的问题用例
	latestCase := problemCases[0]

	// 寻找最后一个非数字字符的索引
	i := len(latestCase.CaseName) - 1
	for i >= 0 && unicode.IsDigit(rune(latestCase.CaseName[i])) {
		i--
	}

	// 截取数字部分并转换为整数
	numericPart := latestCase.CaseName[i+1:]
	num, err := strconv.Atoi(numericPart)
	if err != nil {
		// 记录错误并返回
		log.Println("Error converting numeric part to integer:", err)
		return "", e.ErrUnknown
	}

	// 递增数字部分并生成新名称
	num++
	newName := latestCase.CaseName[:i+1] + strconv.Itoa(num)
	return newName, nil
}
