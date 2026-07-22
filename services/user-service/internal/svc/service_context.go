package svc

import (
	"mymall/pkg/health"
	"mymall/pkg/jwt"
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
	JWT            jwt.Config
	Health         *health.Registry

	RequestID            rest.Middleware
	GatewayIdentity      rest.Middleware
	RequirePlatformAdmin rest.Middleware
	UserFlexibleAuth     rest.Middleware
}

func NewServiceContext(cfg *config.Config, conn sqlx.SqlConn) *ServiceContext {
	flex := middleware.NewUserFlexibleAuthMiddlewareWithSecret(cfg.Auth.AccessSecret)
	return &ServiceContext{
		Config:         cfg,
		Conn:           conn,
		Repo:           repository.NewUserRepository(conn),
		RBAC:           repository.NewRBACRepository(conn),
		Tasks:          repository.NewTaskRepository(conn),
		PointsProducts: repository.NewPointsProductRepository(conn),
		PointsOrders:   repository.NewPointsOrderRepository(conn),
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
