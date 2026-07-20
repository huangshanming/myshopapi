package svc

import (
	"mymall/pkg/config"
	"mymall/pkg/health"
	"mymall/services/merchant-service/internal/client/userhttp"
	"mymall/services/merchant-service/internal/middleware"
	"mymall/services/merchant-service/internal/repository"

	"github.com/zeromicro/go-zero/rest"
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
	Health         *health.Registry

	RequestID            rest.Middleware
	GatewayIdentity      rest.Middleware
	RequireMerchantOwner rest.Middleware
	RequirePlatformAdmin rest.Middleware
}

func NewServiceContext(cfg *config.Config, db *gorm.DB, healthReg *health.Registry) *ServiceContext {
	return &ServiceContext{
		Config:               cfg,
		DB:                   db,
		Repo:                 repository.NewMerchantRepository(db),
		PointsProducts:       repository.NewPointsProductRepository(db),
		PointsOrders:         repository.NewPointsOrderRepository(db),
		UserHTTP:             userhttp.New(""),
		Health:               healthReg,
		RequestID:            middleware.NewRequestIDMiddleware().Handle,
		GatewayIdentity:      middleware.NewGatewayIdentityMiddleware().Handle,
		RequireMerchantOwner: middleware.NewRequireMerchantOwnerMiddleware().Handle,
		RequirePlatformAdmin: middleware.NewRequirePlatformAdminMiddleware().Handle,
	}
}
