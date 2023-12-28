package request

// SysUserForList 获取用户列表请求结构
type SysUserForList struct {
	ID        uint   `json:"id"`
	LoginName string `json:"loginName"`
	UserName  string `json:"username"`
	Gender    int    `json:"gender"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
