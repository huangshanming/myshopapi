package svc

import (
	"mymall/pkg/health"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ServiceContext goctl 惯例依赖容器。
type ServiceContext struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Health *health.Registry
}

func NewServiceContext(db *gorm.DB, rdb *redis.Client, healthReg *health.Registry) *ServiceContext {
	return &ServiceContext{
		DB:     db,
		Redis:  rdb,
		Health: healthReg,
	}
}
