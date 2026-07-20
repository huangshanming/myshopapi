package handler

import (
	"net/http"

	"mymall/pkg/health"
	"mymall/pkg/httpserver"
	hadmin "mymall/services/merchant-service/internal/handler/admin"
	hinternal "mymall/services/merchant-service/internal/handler/internalapi"
	hmerchant "mymall/services/merchant-service/internal/handler/merchant"
	hpublic "mymall/services/merchant-service/internal/handler/public"
	huser "mymall/services/merchant-service/internal/handler/user"
	svcMW "mymall/services/merchant-service/internal/middleware"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/uploadpath"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterHandlers wires routes declared in api/merchant.api.
func RegisterHandlers(server *rest.Server, svcCtx *svc.ServiceContext, healthReg *health.Registry, mws svcMW.Bundle) {
	shopAdmin := hadmin.NewShopHandler(svcCtx)
	shopMerchant := hmerchant.NewShopHandler(svcCtx)
	shopPublic := hpublic.NewShopHandler(svcCtx)
	walletAdmin := hadmin.NewWalletHandler(svcCtx)
	walletMerchant := hmerchant.NewWalletHandler(svcCtx)
	seckillAdmin := hadmin.NewSeckillHandler(svcCtx)
	seckillMerchant := hmerchant.NewSeckillHandler(svcCtx)
	seckillPublic := hpublic.NewSeckillHandler(svcCtx)
	seckillInternal := hinternal.NewSeckillHandler(svcCtx)
	couponAdmin := hadmin.NewCouponHandler(svcCtx)
	couponMerchant := hmerchant.NewCouponHandler(svcCtx)
	couponPublic := hpublic.NewCouponHandler(svcCtx)
	couponUser := huser.NewCouponHandler(svcCtx)
	couponInternal := hinternal.NewCouponHandler(svcCtx)
	slotAdmin := hadmin.NewHomepageSlotHandler(svcCtx)
	slotMerchant := hmerchant.NewHomepageSlotHandler(svcCtx)
	slotPublic := hpublic.NewHomepageSlotHandler(svcCtx)
	themeAdmin := hadmin.NewHomepageThemeHandler(svcCtx)
	themeMerchant := hmerchant.NewHomepageThemeHandler(svcCtx)
	themePublic := hpublic.NewHomepageThemeHandler(svcCtx)
	pointsProductAdmin := hadmin.NewPointsProductHandler(svcCtx)
	pointsOrderAdmin := hadmin.NewPointsOrderHandler(svcCtx)
	pointsOrderUser := huser.NewPointsOrderHandler(svcCtx)

	server.AddRoutes(mws.Public([]rest.Route{
		{Method: http.MethodGet, Path: "/healthz", Handler: httpserver.Healthz("merchant-service")},
		{Method: http.MethodGet, Path: "/readyz", Handler: healthReg.ReadyHandler()},
		{Method: http.MethodGet, Path: "/api/v1/shops/list", Handler: shopPublic.PublicListShops},
		{Method: http.MethodGet, Path: "/api/v1/shops/home-slots", Handler: slotPublic.PublicHomeSlots},
		{Method: http.MethodGet, Path: "/api/v1/home/theme-tiles", Handler: themePublic.PublicThemeTiles},
		{Method: http.MethodGet, Path: "/api/v1/shops/:id", Handler: shopPublic.PublicGetShop},
		{Method: http.MethodGet, Path: "/api/v1/seckill/current", Handler: seckillPublic.PublicSeckillCurrent},
		{Method: http.MethodGet, Path: "/api/v1/seckill/list", Handler: seckillPublic.PublicSeckillList},
		{Method: http.MethodGet, Path: "/api/v1/seckill/entries/:id", Handler: seckillPublic.PublicSeckillEntry},
		{Method: http.MethodPost, Path: "/api/v1/seckill/consume", Handler: seckillInternal.SeckillConsume},
		{Method: http.MethodPost, Path: "/api/v1/seckill/restore", Handler: seckillInternal.SeckillRestore},
		{Method: http.MethodGet, Path: "/api/v1/coupons/center", Handler: couponPublic.PublicCouponCenter},
		{Method: http.MethodGet, Path: "/api/v1/coupons/popup", Handler: couponPublic.PublicCouponPopup},
		{Method: http.MethodPost, Path: "/api/v1/internal/coupons/match", Handler: couponInternal.InternalMatchCoupons},
		{Method: http.MethodPost, Path: "/api/v1/internal/coupons/lock", Handler: couponInternal.InternalLockCoupon},
		{Method: http.MethodPost, Path: "/api/v1/internal/coupons/unlock", Handler: couponInternal.InternalUnlockCoupon},
		{Method: http.MethodPost, Path: "/api/v1/internal/coupons/redeem", Handler: couponInternal.InternalRedeemCoupon},
		{Method: http.MethodPost, Path: "/api/v1/internal/coupons/return", Handler: couponInternal.InternalReturnCoupon},
		{Method: http.MethodPost, Path: "/api/v1/internal/coupons/order-gift", Handler: couponInternal.InternalOrderGift},
		{Method: http.MethodGet, Path: "/uploads/points-mall/:file", Handler: servePointsMallUpload},
	}))

	server.AddRoutes(mws.Authed([]rest.Route{
		{Method: http.MethodPost, Path: "/api/v1/coupons/:id/claim", Handler: couponUser.ClaimCoupon},
		{Method: http.MethodGet, Path: "/api/v1/user/coupons", Handler: couponUser.ListMyCoupons},
		{Method: http.MethodPost, Path: "/api/v1/merchant/apply", Handler: shopMerchant.Apply},
		{Method: http.MethodGet, Path: "/api/v1/merchant/shops", Handler: shopMerchant.MyShops},
		{Method: http.MethodPost, Path: "/api/v1/user/points-mall/exchange", Handler: pointsOrderUser.Exchange},
		{Method: http.MethodGet, Path: "/api/v1/user/points-mall/orders", Handler: pointsOrderUser.List},
		{Method: http.MethodGet, Path: "/api/v1/user/points-mall/orders/:id", Handler: pointsOrderUser.Detail},
	}))

	server.AddRoutes(mws.MerchantOwner([]rest.Route{
		{Method: http.MethodPut, Path: "/api/v1/merchant/shops/:id", Handler: shopMerchant.UpdateMyShop},
		{Method: http.MethodGet, Path: "/api/v1/merchant/wallet", Handler: walletMerchant.MerchantGetWallet},
		{Method: http.MethodGet, Path: "/api/v1/merchant/wallet/logs", Handler: walletMerchant.MerchantWalletLogs},
		{Method: http.MethodGet, Path: "/api/v1/merchant/seckill/sessions", Handler: seckillMerchant.MerchantSeckillSessions},
		{Method: http.MethodPost, Path: "/api/v1/merchant/seckill/entries", Handler: seckillMerchant.MerchantApplySeckill},
		{Method: http.MethodGet, Path: "/api/v1/merchant/seckill/entries", Handler: seckillMerchant.MerchantListSeckillEntries},
		{Method: http.MethodPut, Path: "/api/v1/merchant/seckill/entries/:id/auto-renew", Handler: seckillMerchant.MerchantSetSeckillAutoRenew},
		{Method: http.MethodGet, Path: "/api/v1/merchant/homepage-packages", Handler: slotMerchant.MerchantListSlotPackages},
		{Method: http.MethodPost, Path: "/api/v1/merchant/homepage-orders", Handler: slotMerchant.MerchantBuySlot},
		{Method: http.MethodGet, Path: "/api/v1/merchant/homepage-orders", Handler: slotMerchant.MerchantListSlotOrders},
		{Method: http.MethodGet, Path: "/api/v1/merchant/theme-slots", Handler: themeMerchant.MerchantListThemeSlots},
		{Method: http.MethodGet, Path: "/api/v1/merchant/theme-packages", Handler: themeMerchant.MerchantListThemePackages},
		{Method: http.MethodPost, Path: "/api/v1/merchant/theme-orders", Handler: themeMerchant.MerchantBuyTheme},
		{Method: http.MethodGet, Path: "/api/v1/merchant/theme-orders", Handler: themeMerchant.MerchantListThemeOrders},
		{Method: http.MethodGet, Path: "/api/v1/merchant/coupons", Handler: couponMerchant.MerchantListCoupons},
		{Method: http.MethodPost, Path: "/api/v1/merchant/coupons", Handler: couponMerchant.MerchantCreateCoupon},
		{Method: http.MethodPut, Path: "/api/v1/merchant/coupons/:id", Handler: couponMerchant.MerchantUpdateCoupon},
		{Method: http.MethodPut, Path: "/api/v1/merchant/coupons/:id/off", Handler: couponMerchant.MerchantOffCoupon},
		{Method: http.MethodPost, Path: "/api/v1/merchant/coupons/:id/copy", Handler: couponMerchant.MerchantCopyCoupon},
		{Method: http.MethodPost, Path: "/api/v1/merchant/coupons/grant", Handler: couponMerchant.MerchantGrantCoupon},
		{Method: http.MethodGet, Path: "/api/v1/merchant/coupons/:id/claims", Handler: couponMerchant.MerchantCouponClaims},
		{Method: http.MethodGet, Path: "/api/v1/merchant/coupons/:id/redeems", Handler: couponMerchant.MerchantCouponRedeems},
		{Method: http.MethodGet, Path: "/api/v1/merchant/coupons/:id/stats", Handler: couponMerchant.MerchantCouponStats},
	}))

	server.AddRoutes(mws.PlatformAdmin([]rest.Route{
		{Method: http.MethodGet, Path: "/api/v1/admin/applications", Handler: shopAdmin.AdminListApplications},
		{Method: http.MethodPost, Path: "/api/v1/admin/applications/:id/approve", Handler: shopAdmin.AdminApprove},
		{Method: http.MethodPost, Path: "/api/v1/admin/applications/:id/reject", Handler: shopAdmin.AdminReject},
		{Method: http.MethodGet, Path: "/api/v1/admin/shops", Handler: shopAdmin.AdminListShops},
		{Method: http.MethodPost, Path: "/api/v1/admin/shops", Handler: shopAdmin.AdminCreateShop},
		{Method: http.MethodGet, Path: "/api/v1/admin/shops/:id", Handler: shopAdmin.AdminGetShop},
		{Method: http.MethodPut, Path: "/api/v1/admin/shops/:id", Handler: shopAdmin.AdminUpdateShop},
		{Method: http.MethodPut, Path: "/api/v1/admin/shops/:id/owner-password", Handler: shopAdmin.AdminResetOwnerPassword},
		{Method: http.MethodPut, Path: "/api/v1/admin/shops/:id/disable", Handler: shopAdmin.AdminDisableShop},
		{Method: http.MethodPut, Path: "/api/v1/admin/shops/:id/enable", Handler: shopAdmin.AdminEnableShop},
		{Method: http.MethodGet, Path: "/api/v1/admin/shops/:id/wallet", Handler: walletAdmin.AdminGetWallet},
		{Method: http.MethodPost, Path: "/api/v1/admin/shops/:id/wallet/adjust", Handler: walletAdmin.AdminAdjustWallet},
		{Method: http.MethodGet, Path: "/api/v1/admin/shops/:id/wallet/logs", Handler: walletAdmin.AdminWalletLogs},
		{Method: http.MethodGet, Path: "/api/v1/admin/seckill/rule", Handler: seckillAdmin.AdminGetSeckillRule},
		{Method: http.MethodPut, Path: "/api/v1/admin/seckill/rule", Handler: seckillAdmin.AdminUpdateSeckillRule},
		{Method: http.MethodGet, Path: "/api/v1/admin/seckill/sessions", Handler: seckillAdmin.AdminListSeckillSessions},
		{Method: http.MethodGet, Path: "/api/v1/admin/seckill/entries", Handler: seckillAdmin.AdminListSeckillEntries},
		{Method: http.MethodGet, Path: "/api/v1/admin/homepage-packages", Handler: slotAdmin.AdminListSlotPackages},
		{Method: http.MethodPost, Path: "/api/v1/admin/homepage-packages", Handler: slotAdmin.AdminCreateSlotPackage},
		{Method: http.MethodPut, Path: "/api/v1/admin/homepage-packages/:id", Handler: slotAdmin.AdminUpdateSlotPackage},
		{Method: http.MethodGet, Path: "/api/v1/admin/homepage-settings", Handler: slotAdmin.AdminListSlotSettings},
		{Method: http.MethodPut, Path: "/api/v1/admin/homepage-settings", Handler: slotAdmin.AdminUpdateSlotSettings},
		{Method: http.MethodGet, Path: "/api/v1/admin/homepage-orders", Handler: slotAdmin.AdminListSlotOrders},
		{Method: http.MethodPost, Path: "/api/v1/admin/homepage-orders/grant", Handler: slotAdmin.AdminGrantSlot},
		{Method: http.MethodGet, Path: "/api/v1/admin/theme-slots", Handler: themeAdmin.AdminListThemeSlots},
		{Method: http.MethodPut, Path: "/api/v1/admin/theme-slots/:id", Handler: themeAdmin.AdminUpdateThemeSlot},
		{Method: http.MethodGet, Path: "/api/v1/admin/theme-packages", Handler: themeAdmin.AdminListThemePackages},
		{Method: http.MethodPost, Path: "/api/v1/admin/theme-packages", Handler: themeAdmin.AdminCreateThemePackage},
		{Method: http.MethodPut, Path: "/api/v1/admin/theme-packages/:id", Handler: themeAdmin.AdminUpdateThemePackage},
		{Method: http.MethodGet, Path: "/api/v1/admin/theme-orders", Handler: themeAdmin.AdminListThemeOrders},
		{Method: http.MethodPost, Path: "/api/v1/admin/theme-orders/grant", Handler: themeAdmin.AdminGrantTheme},
		{Method: http.MethodGet, Path: "/api/v1/admin/coupons", Handler: couponAdmin.AdminListCoupons},
		{Method: http.MethodPost, Path: "/api/v1/admin/coupons", Handler: couponAdmin.AdminCreateCoupon},
		{Method: http.MethodPut, Path: "/api/v1/admin/coupons/:id", Handler: couponAdmin.AdminUpdateCoupon},
		{Method: http.MethodPut, Path: "/api/v1/admin/coupons/:id/off", Handler: couponAdmin.AdminOffCoupon},
		{Method: http.MethodPost, Path: "/api/v1/admin/coupons/:id/copy", Handler: couponAdmin.AdminCopyCoupon},
		{Method: http.MethodPost, Path: "/api/v1/admin/coupons/grant", Handler: couponAdmin.AdminGrantCoupon},
		{Method: http.MethodGet, Path: "/api/v1/admin/coupons/:id/claims", Handler: couponAdmin.AdminCouponClaims},
		{Method: http.MethodGet, Path: "/api/v1/admin/coupons/:id/redeems", Handler: couponAdmin.AdminCouponRedeems},
		{Method: http.MethodGet, Path: "/api/v1/admin/coupons/:id/stats", Handler: couponAdmin.AdminCouponStats},
		{Method: http.MethodGet, Path: "/api/v1/admin/points-products", Handler: pointsProductAdmin.List},
		{Method: http.MethodPost, Path: "/api/v1/admin/points-products", Handler: pointsProductAdmin.Create},
		{Method: http.MethodPost, Path: "/api/v1/admin/points-products/upload", Handler: pointsProductAdmin.Upload},
		{Method: http.MethodGet, Path: "/api/v1/admin/points-products/:id", Handler: pointsProductAdmin.Detail},
		{Method: http.MethodPut, Path: "/api/v1/admin/points-products/:id", Handler: pointsProductAdmin.Update},
		{Method: http.MethodPut, Path: "/api/v1/admin/points-products/:id/status", Handler: pointsProductAdmin.SetStatus},
		{Method: http.MethodDelete, Path: "/api/v1/admin/points-products/:id", Handler: pointsProductAdmin.Delete},
		{Method: http.MethodGet, Path: "/api/v1/admin/points-orders", Handler: pointsOrderAdmin.List},
		{Method: http.MethodGet, Path: "/api/v1/admin/points-orders/:id", Handler: pointsOrderAdmin.Detail},
		{Method: http.MethodPost, Path: "/api/v1/admin/points-orders/:id/ship", Handler: pointsOrderAdmin.Ship},
		{Method: http.MethodPost, Path: "/api/v1/admin/points-orders/:id/complete", Handler: pointsOrderAdmin.Complete},
		{Method: http.MethodPost, Path: "/api/v1/admin/points-orders/:id/cancel", Handler: pointsOrderAdmin.Cancel},
		{Method: http.MethodPut, Path: "/api/v1/admin/points-orders/:id/remark", Handler: pointsOrderAdmin.Remark},
	}))
}

func servePointsMallUpload(w http.ResponseWriter, r *http.Request) {
	p := uploadpath.Abs("points-mall", httpserver.PathParam(r, "file"))
	http.ServeFile(w, r, p)
}
