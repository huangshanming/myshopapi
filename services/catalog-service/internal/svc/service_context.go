package svc

import (
	"mymall/pkg/config"
	"mymall/pkg/health"
	"mymall/pkg/mq"
	"mymall/services/catalog-service/internal/client/userhttp"
	contentrepo "mymall/services/catalog-service/internal/content/repository"
	"mymall/services/catalog-service/internal/middleware"
	notifyrepo "mymall/services/catalog-service/internal/notify/repository"
	productrepo "mymall/services/catalog-service/internal/product/repository"
	shopopsrepo "mymall/services/catalog-service/internal/shopops/repository"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
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
	UserHTTP      *userhttp.Client
	Health        *health.Registry

	RequestID            rest.Middleware
	GatewayIdentity      rest.Middleware
	GatewayIdentityShop  rest.Middleware
	RequireMerchantOwner rest.Middleware
	RequirePlatformAdmin rest.Middleware
}

func NewServiceContext(cfg *config.Config, db *gorm.DB, redisClient *redis.Client, mqClient *mq.Client) *ServiceContext {
	return &ServiceContext{
		Config:               cfg,
		DB:                   db,
		Redis:                redisClient,
		MQ:                   mqClient,
		Products:             productrepo.NewProductRepository(db),
		Categories:           productrepo.NewCategoryRepository(db),
		ProductAdmin:         productrepo.NewProductAdminRepository(db),
		Favorites:            productrepo.NewFavoriteRepository(db),
		ShopRBAC:             shopopsrepo.NewShopRBACRepository(db),
		Articles:             contentrepo.NewArticleRepository(db),
		Notifications:        notifyrepo.NewNotificationRepository(db),
		UserHTTP:             userhttp.New(cfg.UserHTTP),
		RequestID:            middleware.NewRequestIDMiddleware().Handle,
		GatewayIdentity:      middleware.NewGatewayIdentityMiddleware().Handle,
		GatewayIdentityShop:  middleware.NewGatewayIdentityShopMiddleware().Handle,
		RequireMerchantOwner: middleware.NewRequireMerchantOwnerMiddleware().Handle,
		RequirePlatformAdmin: middleware.NewRequirePlatformAdminMiddleware().Handle,
	}
}
