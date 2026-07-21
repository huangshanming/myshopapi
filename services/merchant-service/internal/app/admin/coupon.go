package admin

import (
	"encoding/json"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
)

func (h *CouponHandler) AdminListCoupons(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	list, total, err := h.logic.ListCoupons("platform", 0, r.URL.Query().Get("status"), r.URL.Query().Get("keyword"), page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}

func (h *CouponHandler) AdminCreateCoupon(w http.ResponseWriter, r *http.Request) {
	adminID, _ := middleware.GetUserID(r.Context())
	var req biz.CouponSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	c, err := h.logic.AdminCreateCoupon(adminID, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, c)
}

func (h *CouponHandler) AdminUpdateCoupon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req biz.CouponSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UpdateCoupon(id, 0, true, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CouponHandler) AdminOffCoupon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.OffCoupon(id, 0, true); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CouponHandler) AdminCopyCoupon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	adminID, _ := middleware.GetUserID(r.Context())
	c, err := h.logic.CopyCoupon(id, 0, adminID, true)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, c)
}

func (h *CouponHandler) AdminGrantCoupon(w http.ResponseWriter, r *http.Request) {
	adminID, _ := middleware.GetUserID(r.Context())
	var body struct {
		CouponID uint64   `json:"coupon_id"`
		UserIDs  []uint64 `json:"user_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	g, err := h.logic.GrantCoupon(adminID, body.CouponID, body.UserIDs, 0, true)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, g)
}

func (h *CouponHandler) AdminCouponClaims(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	list, total, err := h.logic.CouponClaims(id, page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}

func (h *CouponHandler) AdminCouponRedeems(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	list, total, err := h.logic.CouponRedeems(id, page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}

func (h *CouponHandler) AdminCouponStats(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	st, err := h.logic.CouponStats(id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, st)
}
