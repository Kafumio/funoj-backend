package db

import (
	"fmt"
	"funoj-backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitMysql
//
//	@Description: 初始化mysql
//	@param cfg
//	@return error
func InitMysql(cfg *config.MySqlConfig) error {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	var err error
	Mysql, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}
