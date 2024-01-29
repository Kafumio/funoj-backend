package controller

import (
	e "funoj-backend/consts/error"
	"funoj-backend/model/form/request"
	"funoj-backend/model/form/response"
	"funoj-backend/model/repository"
	"funoj-backend/services"
	"funoj-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func (ctl *SysUserController) InsertSysUser(ctx *gin.Context) {
	result := response.NewResult(ctx)
	id, err := ctl.sysUserService.InsertUser(&repository.SysUser{
		LoginName: ctx.PostForm("loginName"),
		UserName:  ctx.PostForm("userName"),
		Password:  ctx.PostForm("password"),
		Email:     ctx.PostForm("email"),
		Phone:     ctx.PostForm("phone"),
	})
	if err != nil {
		result.Error(err)
		return
	}
	result.Success("添加成功", id)
}

func (ctl *SysUserController) UpdateSysUser(ctx *gin.Context) {
	result := response.NewResult(ctx)
	id := utils.AtoiOrDefault(ctx.PostForm("id"), 0)
	user := &repository.SysUser{
		Model: gorm.Model{
			ID: uint(id),
		},
		LoginName: ctx.PostForm("loginName"),
		UserName:  ctx.PostForm("userName"),
		Password:  ctx.PostForm("password"),
		Email:     ctx.PostForm("email"),
		Phone:     ctx.PostForm("phone"),
	}
	if err2 := ctl.sysUserService.UpdateUser(user); err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessMessage("更新成功")
}

func (ctl *SysUserController) DeleteSysUser(ctx *gin.Context) {
	result := response.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	err2 := ctl.sysUserService.DeleteUser(uint(id))
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessMessage("删除成功")
}

func (ctl *SysUserController) GetSysUserList(ctx *gin.Context) {
	result := response.NewResult(ctx)
	pageQuery, err := utils.GetPageQueryByQuery(ctx)
	if err != nil {
		result.Error(err)
		return
	}
	user := &repository.SysUser{
		UserName:     ctx.Query("userName"),
		LoginName:    ctx.Query("loginName"),
		Email:        ctx.Query("email"),
		Phone:        ctx.Query("phone"),
		Introduction: ctx.Query("introduction"),
	}
	sexStr := ctx.Query("sex")
	if sexStr == "1" {
		user.Gender = 1
	} else if sexStr == "2" {
		user.Gender = 2
	}
	pageQuery.Query = user
	pageInfo, err2 := ctl.sysUserService.GetUserList(pageQuery)
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData(pageInfo)
}

func (ctl *SysUserController) UpdateUserRoles(ctx *gin.Context) {
	result := response.NewResult(ctx)
	var req request.UpdateUserRolesRequest
	if err := ctx.BindJSON(&req); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	if err2 := ctl.sysUserService.UpdateUserRoles(req.UserID, req.RoleIDs); err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessMessage("更新成功")
}

func (ctl *SysUserController) GetRoleIDsByUserID(ctx *gin.Context) {
	result := response.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	roleIDs, err2 := ctl.sysUserService.GetRoleIDsByUserID(uint(id))
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData(roleIDs)
}

func (ctl *SysUserController) GetAllSimpleRole(ctx *gin.Context) {
	result := response.NewResult(ctx)
	roles, err := ctl.sysUserService.GetAllSimpleRole()
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(roles)
}
