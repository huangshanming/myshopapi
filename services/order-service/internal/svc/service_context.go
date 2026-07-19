package svc

import (
	"mymall/pkg/cache"
	"mymall/pkg/config"
	"mymall/pkg/mq"
	"mymall/services/order-service/internal/client/catalogrpc"
	"mymall/services/order-service/internal/client/merchanthttp"
	"mymall/services/order-service/internal/client/userhttp"
	"mymall/services/order-service/internal/client/userrpc"
	ordermq "mymall/services/order-service/internal/mq"
	"mymall/services/order-service/internal/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ServiceContext 全局依赖（go-zero 惯例）
type ServiceContext struct {
	Config        *config.Config
	DB            *gorm.DB
	Redis         *redis.Client
	Repo          *repository.OrderRepository
	LogisticsRepo *repository.LogisticsRepository
	UserRPC       *userrpc.Client
	CatalogRPC    *catalogrpc.Client
	MerchantHTTP  *merchanthttp.Client
	UserHTTP      *userhttp.Client
	MQ            *ordermq.Publisher
	MQClient      *mq.Client
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
		Config:        cfg,
		DB:            db,
		Redis:         rdb,
		Repo:          repository.NewOrderRepository(db),
		LogisticsRepo: logisticsRepo,
		UserRPC:       userRPC,
		CatalogRPC:    catalogRPC,
		MerchantHTTP:  merchanthttp.New(cfg.MerchantHTTP),
		UserHTTP:      userhttp.New(cfg.UserHTTP),
		MQ:            publisher,
		MQClient:      mqClient,
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
