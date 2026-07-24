package svc

import (
	"mymall/pkg/health"
	"mymall/services/merchant-service/internal/client/tencentmap"
	"mymall/services/merchant-service/internal/config"
	"mymall/services/merchant-service/internal/middleware"
	"mymall/services/merchant-service/internal/repository"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

// ServiceContext 全局依赖（go-zero 惯例）
type ServiceContext struct {
	Config     *config.Config
	Conn       sqlx.SqlConn
	Repo       *repository.MerchantRepository
	Health     *health.Registry
	TencentMap *tencentmap.Client

	RequestID            rest.Middleware
	GatewayIdentity      rest.Middleware
	RequireMerchantOwner rest.Middleware
	RequirePlatformAdmin rest.Middleware
}

func NewServiceContext(cfg *config.Config, conn sqlx.SqlConn, healthReg *health.Registry) *ServiceContext {
	return &ServiceContext{
		Config:               cfg,
		Conn:                 conn,
		Repo:                 repository.NewMerchantRepository(conn),
		Health:               healthReg,
		TencentMap:           tencentmap.New(cfg.TencentMap.Key, cfg.TencentMap.BaseURL),
		RequestID:            middleware.NewRequestIDMiddleware().Handle,
		GatewayIdentity:      middleware.NewGatewayIdentityMiddleware().Handle,
		RequireMerchantOwner: middleware.NewRequireMerchantOwnerMiddleware().Handle,
		RequirePlatformAdmin: middleware.NewRequirePlatformAdminMiddleware().Handle,
	}
}
