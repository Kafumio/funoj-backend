package services

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAccountService,
	NewAuthService,
	NewProblemMenuService,
	NewProblemService,
	NewProblemCaseService,
	NewSubmissionService,
	NewSysPermissionService,
	NewSysRoleService,
	NewSysUserService,
)
