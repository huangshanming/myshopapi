package admin

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/admin"
	"mymall/services/merchant-service/internal/svc"
)

func AdminAdjustWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminAdjustWalletLogic(r.Context(), svcCtx)
		l.AdminAdjustWallet(w, r)
	}
}

func AdminApproveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminApproveLogic(r.Context(), svcCtx)
		l.AdminApprove(w, r)
	}
}

func AdminCopyCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminCopyCouponLogic(r.Context(), svcCtx)
		l.AdminCopyCoupon(w, r)
	}
}

func AdminCouponClaimsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminCouponClaimsLogic(r.Context(), svcCtx)
		l.AdminCouponClaims(w, r)
	}
}

func AdminCouponRedeemsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminCouponRedeemsLogic(r.Context(), svcCtx)
		l.AdminCouponRedeems(w, r)
	}
}

func AdminCouponStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminCouponStatsLogic(r.Context(), svcCtx)
		l.AdminCouponStats(w, r)
	}
}

func AdminCreateCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminCreateCouponLogic(r.Context(), svcCtx)
		l.AdminCreateCoupon(w, r)
	}
}

func AdminCreateShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminCreateShopLogic(r.Context(), svcCtx)
		l.AdminCreateShop(w, r)
	}
}

func AdminCreateSlotPackageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminCreateSlotPackageLogic(r.Context(), svcCtx)
		l.AdminCreateSlotPackage(w, r)
	}
}

func AdminCreateThemePackageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminCreateThemePackageLogic(r.Context(), svcCtx)
		l.AdminCreateThemePackage(w, r)
	}
}

func AdminDisableShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminDisableShopLogic(r.Context(), svcCtx)
		l.AdminDisableShop(w, r)
	}
}

func AdminEnableShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminEnableShopLogic(r.Context(), svcCtx)
		l.AdminEnableShop(w, r)
	}
}

func AdminGetSeckillRuleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminGetSeckillRuleLogic(r.Context(), svcCtx)
		l.AdminGetSeckillRule(w, r)
	}
}

func AdminGetShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminGetShopLogic(r.Context(), svcCtx)
		l.AdminGetShop(w, r)
	}
}

func AdminGetWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminGetWalletLogic(r.Context(), svcCtx)
		l.AdminGetWallet(w, r)
	}
}

func AdminGrantCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminGrantCouponLogic(r.Context(), svcCtx)
		l.AdminGrantCoupon(w, r)
	}
}

func AdminGrantSlotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminGrantSlotLogic(r.Context(), svcCtx)
		l.AdminGrantSlot(w, r)
	}
}

func AdminGrantThemeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminGrantThemeLogic(r.Context(), svcCtx)
		l.AdminGrantTheme(w, r)
	}
}

func AdminListApplicationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListApplicationsLogic(r.Context(), svcCtx)
		l.AdminListApplications(w, r)
	}
}

func AdminListCouponsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListCouponsLogic(r.Context(), svcCtx)
		l.AdminListCoupons(w, r)
	}
}

func AdminListSeckillEntriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListSeckillEntriesLogic(r.Context(), svcCtx)
		l.AdminListSeckillEntries(w, r)
	}
}

func AdminListSeckillSessionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListSeckillSessionsLogic(r.Context(), svcCtx)
		l.AdminListSeckillSessions(w, r)
	}
}

func AdminListShopsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListShopsLogic(r.Context(), svcCtx)
		l.AdminListShops(w, r)
	}
}

func AdminListSlotOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListSlotOrdersLogic(r.Context(), svcCtx)
		l.AdminListSlotOrders(w, r)
	}
}

func AdminListSlotPackagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListSlotPackagesLogic(r.Context(), svcCtx)
		l.AdminListSlotPackages(w, r)
	}
}

func AdminListSlotSettingsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListSlotSettingsLogic(r.Context(), svcCtx)
		l.AdminListSlotSettings(w, r)
	}
}

func AdminListThemeOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListThemeOrdersLogic(r.Context(), svcCtx)
		l.AdminListThemeOrders(w, r)
	}
}

func AdminListThemePackagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListThemePackagesLogic(r.Context(), svcCtx)
		l.AdminListThemePackages(w, r)
	}
}

func AdminListThemeSlotsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminListThemeSlotsLogic(r.Context(), svcCtx)
		l.AdminListThemeSlots(w, r)
	}
}

func AdminOffCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminOffCouponLogic(r.Context(), svcCtx)
		l.AdminOffCoupon(w, r)
	}
}

func AdminRejectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminRejectLogic(r.Context(), svcCtx)
		l.AdminReject(w, r)
	}
}

func AdminResetOwnerPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminResetOwnerPasswordLogic(r.Context(), svcCtx)
		l.AdminResetOwnerPassword(w, r)
	}
}

func AdminUpdateCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminUpdateCouponLogic(r.Context(), svcCtx)
		l.AdminUpdateCoupon(w, r)
	}
}

func AdminUpdateSeckillRuleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminUpdateSeckillRuleLogic(r.Context(), svcCtx)
		l.AdminUpdateSeckillRule(w, r)
	}
}

func AdminUpdateShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminUpdateShopLogic(r.Context(), svcCtx)
		l.AdminUpdateShop(w, r)
	}
}

func AdminUpdateSlotPackageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminUpdateSlotPackageLogic(r.Context(), svcCtx)
		l.AdminUpdateSlotPackage(w, r)
	}
}

func AdminUpdateSlotSettingsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminUpdateSlotSettingsLogic(r.Context(), svcCtx)
		l.AdminUpdateSlotSettings(w, r)
	}
}

func AdminUpdateThemePackageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminUpdateThemePackageLogic(r.Context(), svcCtx)
		l.AdminUpdateThemePackage(w, r)
	}
}

func AdminUpdateThemeSlotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminUpdateThemeSlotLogic(r.Context(), svcCtx)
		l.AdminUpdateThemeSlot(w, r)
	}
}

func AdminWalletLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewAdminWalletLogsLogic(r.Context(), svcCtx)
		l.AdminWalletLogs(w, r)
	}
}

func CancelPointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewCancelPointsOrderLogic(r.Context(), svcCtx)
		l.CancelPointsOrder(w, r)
	}
}

func CompletePointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewCompletePointsOrderLogic(r.Context(), svcCtx)
		l.CompletePointsOrder(w, r)
	}
}

func CreatePointsProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewCreatePointsProductLogic(r.Context(), svcCtx)
		l.CreatePointsProduct(w, r)
	}
}

func DeletePointsProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewDeletePointsProductLogic(r.Context(), svcCtx)
		l.DeletePointsProduct(w, r)
	}
}

func DetailPointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewDetailPointsOrderLogic(r.Context(), svcCtx)
		l.DetailPointsOrder(w, r)
	}
}

func DetailPointsProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewDetailPointsProductLogic(r.Context(), svcCtx)
		l.DetailPointsProduct(w, r)
	}
}

func ListPointsOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewListPointsOrdersLogic(r.Context(), svcCtx)
		l.ListPointsOrders(w, r)
	}
}

func ListPointsProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewListPointsProductsLogic(r.Context(), svcCtx)
		l.ListPointsProducts(w, r)
	}
}

func RemarkPointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewRemarkPointsOrderLogic(r.Context(), svcCtx)
		l.RemarkPointsOrder(w, r)
	}
}

func SetPointsProductStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewSetPointsProductStatusLogic(r.Context(), svcCtx)
		l.SetPointsProductStatus(w, r)
	}
}

func ShipPointsOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewShipPointsOrderLogic(r.Context(), svcCtx)
		l.ShipPointsOrder(w, r)
	}
}

func UpdatePointsProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewUpdatePointsProductLogic(r.Context(), svcCtx)
		l.UpdatePointsProduct(w, r)
	}
}

func UploadPointsProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := admin.NewUploadPointsProductLogic(r.Context(), svcCtx)
		l.UploadPointsProduct(w, r)
	}
}
