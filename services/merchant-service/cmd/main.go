package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mymall/pkg/config"
	"mymall/pkg/database"
	"mymall/pkg/health"
	"mymall/pkg/httpserver"
	"mymall/pkg/jwt"
	applog "mymall/pkg/log"
	"mymall/pkg/middleware"
	"mymall/services/merchant-service/internal/handler"
	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/model"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func main() {
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
	); err != nil {
		log.Fatalf("AutoMigrate 失败：%v", err)
	}

	svcCtx := svc.NewServiceContext(cfg, db)
	h := handler.NewMerchantHandler(svcCtx)
	seckillLogic := logic.NewMerchantLogic(svcCtx)
	_, _, _ = seckillLogic.EnsureActiveSession()

	go func() {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for range t.C {
			seckillLogic.RotateSeckillSessions()
		}
	}()

	healthReg := health.NewRegistry()
	healthReg.Register("mysql", func(ctx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.PingContext(ctx)
	})

	server := httpserver.NewRest(cfg.Server.HTTPPort, cfg.Server.Mode)
	defer server.Stop()

	rid := middleware.RequestID()
	gw := middleware.GatewayIdentity(false)
	owner := middleware.RequireRoles(jwt.RoleMerchantOwner, jwt.RoleMerchantStaff)
	plat := middleware.RequireRoles(jwt.RolePlatformAdmin)

	server.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/healthz", Handler: rid(httpserver.Healthz("merchant-service"))},
		{Method: http.MethodGet, Path: "/readyz", Handler: rid(healthReg.ReadyHandler())},

		{Method: http.MethodGet, Path: "/api/v1/shops/list", Handler: rid(h.PublicListShops)},
		{Method: http.MethodGet, Path: "/api/v1/shops/:id", Handler: rid(h.PublicGetShop)},
		{Method: http.MethodGet, Path: "/api/v1/seckill/current", Handler: rid(h.PublicSeckillCurrent)},
		{Method: http.MethodGet, Path: "/api/v1/seckill/list", Handler: rid(h.PublicSeckillList)},
		{Method: http.MethodGet, Path: "/api/v1/seckill/entries/:id", Handler: rid(h.PublicSeckillEntry)},
		{Method: http.MethodPost, Path: "/api/v1/seckill/consume", Handler: rid(h.SeckillConsume)},
		{Method: http.MethodPost, Path: "/api/v1/seckill/restore", Handler: rid(h.SeckillRestore)},

		{Method: http.MethodPost, Path: "/api/v1/merchant/apply", Handler: rid(gw(h.Apply))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/shops", Handler: rid(gw(h.MyShops))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/shops/:id", Handler: rid(gw(owner(h.UpdateMyShop)))},

		{Method: http.MethodGet, Path: "/api/v1/merchant/wallet", Handler: rid(gw(owner(h.MerchantGetWallet)))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/wallet/logs", Handler: rid(gw(owner(h.MerchantWalletLogs)))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/seckill/sessions", Handler: rid(gw(owner(h.MerchantSeckillSessions)))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/seckill/entries", Handler: rid(gw(owner(h.MerchantApplySeckill)))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/seckill/entries", Handler: rid(gw(owner(h.MerchantListSeckillEntries)))},

		{Method: http.MethodGet, Path: "/api/v1/admin/applications", Handler: rid(middleware.Chain(h.AdminListApplications, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/applications/:id/approve", Handler: rid(middleware.Chain(h.AdminApprove, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/applications/:id/reject", Handler: rid(middleware.Chain(h.AdminReject, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/shops", Handler: rid(middleware.Chain(h.AdminListShops, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/shops", Handler: rid(middleware.Chain(h.AdminCreateShop, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/shops/:id", Handler: rid(middleware.Chain(h.AdminGetShop, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/shops/:id", Handler: rid(middleware.Chain(h.AdminUpdateShop, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/shops/:id/owner-password", Handler: rid(middleware.Chain(h.AdminResetOwnerPassword, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/shops/:id/disable", Handler: rid(middleware.Chain(h.AdminDisableShop, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/shops/:id/enable", Handler: rid(middleware.Chain(h.AdminEnableShop, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/shops/:id/wallet", Handler: rid(middleware.Chain(h.AdminGetWallet, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/shops/:id/wallet/adjust", Handler: rid(middleware.Chain(h.AdminAdjustWallet, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/shops/:id/wallet/logs", Handler: rid(middleware.Chain(h.AdminWalletLogs, gw, plat))},

		{Method: http.MethodGet, Path: "/api/v1/admin/seckill/rule", Handler: rid(middleware.Chain(h.AdminGetSeckillRule, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/seckill/rule", Handler: rid(middleware.Chain(h.AdminUpdateSeckillRule, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/seckill/sessions", Handler: rid(middleware.Chain(h.AdminListSeckillSessions, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/seckill/entries", Handler: rid(middleware.Chain(h.AdminListSeckillEntries, gw, plat))},
	})

	go func() {
		logger.Info(fmt.Sprintf("merchant-service HTTP(go-zero) 启动 :%d", cfg.Server.HTTPPort))
		server.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	server.Stop()
}
