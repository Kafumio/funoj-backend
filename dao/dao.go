package dao

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewProblemAttemptDao,
	NewProblemMenuDao,
	NewProblemDao,
	NewProblemCaseDao,
	NewSubmissionDao,
	NewSysPermissionDao,
	NewSysRoleDao,
	NewSysUserDao,
)
