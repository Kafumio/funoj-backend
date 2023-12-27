package request

type Problem struct {
	Number     string `json:"number"`
	Name       string `json:"name"`
	Difficulty int    `json:"difficulty"`
	// 0空值，1启用，-1停用
	Enable int   `gorm:"column:enable" json:"enable"`
	MenuID *uint `json:"menuID"`
}
