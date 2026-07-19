package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"mymall/pkg/config"
	"mymall/pkg/database"
	"mymall/pkg/health"
	"mymall/pkg/httpserver"
	"mymall/pkg/jwt"
	applog "mymall/pkg/log"
	"mymall/pkg/metrics"
	"mymall/pkg/middleware"
	"mymall/pkg/telemetry"
	"mymall/services/user-service/internal/data"
	"mymall/services/user-service/internal/handler"
	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/model"
	"mymall/services/user-service/internal/server"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./etc/user-service.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("加载配置失败：%v", err)
	}

	logger, err := applog.New("user-service")
	if err != nil {
		log.Fatalf("初始化日志失败：%v", err)
	}
	defer logger.Sync()

	ctx := context.Background()
	shutdownTrace, err := telemetry.Init(ctx, cfg.Telemetry)
	if err != nil {
		logger.Warn("telemetry init skipped")
	}
	defer shutdownTrace(context.Background())

	db, err := database.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("连接数据库失败：%v", err)
	}
	if err := database.AutoMigrateIfDebug(cfg.Server.Mode, db,
		&model.User{},
		&model.SysMenu{},
		&model.SysRole{},
		&model.SysRoleMenu{},
		&model.SysUserRole{},
		&model.SysConfig{},
		&model.UserWallet{},
		&model.UserWalletLog{},
		&model.UserAddress{},
		&model.Region{},
	); err != nil {
		log.Fatalf("AutoMigrate 失败：%v", err)
	}

	svcCtx := svc.NewServiceContext(cfg, db)
	if n, err := svcCtx.Repo.CountRegions(); err == nil && n == 0 {
		if err := svcCtx.Repo.SeedRegionsFromPCA(data.PCACodeJSON); err != nil {
			logger.Warn(fmt.Sprintf("seed regions failed: %v", err))
		} else {
			logger.Info("regions seeded from pca-code.json")
		}
	}
	userLogic := logic.NewUserLogic(svcCtx)
	userHandler := handler.NewUserHandler(svcCtx)
	adminHandler := handler.NewAdminHandler(svcCtx)
	walletHandler := handler.NewWalletHandler(svcCtx)
	addressHandler := handler.NewAddressHandler(svcCtx)
	regionHandler := handler.NewRegionHandler(svcCtx)

	healthReg := health.NewRegistry()
	healthReg.Register("mysql", func(ctx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.PingContext(ctx)
	})

	rpcServer := server.StartZRPC(cfg.Server.GRPCPort, userLogic, logger)
	go func() {
		logger.Info(fmt.Sprintf("user-service zRPC 启动 :%d", cfg.Server.GRPCPort))
		rpcServer.Start()
	}()
	defer rpcServer.Stop()

	serverHTTP := httpserver.NewRest(cfg.Server.HTTPPort, cfg.Server.Mode)
	defer serverHTTP.Stop()

	rid := middleware.RequestID()
	authJWT := jwt.AuthMiddleware(svcCtx.JWT.Secret)
	authGW := middleware.GatewayIdentity(false)
	plat := middleware.RequireRoles(jwt.RolePlatformAdmin)
	perm := func(code string) middleware.Middleware {
		return middleware.RequirePermission(adminHandler, code)
	}
	adminAuth := func(code string, h http.HandlerFunc) http.HandlerFunc {
		return rid(middleware.Chain(h, authGW, plat, perm(code)))
	}

	profile := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(middleware.GatewayUserIDHeader) != "" {
			authGW(userHandler.Profile)(w, r)
			return
		}
		authJWT(userHandler.Profile)(w, r)
	}
	userAuth := func(h http.HandlerFunc) http.HandlerFunc {
		return rid(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(middleware.GatewayUserIDHeader) != "" {
				authGW(h)(w, r)
				return
			}
			authJWT(h)(w, r)
		})
	}

	serverHTTP.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/healthz", Handler: rid(httpserver.Healthz("user-service"))},
		{Method: http.MethodGet, Path: "/readyz", Handler: rid(healthReg.ReadyHandler())},
		{Method: http.MethodGet, Path: "/metrics", Handler: rid(metrics.Handler())},

		{Method: http.MethodPost, Path: "/api/v1/user/login", Handler: rid(userHandler.Login)},
		{Method: http.MethodPost, Path: "/api/v1/user/register", Handler: rid(userHandler.Register)},
		{Method: http.MethodGet, Path: "/api/v1/user/profile", Handler: rid(profile)},
		{Method: http.MethodGet, Path: "/api/v1/user/wallet", Handler: userAuth(walletHandler.UserGetWallet)},
		{Method: http.MethodGet, Path: "/api/v1/user/wallet/logs", Handler: userAuth(walletHandler.UserWalletLogs)},
		{Method: http.MethodPost, Path: "/api/v1/user/wallet/freeze", Handler: rid(walletHandler.Freeze)},
		{Method: http.MethodPost, Path: "/api/v1/user/wallet/unfreeze", Handler: rid(walletHandler.Unfreeze)},
		{Method: http.MethodPost, Path: "/api/v1/user/wallet/settle", Handler: rid(walletHandler.Settle)},

		{Method: http.MethodGet, Path: "/api/v1/user/addresses", Handler: userAuth(addressHandler.List)},
		{Method: http.MethodPost, Path: "/api/v1/user/addresses", Handler: userAuth(addressHandler.Create)},
		{Method: http.MethodPut, Path: "/api/v1/user/addresses/:id", Handler: userAuth(addressHandler.Update)},
		{Method: http.MethodDelete, Path: "/api/v1/user/addresses/:id", Handler: userAuth(addressHandler.Delete)},
		{Method: http.MethodPut, Path: "/api/v1/user/addresses/:id/default", Handler: userAuth(addressHandler.SetDefault)},
		{Method: http.MethodGet, Path: "/api/v1/user/addresses/internal", Handler: rid(addressHandler.InternalGet)},

		{Method: http.MethodGet, Path: "/api/v1/regions", Handler: rid(regionHandler.List)},
		{Method: http.MethodGet, Path: "/api/v1/regions/tree", Handler: rid(regionHandler.Tree)},

		{Method: http.MethodGet, Path: "/api/v1/admin/auth/me", Handler: adminAuth("", adminHandler.AuthMe)},

		{Method: http.MethodGet, Path: "/api/v1/admin/menus", Handler: adminAuth("system:menu:list", adminHandler.MenuTree)},
		{Method: http.MethodPost, Path: "/api/v1/admin/menus", Handler: adminAuth("system:menu:add", adminHandler.CreateMenu)},
		{Method: http.MethodPut, Path: "/api/v1/admin/menus/:id", Handler: adminAuth("system:menu:edit", adminHandler.UpdateMenu)},
		{Method: http.MethodDelete, Path: "/api/v1/admin/menus/:id", Handler: adminAuth("system:menu:delete", adminHandler.DeleteMenu)},

		{Method: http.MethodGet, Path: "/api/v1/admin/roles", Handler: adminAuth("system:role:list", adminHandler.ListRoles)},
		{Method: http.MethodPost, Path: "/api/v1/admin/roles", Handler: adminAuth("system:role:add", adminHandler.CreateRole)},
		{Method: http.MethodPut, Path: "/api/v1/admin/roles/:id", Handler: adminAuth("system:role:edit", adminHandler.UpdateRole)},
		{Method: http.MethodDelete, Path: "/api/v1/admin/roles/:id", Handler: adminAuth("system:role:delete", adminHandler.DeleteRole)},
		{Method: http.MethodGet, Path: "/api/v1/admin/roles/:id/menus", Handler: adminAuth("system:role:list", adminHandler.GetRoleMenus)},
		{Method: http.MethodPut, Path: "/api/v1/admin/roles/:id/menus", Handler: adminAuth("system:role:assign", adminHandler.AssignRoleMenus)},

		{Method: http.MethodGet, Path: "/api/v1/admin/users", Handler: adminAuth("system:user:list", adminHandler.ListUsers)},
		{Method: http.MethodGet, Path: "/api/v1/admin/users/:id", Handler: adminAuth("system:user:list", adminHandler.GetUser)},
		{Method: http.MethodPut, Path: "/api/v1/admin/users/:id", Handler: adminAuth("system:user:edit", adminHandler.UpdateUser)},
		{Method: http.MethodPut, Path: "/api/v1/admin/users/:id/status", Handler: adminAuth("system:user:status", adminHandler.SetUserStatus)},
		{Method: http.MethodPut, Path: "/api/v1/admin/users/:id/password", Handler: adminAuth("system:user:reset", adminHandler.ResetUserPassword)},
		{Method: http.MethodPost, Path: "/api/v1/admin/users/:id/token", Handler: adminAuth("system:user:list", adminHandler.GenerateUserToken)},
		{Method: http.MethodGet, Path: "/api/v1/admin/users/:id/wallet", Handler: adminAuth("system:user:wallet", walletHandler.AdminGetWallet)},
		{Method: http.MethodPost, Path: "/api/v1/admin/users/:id/wallet/adjust", Handler: adminAuth("system:user:wallet", walletHandler.AdminAdjustWallet)},
		{Method: http.MethodGet, Path: "/api/v1/admin/users/:id/wallet/logs", Handler: adminAuth("system:user:wallet", walletHandler.AdminWalletLogs)},
		{Method: http.MethodGet, Path: "/api/v1/admin/users/:id/addresses", Handler: adminAuth("system:user:list", addressHandler.AdminList)},

		{Method: http.MethodGet, Path: "/api/v1/admin/admins", Handler: adminAuth("system:admin:list", adminHandler.ListAdmins)},
		{Method: http.MethodPost, Path: "/api/v1/admin/admins", Handler: adminAuth("system:admin:add", adminHandler.CreateAdmin)},
		{Method: http.MethodGet, Path: "/api/v1/admin/admins/:id/roles", Handler: adminAuth("system:admin:list", adminHandler.GetAdminRoles)},
		{Method: http.MethodPut, Path: "/api/v1/admin/admins/:id/roles", Handler: adminAuth("system:admin:assign", adminHandler.AssignAdminRoles)},
		{Method: http.MethodPut, Path: "/api/v1/admin/admins/:id/password", Handler: adminAuth("system:admin:reset", adminHandler.ResetAdminPassword)},

		{Method: http.MethodGet, Path: "/api/v1/admin/configs", Handler: adminAuth("system:config:list", adminHandler.ListConfigs)},
		{Method: http.MethodPut, Path: "/api/v1/admin/configs", Handler: adminAuth("system:config:edit", adminHandler.SaveConfigs)},
	})

	go func() {
		logger.Info(fmt.Sprintf("user-service HTTP(go-zero) 启动 :%d", cfg.Server.HTTPPort))
		serverHTTP.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
