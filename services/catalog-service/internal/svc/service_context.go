package svc

import (
	"mymall/pkg/health"
	"mymall/pkg/mq"
	"mymall/services/catalog-service/internal/client/userrpc"
	"mymall/services/catalog-service/internal/config"
	contentrepo "mymall/services/catalog-service/internal/content/repository"
	"mymall/services/catalog-service/internal/middleware"
	notifyrepo "mymall/services/catalog-service/internal/notify/repository"
	productrepo "mymall/services/catalog-service/internal/product/repository"
	shopopsrepo "mymall/services/catalog-service/internal/shopops/repository"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

// ServiceContext 全局依赖（go-zero 惯例）
type ServiceContext struct {
	Config        *config.Config
	Conn          sqlx.SqlConn
	Redis         *redis.Client
	MQ            *mq.Client
	Products      *productrepo.ProductRepository
	Categories    *productrepo.CategoryRepository
	ProductAdmin  *productrepo.ProductAdminRepository
	Favorites     *productrepo.FavoriteRepository
	ShopRBAC      *shopopsrepo.ShopRBACRepository
	Articles      *contentrepo.ArticleRepository
	Notifications *notifyrepo.NotificationRepository
	UserRPC       *userrpc.Client
	Health        *health.Registry

	RequestID            rest.Middleware
	GatewayIdentity      rest.Middleware
	GatewayIdentityShop  rest.Middleware
	RequireMerchantOwner rest.Middleware
	RequirePlatformAdmin rest.Middleware
}

func NewServiceContext(cfg *config.Config, conn sqlx.SqlConn, redisClient *redis.Client, mqClient *mq.Client) *ServiceContext {
	var userRPC *userrpc.Client
	if u, err := userrpc.New(cfg.UserRpc); err == nil {
		userRPC = u
	}
	return &ServiceContext{
		Config:               cfg,
		Conn:                 conn,
		Redis:                redisClient,
		MQ:                   mqClient,
		Products:             productrepo.NewProductRepository(conn),
		Categories:           productrepo.NewCategoryRepository(conn),
		ProductAdmin:         productrepo.NewProductAdminRepository(conn),
		Favorites:            productrepo.NewFavoriteRepository(conn),
		ShopRBAC:             shopopsrepo.NewShopRBACRepository(conn),
		Articles:             contentrepo.NewArticleRepository(conn),
		Notifications:        notifyrepo.NewNotificationRepository(conn),
		UserRPC:              userRPC,
		RequestID:            middleware.NewRequestIDMiddleware().Handle,
		GatewayIdentity:      middleware.NewGatewayIdentityMiddleware().Handle,
		GatewayIdentityShop:  middleware.NewGatewayIdentityShopMiddleware().Handle,
		RequireMerchantOwner: middleware.NewRequireMerchantOwnerMiddleware().Handle,
		RequirePlatformAdmin: middleware.NewRequirePlatformAdminMiddleware().Handle,
	}
}
