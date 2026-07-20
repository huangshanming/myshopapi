package svc

import (
	"mymall/pkg/cache"
	"mymall/pkg/config"
	"mymall/pkg/health"
	"mymall/pkg/mq"
	"mymall/services/order-service/internal/client/catalogrpc"
	"mymall/services/order-service/internal/client/merchanthttp"
	"mymall/services/order-service/internal/client/userhttp"
	"mymall/services/order-service/internal/client/userrpc"
	"mymall/services/order-service/internal/middleware"
	ordermq "mymall/services/order-service/internal/mq"
	"mymall/services/order-service/internal/repository"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

// ServiceContext 全局依赖（go-zero 惯例）
type ServiceContext struct {
	Config        *config.Config
	DB            *gorm.DB
	Redis         *redis.Client
	Repo          *repository.OrderRepository
	Reviews       *repository.ReviewRepository
	LogisticsRepo *repository.LogisticsRepository
	UserRPC       *userrpc.Client
	CatalogRPC    *catalogrpc.Client
	MerchantHTTP  *merchanthttp.Client
	UserHTTP      *userhttp.Client
	MQ            *ordermq.Publisher
	MQClient      *mq.Client
	Health        *health.Registry

	RequestID                 rest.Middleware
	GatewayIdentity           rest.Middleware
	GatewayIdentityShop       rest.Middleware
	RequireMerchantOwner      rest.Middleware
	RequirePlatformAdmin      rest.Middleware
	RequirePlatformOrMerchant rest.Middleware
}

func NewServiceContext(cfg *config.Config, db *gorm.DB) (*ServiceContext, error) {
	catalogRPC, err := catalogrpc.New(cfg.GRPC.CatalogService)
	if err != nil {
		return nil, err
	}

	var userRPC *userrpc.Client
	if u, err := userrpc.New(cfg.GRPC.UserService); err == nil {
		userRPC = u
	}

	var mqClient *mq.Client
	var publisher *ordermq.Publisher
	if mqc, err := mq.New(cfg.RabbitMQ); err == nil {
		mqClient = mqc
		publisher = ordermq.NewPublisher(mqc)
	}

	var rdb *redis.Client
	if c, err := cache.NewRedis(cfg.Redis); err == nil {
		rdb = c
	}

	logisticsRepo := repository.NewLogisticsRepository(db)
	_ = logisticsRepo.SeedDefaults()

	return &ServiceContext{
		Config:                    cfg,
		DB:                        db,
		Redis:                     rdb,
		Repo:                      repository.NewOrderRepository(db),
		Reviews:                   repository.NewReviewRepository(db),
		LogisticsRepo:             logisticsRepo,
		UserRPC:                   userRPC,
		CatalogRPC:                catalogRPC,
		MerchantHTTP:              merchanthttp.New(cfg.MerchantHTTP),
		UserHTTP:                  userhttp.New(cfg.UserHTTP),
		MQ:                        publisher,
		MQClient:                  mqClient,
		RequestID:                 middleware.NewRequestIDMiddleware().Handle,
		GatewayIdentity:           middleware.NewGatewayIdentityMiddleware().Handle,
		GatewayIdentityShop:       middleware.NewGatewayIdentityShopMiddleware().Handle,
		RequireMerchantOwner:      middleware.NewRequireMerchantOwnerMiddleware().Handle,
		RequirePlatformAdmin:      middleware.NewRequirePlatformAdminMiddleware().Handle,
		RequirePlatformOrMerchant: middleware.NewRequirePlatformOrMerchantMiddleware().Handle,
	}, nil
}

func (s *ServiceContext) Close() {
	if s.CatalogRPC != nil {
		_ = s.CatalogRPC.Close()
	}
	if s.UserRPC != nil {
		_ = s.UserRPC.Close()
	}
	if s.MQClient != nil {
		_ = s.MQClient.Close()
	}
	if s.Redis != nil {
		_ = s.Redis.Close()
	}
}
