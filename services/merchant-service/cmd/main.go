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
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/handler"
	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/model"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/uploadpath"

	"github.com/zeromicro/go-zero/rest"
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

	svcCtx := svc.NewServiceContext(cfg, db)
	shopAdmin := handler.NewShopAdminHandler(svcCtx)
	shopMerchant := handler.NewShopMerchantHandler(svcCtx)
	shopPublic := handler.NewShopPublicHandler(svcCtx)
	walletAdmin := handler.NewWalletAdminHandler(svcCtx)
	walletMerchant := handler.NewWalletMerchantHandler(svcCtx)
	seckillAdmin := handler.NewSeckillAdminHandler(svcCtx)
	seckillMerchant := handler.NewSeckillMerchantHandler(svcCtx)
	seckillPublic := handler.NewSeckillPublicHandler(svcCtx)
	seckillInternal := handler.NewSeckillInternalHandler(svcCtx)
	couponAdmin := handler.NewCouponAdminHandler(svcCtx)
	couponMerchant := handler.NewCouponMerchantHandler(svcCtx)
	couponPublic := handler.NewCouponPublicHandler(svcCtx)
	couponUser := handler.NewCouponUserHandler(svcCtx)
	couponInternal := handler.NewCouponInternalHandler(svcCtx)
	slotAdmin := handler.NewHomepageSlotAdminHandler(svcCtx)
	slotMerchant := handler.NewHomepageSlotMerchantHandler(svcCtx)
	slotPublic := handler.NewHomepageSlotPublicHandler(svcCtx)
	themeAdmin := handler.NewHomepageThemeAdminHandler(svcCtx)
	themeMerchant := handler.NewHomepageThemeMerchantHandler(svcCtx)
	themePublic := handler.NewHomepageThemePublicHandler(svcCtx)
	pointsProductAdmin := handler.NewPointsProductAdminHandler(svcCtx)
	pointsOrderAdmin := handler.NewPointsOrderAdminHandler(svcCtx)
	pointsOrderUser := handler.NewPointsOrderUserHandler(svcCtx)
	seckillLogic := logic.NewMerchantLogic(context.Background(), svcCtx)
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
	admin := func(h http.HandlerFunc) http.HandlerFunc {
		return rid(middleware.Chain(h, gw, plat))
	}
	userAuth := func(h http.HandlerFunc) http.HandlerFunc {
		return rid(gw(h))
	}

	server.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/healthz", Handler: rid(httpserver.Healthz("merchant-service"))},
		{Method: http.MethodGet, Path: "/readyz", Handler: rid(healthReg.ReadyHandler())},

		{Method: http.MethodGet, Path: "/api/v1/shops/list", Handler: rid(shopPublic.PublicListShops)},
		{Method: http.MethodGet, Path: "/api/v1/shops/home-slots", Handler: rid(slotPublic.PublicHomeSlots)},
		{Method: http.MethodGet, Path: "/api/v1/home/theme-tiles", Handler: rid(themePublic.PublicThemeTiles)},
		{Method: http.MethodGet, Path: "/api/v1/shops/:id", Handler: rid(shopPublic.PublicGetShop)},
		{Method: http.MethodGet, Path: "/api/v1/seckill/current", Handler: rid(seckillPublic.PublicSeckillCurrent)},
		{Method: http.MethodGet, Path: "/api/v1/seckill/list", Handler: rid(seckillPublic.PublicSeckillList)},
		{Method: http.MethodGet, Path: "/api/v1/seckill/entries/:id", Handler: rid(seckillPublic.PublicSeckillEntry)},
		{Method: http.MethodPost, Path: "/api/v1/seckill/consume", Handler: rid(seckillInternal.SeckillConsume)},
		{Method: http.MethodPost, Path: "/api/v1/seckill/restore", Handler: rid(seckillInternal.SeckillRestore)},

		{Method: http.MethodGet, Path: "/api/v1/coupons/center", Handler: rid(couponPublic.PublicCouponCenter)},
		{Method: http.MethodGet, Path: "/api/v1/coupons/popup", Handler: rid(couponPublic.PublicCouponPopup)},
		{Method: http.MethodPost, Path: "/api/v1/coupons/:id/claim", Handler: rid(gw(couponUser.ClaimCoupon))},
		{Method: http.MethodGet, Path: "/api/v1/user/coupons", Handler: rid(gw(couponUser.ListMyCoupons))},

		{Method: http.MethodPost, Path: "/api/v1/internal/coupons/match", Handler: rid(couponInternal.InternalMatchCoupons)},
		{Method: http.MethodPost, Path: "/api/v1/internal/coupons/lock", Handler: rid(couponInternal.InternalLockCoupon)},
		{Method: http.MethodPost, Path: "/api/v1/internal/coupons/unlock", Handler: rid(couponInternal.InternalUnlockCoupon)},
		{Method: http.MethodPost, Path: "/api/v1/internal/coupons/redeem", Handler: rid(couponInternal.InternalRedeemCoupon)},
		{Method: http.MethodPost, Path: "/api/v1/internal/coupons/return", Handler: rid(couponInternal.InternalReturnCoupon)},
		{Method: http.MethodPost, Path: "/api/v1/internal/coupons/order-gift", Handler: rid(couponInternal.InternalOrderGift)},

		{Method: http.MethodPost, Path: "/api/v1/merchant/apply", Handler: rid(gw(shopMerchant.Apply))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/shops", Handler: rid(gw(shopMerchant.MyShops))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/shops/:id", Handler: rid(gw(owner(shopMerchant.UpdateMyShop)))},

		{Method: http.MethodGet, Path: "/api/v1/merchant/wallet", Handler: rid(gw(owner(walletMerchant.MerchantGetWallet)))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/wallet/logs", Handler: rid(gw(owner(walletMerchant.MerchantWalletLogs)))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/seckill/sessions", Handler: rid(gw(owner(seckillMerchant.MerchantSeckillSessions)))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/seckill/entries", Handler: rid(gw(owner(seckillMerchant.MerchantApplySeckill)))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/seckill/entries", Handler: rid(gw(owner(seckillMerchant.MerchantListSeckillEntries)))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/seckill/entries/:id/auto-renew", Handler: rid(gw(owner(seckillMerchant.MerchantSetSeckillAutoRenew)))},

		{Method: http.MethodGet, Path: "/api/v1/merchant/homepage-packages", Handler: rid(gw(owner(slotMerchant.MerchantListSlotPackages)))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/homepage-orders", Handler: rid(gw(owner(slotMerchant.MerchantBuySlot)))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/homepage-orders", Handler: rid(gw(owner(slotMerchant.MerchantListSlotOrders)))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/theme-slots", Handler: rid(gw(owner(themeMerchant.MerchantListThemeSlots)))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/theme-packages", Handler: rid(gw(owner(themeMerchant.MerchantListThemePackages)))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/theme-orders", Handler: rid(gw(owner(themeMerchant.MerchantBuyTheme)))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/theme-orders", Handler: rid(gw(owner(themeMerchant.MerchantListThemeOrders)))},

		{Method: http.MethodGet, Path: "/api/v1/merchant/coupons", Handler: rid(gw(owner(couponMerchant.MerchantListCoupons)))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/coupons", Handler: rid(gw(owner(couponMerchant.MerchantCreateCoupon)))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/coupons/:id", Handler: rid(gw(owner(couponMerchant.MerchantUpdateCoupon)))},
		{Method: http.MethodPut, Path: "/api/v1/merchant/coupons/:id/off", Handler: rid(gw(owner(couponMerchant.MerchantOffCoupon)))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/coupons/:id/copy", Handler: rid(gw(owner(couponMerchant.MerchantCopyCoupon)))},
		{Method: http.MethodPost, Path: "/api/v1/merchant/coupons/grant", Handler: rid(gw(owner(couponMerchant.MerchantGrantCoupon)))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/coupons/:id/claims", Handler: rid(gw(owner(couponMerchant.MerchantCouponClaims)))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/coupons/:id/redeems", Handler: rid(gw(owner(couponMerchant.MerchantCouponRedeems)))},
		{Method: http.MethodGet, Path: "/api/v1/merchant/coupons/:id/stats", Handler: rid(gw(owner(couponMerchant.MerchantCouponStats)))},

		{Method: http.MethodGet, Path: "/api/v1/admin/applications", Handler: rid(middleware.Chain(shopAdmin.AdminListApplications, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/applications/:id/approve", Handler: rid(middleware.Chain(shopAdmin.AdminApprove, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/applications/:id/reject", Handler: rid(middleware.Chain(shopAdmin.AdminReject, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/shops", Handler: rid(middleware.Chain(shopAdmin.AdminListShops, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/shops", Handler: rid(middleware.Chain(shopAdmin.AdminCreateShop, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/shops/:id", Handler: rid(middleware.Chain(shopAdmin.AdminGetShop, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/shops/:id", Handler: rid(middleware.Chain(shopAdmin.AdminUpdateShop, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/shops/:id/owner-password", Handler: rid(middleware.Chain(shopAdmin.AdminResetOwnerPassword, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/shops/:id/disable", Handler: rid(middleware.Chain(shopAdmin.AdminDisableShop, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/shops/:id/enable", Handler: rid(middleware.Chain(shopAdmin.AdminEnableShop, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/shops/:id/wallet", Handler: rid(middleware.Chain(walletAdmin.AdminGetWallet, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/shops/:id/wallet/adjust", Handler: rid(middleware.Chain(walletAdmin.AdminAdjustWallet, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/shops/:id/wallet/logs", Handler: rid(middleware.Chain(walletAdmin.AdminWalletLogs, gw, plat))},

		{Method: http.MethodGet, Path: "/api/v1/admin/seckill/rule", Handler: rid(middleware.Chain(seckillAdmin.AdminGetSeckillRule, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/seckill/rule", Handler: rid(middleware.Chain(seckillAdmin.AdminUpdateSeckillRule, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/seckill/sessions", Handler: rid(middleware.Chain(seckillAdmin.AdminListSeckillSessions, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/seckill/entries", Handler: rid(middleware.Chain(seckillAdmin.AdminListSeckillEntries, gw, plat))},

		{Method: http.MethodGet, Path: "/api/v1/admin/homepage-packages", Handler: rid(middleware.Chain(slotAdmin.AdminListSlotPackages, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/homepage-packages", Handler: rid(middleware.Chain(slotAdmin.AdminCreateSlotPackage, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/homepage-packages/:id", Handler: rid(middleware.Chain(slotAdmin.AdminUpdateSlotPackage, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/homepage-settings", Handler: rid(middleware.Chain(slotAdmin.AdminListSlotSettings, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/homepage-settings", Handler: rid(middleware.Chain(slotAdmin.AdminUpdateSlotSettings, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/homepage-orders", Handler: rid(middleware.Chain(slotAdmin.AdminListSlotOrders, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/homepage-orders/grant", Handler: rid(middleware.Chain(slotAdmin.AdminGrantSlot, gw, plat))},

		{Method: http.MethodGet, Path: "/api/v1/admin/theme-slots", Handler: rid(middleware.Chain(themeAdmin.AdminListThemeSlots, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/theme-slots/:id", Handler: rid(middleware.Chain(themeAdmin.AdminUpdateThemeSlot, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/theme-packages", Handler: rid(middleware.Chain(themeAdmin.AdminListThemePackages, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/theme-packages", Handler: rid(middleware.Chain(themeAdmin.AdminCreateThemePackage, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/theme-packages/:id", Handler: rid(middleware.Chain(themeAdmin.AdminUpdateThemePackage, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/theme-orders", Handler: rid(middleware.Chain(themeAdmin.AdminListThemeOrders, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/theme-orders/grant", Handler: rid(middleware.Chain(themeAdmin.AdminGrantTheme, gw, plat))},

		{Method: http.MethodGet, Path: "/api/v1/admin/coupons", Handler: rid(middleware.Chain(couponAdmin.AdminListCoupons, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/coupons", Handler: rid(middleware.Chain(couponAdmin.AdminCreateCoupon, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/coupons/:id", Handler: rid(middleware.Chain(couponAdmin.AdminUpdateCoupon, gw, plat))},
		{Method: http.MethodPut, Path: "/api/v1/admin/coupons/:id/off", Handler: rid(middleware.Chain(couponAdmin.AdminOffCoupon, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/coupons/:id/copy", Handler: rid(middleware.Chain(couponAdmin.AdminCopyCoupon, gw, plat))},
		{Method: http.MethodPost, Path: "/api/v1/admin/coupons/grant", Handler: rid(middleware.Chain(couponAdmin.AdminGrantCoupon, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/coupons/:id/claims", Handler: rid(middleware.Chain(couponAdmin.AdminCouponClaims, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/coupons/:id/redeems", Handler: rid(middleware.Chain(couponAdmin.AdminCouponRedeems, gw, plat))},
		{Method: http.MethodGet, Path: "/api/v1/admin/coupons/:id/stats", Handler: rid(middleware.Chain(couponAdmin.AdminCouponStats, gw, plat))},

		{Method: http.MethodGet, Path: "/api/v1/admin/points-products", Handler: admin(pointsProductAdmin.List)},
		{Method: http.MethodPost, Path: "/api/v1/admin/points-products", Handler: admin(pointsProductAdmin.Create)},
		{Method: http.MethodPost, Path: "/api/v1/admin/points-products/upload", Handler: admin(pointsProductAdmin.Upload)},
		{Method: http.MethodGet, Path: "/api/v1/admin/points-products/:id", Handler: admin(pointsProductAdmin.Detail)},
		{Method: http.MethodPut, Path: "/api/v1/admin/points-products/:id", Handler: admin(pointsProductAdmin.Update)},
		{Method: http.MethodPut, Path: "/api/v1/admin/points-products/:id/status", Handler: admin(pointsProductAdmin.SetStatus)},
		{Method: http.MethodDelete, Path: "/api/v1/admin/points-products/:id", Handler: admin(pointsProductAdmin.Delete)},

		{Method: http.MethodGet, Path: "/api/v1/admin/points-orders", Handler: admin(pointsOrderAdmin.List)},
		{Method: http.MethodGet, Path: "/api/v1/admin/points-orders/:id", Handler: admin(pointsOrderAdmin.Detail)},
		{Method: http.MethodPost, Path: "/api/v1/admin/points-orders/:id/ship", Handler: admin(pointsOrderAdmin.Ship)},
		{Method: http.MethodPost, Path: "/api/v1/admin/points-orders/:id/complete", Handler: admin(pointsOrderAdmin.Complete)},
		{Method: http.MethodPost, Path: "/api/v1/admin/points-orders/:id/cancel", Handler: admin(pointsOrderAdmin.Cancel)},
		{Method: http.MethodPut, Path: "/api/v1/admin/points-orders/:id/remark", Handler: admin(pointsOrderAdmin.Remark)},

		{Method: http.MethodPost, Path: "/api/v1/user/points-mall/exchange", Handler: userAuth(pointsOrderUser.Exchange)},
		{Method: http.MethodGet, Path: "/api/v1/user/points-mall/orders", Handler: userAuth(pointsOrderUser.List)},
		{Method: http.MethodGet, Path: "/api/v1/user/points-mall/orders/:id", Handler: userAuth(pointsOrderUser.Detail)},

		{Method: http.MethodGet, Path: "/uploads/points-mall/:file", Handler: rid(func(w http.ResponseWriter, r *http.Request) {
			p := uploadpath.Abs("points-mall", httpserver.PathParam(r, "file"))
			http.ServeFile(w, r, p)
		})},
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
