package svc

import (
	"mymall/pkg/config"
	"mymall/services/merchant-service/internal/repository"

	"gorm.io/gorm"
)

// ServiceContext 全局依赖（go-zero 惯例）
type ServiceContext struct {
	Config *config.Config
	DB     *gorm.DB
	Repo   *repository.MerchantRepository
}

func NewServiceContext(cfg *config.Config, db *gorm.DB) *ServiceContext {
	return &ServiceContext{
		Config: cfg,
		DB:     db,
		Repo:   repository.NewMerchantRepository(db),
	}
}
