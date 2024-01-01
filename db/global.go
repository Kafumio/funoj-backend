package db

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Mysql *gorm.DB
	Redis *redis.Client
)
