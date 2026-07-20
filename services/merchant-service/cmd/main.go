package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mymall/pkg/config"
	"mymall/pkg/database"
	"mymall/pkg/health"
	"mymall/pkg/httpserver"
	applog "mymall/pkg/log"
	"mymall/pkg/xerr"
	biz "mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/handler"
	"mymall/services/merchant-service/internal/model"
	"mymall/services/merchant-service/internal/svc"
)

func main() {
	xerr.RegisterErrorHandler()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./etc/merchant-service.yaml"
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("加载配置失败：%v", err)
	}

	logger, err := applog.New("merchant-service")
	if err != nil {
		log.Fatalf("初始化日志失败：%v", err)
	}
	defer logger.Sync()

	db, err := database.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("连接数据库失败：%v", err)
	}
	if err := database.AutoMigrateIfDebug(cfg.Server.Mode, db,
		&model.Shop{},
		&model.ShopApplication{},
		&model.ShopMember{},
		&model.ShopWallet{},
		&model.ShopWalletLog{},
		&model.SeckillRule{},
		&model.SeckillSession{},
		&model.SeckillEntry{},
		&model.HomepageSlotPackage{},
		&model.HomepageSlotSetting{},
		&model.HomepageSlotOrder{},
		&model.HomepageThemeSlot{},
		&model.HomepageThemePackage{},
		&model.HomepageThemeOrder{},
		&model.Coupon{},
		&model.CouponScope{},
		&model.UserCoupon{},
		&model.CouponGrant{},
		&model.CouponRedeemLog{},
		&model.PointsProduct{},
		&model.PointsExchangeOrder{},
	); err != nil {
		log.Fatalf("AutoMigrate 失败：%v", err)
	}

	healthReg := health.NewRegistry()
	healthReg.Register("mysql", func(ctx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.PingContext(ctx)
	})
	svcCtx := svc.NewServiceContext(cfg, db, healthReg)
	seckillLogic := biz.NewMerchantLogic(context.Background(), svcCtx)
	_, _, _ = seckillLogic.EnsureActiveSession()

	go func() {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for range t.C {
			seckillLogic.RotateSeckillSessions()
		}
	}()

	server := httpserver.NewRest(cfg.Server.HTTPPort, cfg.Server.Mode)
	defer server.Stop()

	handler.RegisterHandlers(server, svcCtx)

	go func() {
		logger.Info(fmt.Sprintf("merchant-service HTTP(go-zero) 启动 :%d", cfg.Server.HTTPPort))
		server.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	server.Stop()
}
