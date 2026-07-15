package svc

import (
	"mymall/pkg/config"
	"mymall/pkg/jwt"
	"mymall/services/user-service/internal/repository"

	"gorm.io/gorm"
)

// ServiceContext 全局依赖（go-zero 惯例）
type ServiceContext struct {
	Config *config.Config
	DB     *gorm.DB
	Repo   *repository.UserRepository
	RBAC   *repository.RBACRepository
	JWT    jwt.Config
}

func NewServiceContext(cfg *config.Config, db *gorm.DB) *ServiceContext {
	return &ServiceContext{
		Config: cfg,
		DB:     db,
		Repo:   repository.NewUserRepository(db),
		RBAC:   repository.NewRBACRepository(db),
		JWT: jwt.Config{
			Secret:      cfg.JWT.Secret,
			ConsumerKey: cfg.JWT.ConsumerKey,
			ExpireHours: cfg.JWT.ExpireHours,
			Issuer:      cfg.JWT.Issuer,
		},
	}
}
