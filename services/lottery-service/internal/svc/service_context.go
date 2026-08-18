package svc

import (
	"mymall/pkg/health"
	"mymall/services/lottery-service/internal/client/userhttp"
	"mymall/services/lottery-service/internal/config"
	"mymall/services/lottery-service/internal/middleware"
	"mymall/services/lottery-service/internal/repository"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config   *config.Config
	Conn     sqlx.SqlConn
	Redis    *redis.Client
	Repo     *repository.LotteryRepository
	UserHTTP *userhttp.Client
	Health   *health.Registry

	RequestID            rest.Middleware
	GatewayIdentity      rest.Middleware
	RequirePlatformAdmin rest.Middleware
}

func NewServiceContext(cfg *config.Config, conn sqlx.SqlConn, rdb *redis.Client, healthReg *health.Registry) *ServiceContext {
	return &ServiceContext{
		Config:               cfg,
		Conn:                 conn,
		Redis:                rdb,
		Repo:                 repository.NewLotteryRepository(conn),
		UserHTTP:             userhttp.New(cfg.UserHTTP, cfg.UserHTTPTimeout),
		Health:               healthReg,
		RequestID:            middleware.NewRequestIDMiddleware().Handle,
		GatewayIdentity:      middleware.NewGatewayIdentityMiddleware().Handle,
		RequirePlatformAdmin: middleware.NewRequirePlatformAdminMiddleware().Handle,
	}
}
