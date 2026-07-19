package svc

import (
	"mymall/pkg/config"
	"mymall/pkg/mq"
	contentrepo "mymall/services/catalog-service/internal/content/repository"
	notifyrepo "mymall/services/catalog-service/internal/notify/repository"
	productrepo "mymall/services/catalog-service/internal/product/repository"
	shopopsrepo "mymall/services/catalog-service/internal/shopops/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ServiceContext 全局依赖（go-zero 惯例）
type ServiceContext struct {
	Config        *config.Config
	DB            *gorm.DB
	Redis         *redis.Client
	MQ            *mq.Client
	Products      *productrepo.ProductRepository
	Categories    *productrepo.CategoryRepository
	ProductAdmin  *productrepo.ProductAdminRepository
	Favorites     *productrepo.FavoriteRepository
	ShopRBAC      *shopopsrepo.ShopRBACRepository
	Articles      *contentrepo.ArticleRepository
	Notifications *notifyrepo.NotificationRepository
}

func NewServiceContext(cfg *config.Config, db *gorm.DB, redisClient *redis.Client, mqClient *mq.Client) *ServiceContext {
	return &ServiceContext{
		Config:        cfg,
		DB:            db,
		Redis:         redisClient,
		MQ:            mqClient,
		Products:      productrepo.NewProductRepository(db),
		Categories:    productrepo.NewCategoryRepository(db),
		ProductAdmin:  productrepo.NewProductAdminRepository(db),
		Favorites:     productrepo.NewFavoriteRepository(db),
		ShopRBAC:      shopopsrepo.NewShopRBACRepository(db),
		Articles:      contentrepo.NewArticleRepository(db),
		Notifications: notifyrepo.NewNotificationRepository(db),
	}
}
