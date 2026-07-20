package merchant

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/merchant"
	"mymall/services/merchant-service/internal/svc"
)

func ApplyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewApplyLogic(r.Context(), svcCtx)
		l.Apply(w, r)
	}
}

func MerchantApplySeckillHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantApplySeckillLogic(r.Context(), svcCtx)
		l.MerchantApplySeckill(w, r)
	}
}

func MerchantBuySlotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantBuySlotLogic(r.Context(), svcCtx)
		l.MerchantBuySlot(w, r)
	}
}

func MerchantBuyThemeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantBuyThemeLogic(r.Context(), svcCtx)
		l.MerchantBuyTheme(w, r)
	}
}

func MerchantCopyCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantCopyCouponLogic(r.Context(), svcCtx)
		l.MerchantCopyCoupon(w, r)
	}
}

func MerchantCouponClaimsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantCouponClaimsLogic(r.Context(), svcCtx)
		l.MerchantCouponClaims(w, r)
	}
}

func MerchantCouponRedeemsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantCouponRedeemsLogic(r.Context(), svcCtx)
		l.MerchantCouponRedeems(w, r)
	}
}

func MerchantCouponStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantCouponStatsLogic(r.Context(), svcCtx)
		l.MerchantCouponStats(w, r)
	}
}

func MerchantCreateCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantCreateCouponLogic(r.Context(), svcCtx)
		l.MerchantCreateCoupon(w, r)
	}
}

func MerchantGetWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantGetWalletLogic(r.Context(), svcCtx)
		l.MerchantGetWallet(w, r)
	}
}

func MerchantGrantCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantGrantCouponLogic(r.Context(), svcCtx)
		l.MerchantGrantCoupon(w, r)
	}
}

func MerchantListCouponsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantListCouponsLogic(r.Context(), svcCtx)
		l.MerchantListCoupons(w, r)
	}
}

func MerchantListSeckillEntriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantListSeckillEntriesLogic(r.Context(), svcCtx)
		l.MerchantListSeckillEntries(w, r)
	}
}

func MerchantListSlotOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantListSlotOrdersLogic(r.Context(), svcCtx)
		l.MerchantListSlotOrders(w, r)
	}
}

func MerchantListSlotPackagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantListSlotPackagesLogic(r.Context(), svcCtx)
		l.MerchantListSlotPackages(w, r)
	}
}

func MerchantListThemeOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantListThemeOrdersLogic(r.Context(), svcCtx)
		l.MerchantListThemeOrders(w, r)
	}
}

func MerchantListThemePackagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantListThemePackagesLogic(r.Context(), svcCtx)
		l.MerchantListThemePackages(w, r)
	}
}

func MerchantListThemeSlotsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantListThemeSlotsLogic(r.Context(), svcCtx)
		l.MerchantListThemeSlots(w, r)
	}
}

func MerchantOffCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantOffCouponLogic(r.Context(), svcCtx)
		l.MerchantOffCoupon(w, r)
	}
}

func MerchantSeckillSessionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantSeckillSessionsLogic(r.Context(), svcCtx)
		l.MerchantSeckillSessions(w, r)
	}
}

func MerchantSetSeckillAutoRenewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantSetSeckillAutoRenewLogic(r.Context(), svcCtx)
		l.MerchantSetSeckillAutoRenew(w, r)
	}
}

func MerchantUpdateCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantUpdateCouponLogic(r.Context(), svcCtx)
		l.MerchantUpdateCoupon(w, r)
	}
}

func MerchantWalletLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMerchantWalletLogsLogic(r.Context(), svcCtx)
		l.MerchantWalletLogs(w, r)
	}
}

func MyShopsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewMyShopsLogic(r.Context(), svcCtx)
		l.MyShops(w, r)
	}
}

func UpdateMyShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewUpdateMyShopLogic(r.Context(), svcCtx)
		l.UpdateMyShop(w, r)
	}
}
