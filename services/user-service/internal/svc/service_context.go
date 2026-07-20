package svc

import (
	"mymall/pkg/config"
	"mymall/pkg/health"
	"mymall/pkg/jwt"
	"mymall/services/user-service/internal/middleware"
	"mymall/services/user-service/internal/repository"

	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config         *config.Config
	DB             *gorm.DB
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

func NewServiceContext(cfg *config.Config, db *gorm.DB) *ServiceContext {
	flex := middleware.NewUserFlexibleAuthMiddlewareWithSecret(cfg.JWT.Secret)
	return &ServiceContext{
		Config:         cfg,
		DB:             db,
		Repo:           repository.NewUserRepository(db),
		RBAC:           repository.NewRBACRepository(db),
		Tasks:          repository.NewTaskRepository(db),
		PointsProducts: repository.NewPointsProductRepository(db),
		PointsOrders:   repository.NewPointsOrderRepository(db),
		JWT: jwt.Config{
			Secret:      cfg.JWT.Secret,
			ConsumerKey: cfg.JWT.ConsumerKey,
			ExpireHours: cfg.JWT.ExpireHours,
			Issuer:      cfg.JWT.Issuer,
		},
		RequestID:            middleware.NewRequestIDMiddleware().Handle,
		GatewayIdentity:      middleware.NewGatewayIdentityMiddleware().Handle,
		RequirePlatformAdmin: middleware.NewRequirePlatformAdminMiddleware().Handle,
		UserFlexibleAuth:     flex.Handle,
	}
}
