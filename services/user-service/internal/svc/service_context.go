package svc

import (
	"time"

	"mymall/pkg/health"
	"mymall/pkg/jwt"
	"mymall/services/user-service/internal/client/dingdanxia"
	"mymall/services/user-service/internal/client/jutuike"
	"mymall/services/user-service/internal/config"
	"mymall/services/user-service/internal/middleware"
	"mymall/services/user-service/internal/repository"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config         *config.Config
	Conn           sqlx.SqlConn
	Repo           *repository.UserRepository
	RBAC           *repository.RBACRepository
	Tasks          *repository.TaskRepository
	PointsProducts *repository.PointsProductRepository
	PointsOrders   *repository.PointsOrderRepository
	Jutuike        *jutuike.Client
	Dingdanxia     *dingdanxia.Client
	JWT            jwt.Config
	Health         *health.Registry

	RequestID            rest.Middleware
	GatewayIdentity      rest.Middleware
	RequirePlatformAdmin rest.Middleware
	UserFlexibleAuth     rest.Middleware
}

func NewServiceContext(cfg *config.Config, conn sqlx.SqlConn) *ServiceContext {
	flex := middleware.NewUserFlexibleAuthMiddlewareWithSecret(cfg.Auth.AccessSecret)
	jtkTimeout := time.Duration(cfg.Jutuike.Timeout) * time.Second
	if jtkTimeout <= 0 {
		jtkTimeout = 8 * time.Second
	}
	ddxTimeout := time.Duration(cfg.Dingdanxia.Timeout) * time.Second
	if ddxTimeout <= 0 {
		ddxTimeout = 8 * time.Second
	}
	return &ServiceContext{
		Config:         cfg,
		Conn:           conn,
		Repo:           repository.NewUserRepository(conn),
		RBAC:           repository.NewRBACRepository(conn),
		Tasks:          repository.NewTaskRepository(conn),
		PointsProducts: repository.NewPointsProductRepository(conn),
		PointsOrders:   repository.NewPointsOrderRepository(conn),
		Jutuike: jutuike.NewClient(jutuike.Config{
			ApiKey:  cfg.Jutuike.ApiKey,
			BaseURL: cfg.Jutuike.BaseURL,
			Timeout: jtkTimeout,
		}),
		Dingdanxia: dingdanxia.NewClient(dingdanxia.Config{
			ApiKey:  cfg.Dingdanxia.ApiKey,
			BaseURL: cfg.Dingdanxia.BaseURL,
			Timeout: ddxTimeout,
			PddPid:  cfg.Dingdanxia.PddPid,
		}),
		JWT: jwt.Config{
			Secret:      cfg.Auth.AccessSecret,
			ConsumerKey: cfg.Auth.ConsumerKey,
			ExpireHours: cfg.Auth.ExpireHours(),
			Issuer:      cfg.Auth.Issuer,
		},
		RequestID:            middleware.NewRequestIDMiddleware().Handle,
		GatewayIdentity:      middleware.NewGatewayIdentityMiddleware().Handle,
		RequirePlatformAdmin: middleware.NewRequirePlatformAdminMiddleware().Handle,
		UserFlexibleAuth:     flex.Handle,
	}
}
