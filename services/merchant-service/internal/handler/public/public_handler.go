package public

import (
	"net/http"

	"mymall/services/merchant-service/internal/logic/public"
	"mymall/services/merchant-service/internal/svc"
)

func HealthzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewHealthzLogic(r.Context(), svcCtx)
		l.Healthz(w, r)
	}
}

func PublicCouponCenterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewPublicCouponCenterLogic(r.Context(), svcCtx)
		l.PublicCouponCenter(w, r)
	}
}

func PublicCouponPopupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewPublicCouponPopupLogic(r.Context(), svcCtx)
		l.PublicCouponPopup(w, r)
	}
}

func PublicGetShopHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewPublicGetShopLogic(r.Context(), svcCtx)
		l.PublicGetShop(w, r)
	}
}

func PublicHomeSlotsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewPublicHomeSlotsLogic(r.Context(), svcCtx)
		l.PublicHomeSlots(w, r)
	}
}

func PublicListShopsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewPublicListShopsLogic(r.Context(), svcCtx)
		l.PublicListShops(w, r)
	}
}

func PublicSeckillCurrentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewPublicSeckillCurrentLogic(r.Context(), svcCtx)
		l.PublicSeckillCurrent(w, r)
	}
}

func PublicSeckillEntryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewPublicSeckillEntryLogic(r.Context(), svcCtx)
		l.PublicSeckillEntry(w, r)
	}
}

func PublicSeckillListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewPublicSeckillListLogic(r.Context(), svcCtx)
		l.PublicSeckillList(w, r)
	}
}

func PublicThemeTilesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewPublicThemeTilesLogic(r.Context(), svcCtx)
		l.PublicThemeTiles(w, r)
	}
}

func ReadyzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewReadyzLogic(r.Context(), svcCtx)
		l.Readyz(w, r)
	}
}

func ServePointsMallUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewServePointsMallUploadLogic(r.Context(), svcCtx)
		l.ServePointsMallUpload(w, r)
	}
}
