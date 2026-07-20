package svc

import (
	"mymall/pkg/config"
	"mymall/services/merchant-service/internal/client/userhttp"
	"mymall/services/merchant-service/internal/repository"

	"gorm.io/gorm"
)

// ServiceContext 全局依赖（go-zero 惯例）
type ServiceContext struct {
	Config         *config.Config
	DB             *gorm.DB
	Repo           *repository.MerchantRepository
	PointsProducts *repository.PointsProductRepository
	PointsOrders   *repository.PointsOrderRepository
	UserHTTP       *userhttp.Client
}

func NewServiceContext(cfg *config.Config, db *gorm.DB) *ServiceContext {
	return &ServiceContext{
		Config:         cfg,
		DB:             db,
		Repo:           repository.NewMerchantRepository(db),
		PointsProducts: repository.NewPointsProductRepository(db),
		PointsOrders:   repository.NewPointsOrderRepository(db),
		UserHTTP:       userhttp.New(""),
	}
}
