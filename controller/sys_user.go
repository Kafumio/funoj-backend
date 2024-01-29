package controller

import (
	"funoj-backend/model/form/response"
	"funoj-backend/services"
	"funoj-backend/utils"
	"github.com/gin-gonic/gin"
)

type SysUserController struct {
	sysUserService services.SysUserService
}

func NewSysUserController(userService services.SysUserService) *SysUserController {
	return &SysUserController{
		sysUserService: userService,
	}
}

func (ctl *SysUserController) GetUserByID(ctx *gin.Context) {
	result := response.NewResult(ctx)
	userID := utils.GetIntParamOrDefault(ctx, "id", 0)
	user, err2 := ctl.sysUserService.GetUserByID(uint(userID))
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData(user)
}
