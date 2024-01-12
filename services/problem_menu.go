package services

import (
	conf "funoj-backend/config"
	e "funoj-backend/consts/error"
	"funoj-backend/dao"
	"mime/multipart"
)

const (
	// ProblemMenuIconPath cos中，题库图标存储位置
	ProblemMenuIconPath = "/icon/problemMenu"
)

type ProblemMenuService interface {
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

}
