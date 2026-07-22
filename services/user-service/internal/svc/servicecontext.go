// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"mymall/services/user-service/internal/config"
	"mymall/services/user-service/internal/middleware"
)

type ServiceContext struct {
	Config               config.Config
	RequestID            rest.Middleware
	GatewayIdentity      rest.Middleware
	RequirePlatformAdmin rest.Middleware
	UserFlexibleAuth     rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:               c,
		RequestID:            middleware.NewRequestIDMiddleware().Handle,
		GatewayIdentity:      middleware.NewGatewayIdentityMiddleware().Handle,
		RequirePlatformAdmin: middleware.NewRequirePlatformAdminMiddleware().Handle,
		UserFlexibleAuth:     middleware.NewUserFlexibleAuthMiddleware().Handle,
	}
}
