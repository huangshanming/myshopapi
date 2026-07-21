package coupon

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/services/merchant-service/internal/logic/public/coupon"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
)

func PublicCouponCenterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShopIdQueryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := coupon.NewPublicCouponCenterLogic(r.Context(), svcCtx)
		resp, err := l.PublicCouponCenter(r.Context(), &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func PublicCouponPopupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := coupon.NewPublicCouponPopupLogic(r.Context(), svcCtx)
		resp, err := l.PublicCouponPopup(r.Context())
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
