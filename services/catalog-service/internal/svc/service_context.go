package svc

import (
	"mymall/pkg/config"
	"mymall/pkg/mq"
	"mymall/services/catalog-service/internal/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ServiceContext 全局依赖（go-zero 惯例）
type ServiceContext struct {
	Config       *config.Config
	DB           *gorm.DB
	Redis        *redis.Client
	MQ           *mq.Client
	Products     *repository.ProductRepository
	Categories   *repository.CategoryRepository
	ProductAdmin *repository.ProductAdminRepository
	ShopRBAC     *repository.ShopRBACRepository
	Articles     *repository.ArticleRepository
}

func NewServiceContext(cfg *config.Config, db *gorm.DB, redisClient *redis.Client, mqClient *mq.Client) *ServiceContext {
	return &ServiceContext{
		Config:       cfg,
		DB:           db,
		Redis:        redisClient,
		MQ:           mqClient,
		Products:     repository.NewProductRepository(db),
		Categories:   repository.NewCategoryRepository(db),
		ProductAdmin: repository.NewProductAdminRepository(db),
		ShopRBAC:     repository.NewShopRBACRepository(db),
		Articles:     repository.NewArticleRepository(db),
	}
}
