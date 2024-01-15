package services

import (
	"errors"
	conf "funoj-backend/config"
	e "funoj-backend/consts/error"
	"funoj-backend/dao"
	"funoj-backend/db"
	"funoj-backend/file_store"
	"funoj-backend/model/dto"
	"funoj-backend/model/form/request"
	"funoj-backend/model/form/response"
	"funoj-backend/model/repository"
	"funoj-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mime/multipart"
	"path"
	"time"
)

const (
	// ProblemMenuIconPath cos中，题库图标存储位置
	ProblemMenuIconPath = "/icon/problemMenu"
)

// ProblemMenuService 题单管理的service
type ProblemMenuService interface {
	// UploadProblemMenuIcon 上传题单图标
	UploadProblemMenuIcon(file *multipart.FileHeader) (string, *e.Error)
	// ReadProblemMenuIcon 读取题单图标
	ReadProblemMenuIcon(ctx *gin.Context, iconName string)
	// InsertProblemMenu 添加题单
	InsertProblemMenu(problemMenu *repository.ProblemMenu, ctx *gin.Context) (uint, *e.Error)
	// UpdateProblemMenu 更新题单
	UpdateProblemMenu(problemMenu *repository.ProblemMenu) *e.Error
	// DeleteProblemMenu 删除题单
	DeleteProblemMenu(id uint, forceDelete bool) *e.Error
	// GetProblemMenuList 获取题单列表
	GetProblemMenuList(query *request.PageQuery) (*response.PageInfo, *e.Error)
	// GetAllProblemMenu 获取所有的题单列表
	GetAllProblemMenu() ([]*dto.ProblemMenuDtoForList, *e.Error)
	// GetSimpleProblemMenuList 获取简单的题单列表
	GetSimpleProblemMenuList() ([]*dto.ProblemMenuDtoForSimpleList, *e.Error)
	// GetProblemMenuByID 获取题单信息
	GetProblemMenuByID(id uint) (*repository.ProblemMenu, *e.Error)
}

type ProblemMenuServiceImpl struct {
	config         *conf.AppConfig
	problemMenuDao dao.ProblemMenuDao
	problemDao     dao.ProblemDao
	sysUserDao     dao.SysUserDao
}

func NewProblemMenuService(config *conf.AppConfig, pbm dao.ProblemMenuDao, pb dao.ProblemDao, su dao.SysUserDao) ProblemMenuService {
	return &ProblemMenuServiceImpl{
		config:         config,
		problemMenuDao: pbm,
		problemDao:     pb,
		sysUserDao:     su,
	}
}

func (svc *ProblemMenuServiceImpl) UploadProblemMenuIcon(file *multipart.FileHeader) (string, *e.Error) {
	cos := file_store.NewImageCOS(svc.config.COSConfig)
	fileName := file.Filename
	fileName = utils.GetUUID() + "." + path.Base(fileName)
	file2, err := file.Open()
	if err != nil {
		return "", e.ErrBadRequest
	}
	err = cos.SaveFile(path.Join(ProblemMenuIconPath, fileName), file2)
	if err != nil {
		return "", e.ErrServer
	}
	return svc.config.ProUrl + path.Join("/manage/problemBank/icon", fileName), nil
}

func (svc *ProblemMenuServiceImpl) ReadProblemMenuIcon(ctx *gin.Context, iconName string) {
	result := response.NewResult(ctx)
	cos := file_store.NewImageCOS(svc.config.COSConfig)
	bytes, err := cos.ReadFile(path.Join(ProblemMenuIconPath, iconName))
	if err != nil {
		result.Error(e.ErrServer)
		return
	}
	_, _ = ctx.Writer.Write(bytes)
}

func (svc *ProblemMenuServiceImpl) InsertProblemMenu(problemMenu *repository.ProblemMenu, ctx *gin.Context) (uint, *e.Error) {
	// 对设置值的数据设置默认值
	if problemMenu.Name == "" {
		problemMenu.Name = "未命名题单"
	}
	if problemMenu.Description == "" {
		problemMenu.Description = "无描述信息"
	}
	problemMenu.CreatorID = ctx.Keys["user"].(*dto.UserInfo).ID
	err := svc.problemMenuDao.InsertProblemMenu(db.Mysql, problemMenu)
	if err != nil {
		return 0, e.ErrMysql
	}
	return problemMenu.ID, nil
}

func (svc *ProblemMenuServiceImpl) UpdateProblemMenu(problemMenu *repository.ProblemMenu) *e.Error {
	problemMenu.CreatorID = 0
	problemMenu.UpdatedAt = time.Now()
	err := svc.problemMenuDao.UpdateProblemMenu(db.Mysql, problemMenu)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (svc *ProblemMenuServiceImpl) DeleteProblemMenu(id uint, forceDelete bool) *e.Error {
	var err error
	// 非强制删除
	if !forceDelete {
		var count int64
		count, err = svc.problemDao.GetProblemCount(db.Mysql, &request.ProblemForList{
			MenuID: &id,
		})
		if count != 0 {
			return e.NewCustomMsg("题单不为空，请问是否需要强制删除")
		}
		err = svc.problemMenuDao.DeleteProblemMenuByID(db.Mysql, id)
		if err != nil {
			return e.ErrMysql
		}
		return nil
	}

	// 强制删除
	err = svc.problemMenuDao.DeleteProblemMenuByID(db.Mysql, id)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (svc *ProblemMenuServiceImpl) GetProblemMenuList(query *request.PageQuery) (*response.PageInfo, *e.Error) {
	var menuQuery *request.ProblemMenuForList
	if query.Query != nil {
		menuQuery = query.Query.(*request.ProblemMenuForList)
	}
	// 获取题单列表
	menus, err := svc.problemMenuDao.GetProblemMenuList(db.Mysql, query)
	if err != nil {
		return nil, e.ErrMysql
	}
	newProblemMenus := make([]*dto.ProblemMenuDtoForList, len(menus))
	for i := 0; i < len(menus); i++ {
		newProblemMenus[i] = dto.NewProblemMenuDtoForList(menus[i])
		// 读取题单中的题目总数还有作者
		newProblemMenus[i].ProblemCount, err = svc.problemDao.GetProblemCount(db.Mysql, &request.ProblemForList{
			MenuID: &newProblemMenus[i].ID,
		})
		if err != nil {
			return nil, e.ErrMysql
		}
		newProblemMenus[i].CreatorName, err = svc.sysUserDao.GetUserNameByID(db.Mysql, menus[i].CreatorID)
	}
	// 获取所有题单总数目
	var count int64
	count, err = svc.problemMenuDao.GetProblemMenuCount(db.Mysql, menuQuery)
	if err != nil {
		return nil, e.ErrMysql
	}
	pageInfo := &response.PageInfo{
		Total: count,
		Size:  int64(len(newProblemMenus)),
		List:  newProblemMenus,
	}
	return pageInfo, nil
}

func (svc *ProblemMenuServiceImpl) GetAllProblemMenu() ([]*dto.ProblemMenuDtoForList, *e.Error) {
	menus, err := svc.problemMenuDao.GetAllProblemMenu(db.Mysql)
	if err != nil {
		return nil, e.ErrMysql
	}
	answer := make([]*dto.ProblemMenuDtoForList, len(menus))
	for index, menu := range menus {
		answer[index] = dto.NewProblemMenuDtoForList(menu)
	}
	return answer, nil
}

func (svc *ProblemMenuServiceImpl) GetSimpleProblemMenuList() ([]*dto.ProblemMenuDtoForSimpleList, *e.Error) {
	menus, err := svc.problemMenuDao.GetSimpleProblemMenuList(db.Mysql)
	if err != nil {
		return nil, e.ErrMysql
	}
	newMenus := make([]*dto.ProblemMenuDtoForSimpleList, len(menus))
	for i := 0; i < len(menus); i++ {
		newMenus[i] = dto.NewProblemMenuDtoForSimpleList(menus[i])
	}
	return newMenus, nil
}

func (svc *ProblemMenuServiceImpl) GetProblemMenuByID(id uint) (*repository.ProblemMenu, *e.Error) {
	menu, err := svc.problemMenuDao.GetProblemMenuByID(db.Mysql, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, e.ErrProblemNotExist
	}
	if err != nil {
		return nil, e.ErrMysql
	}
	return menu, nil
}
