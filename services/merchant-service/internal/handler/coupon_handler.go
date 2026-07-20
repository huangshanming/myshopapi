package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/logic"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *CouponAdminHandler) AdminListCoupons(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	list, total, err := h.logic.ListCoupons("platform", 0, r.URL.Query().Get("status"), r.URL.Query().Get("keyword"), page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}

func (h *CouponAdminHandler) AdminCreateCoupon(w http.ResponseWriter, r *http.Request) {
	adminID, _ := middleware.GetUserID(r.Context())
	var req logic.CouponSaveReq
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

func (h *CouponAdminHandler) AdminUpdateCoupon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var req logic.CouponSaveReq
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

func (h *CouponAdminHandler) AdminOffCoupon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.OffCoupon(id, 0, true); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CouponAdminHandler) AdminCopyCoupon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	adminID, _ := middleware.GetUserID(r.Context())
	c, err := h.logic.CopyCoupon(id, 0, adminID, true)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, c)
}

func (h *CouponAdminHandler) AdminGrantCoupon(w http.ResponseWriter, r *http.Request) {
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

func (h *CouponAdminHandler) AdminCouponClaims(w http.ResponseWriter, r *http.Request) {
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

func (h *CouponAdminHandler) AdminCouponRedeems(w http.ResponseWriter, r *http.Request) {
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

func (h *CouponAdminHandler) AdminCouponStats(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	st, err := h.logic.CouponStats(id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, st)
}

func (h *CouponMerchantHandler) MerchantListCoupons(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	list, total, err := h.logic.ListCoupons("shop", shopID, r.URL.Query().Get("status"), r.URL.Query().Get("keyword"), page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}

func (h *CouponMerchantHandler) MerchantCreateCoupon(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())
	var req logic.CouponSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	c, err := h.logic.MerchantCreateCoupon(shopID, userID, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, c)
}

func (h *CouponMerchantHandler) MerchantUpdateCoupon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	shopID := middleware.GetShopID(r.Context())
	var req logic.CouponSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UpdateCoupon(id, shopID, false, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CouponMerchantHandler) MerchantOffCoupon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	shopID := middleware.GetShopID(r.Context())
	if err := h.logic.OffCoupon(id, shopID, false); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CouponMerchantHandler) MerchantCopyCoupon(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	shopID := middleware.GetShopID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())
	c, err := h.logic.CopyCoupon(id, shopID, userID, false)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, c)
}

func (h *CouponMerchantHandler) MerchantGrantCoupon(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())
	var body struct {
		CouponID uint64   `json:"coupon_id"`
		UserIDs  []uint64 `json:"user_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	g, err := h.logic.GrantCoupon(userID, body.CouponID, body.UserIDs, shopID, false)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, g)
}

func (h *CouponMerchantHandler) MerchantCouponClaims(w http.ResponseWriter, r *http.Request) {
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
func (h *CouponMerchantHandler) MerchantCouponRedeems(w http.ResponseWriter, r *http.Request) {
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
func (h *CouponMerchantHandler) MerchantCouponStats(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	st, err := h.logic.CouponStats(id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, st)
}

func (h *CouponPublicHandler) PublicCouponCenter(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())
	if userID == 0 {
		if raw := r.Header.Get("X-User-Id"); raw != "" {
			userID, _ = strconv.ParseUint(raw, 10, 64)
		}
	}
	shopID, _ := strconv.ParseUint(r.URL.Query().Get("shop_id"), 10, 64)
	list, err := h.logic.ListCenter(userID, shopID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list})
}

func (h *CouponPublicHandler) PublicCouponPopup(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())
	list, err := h.logic.ListPopup(userID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list})
}

func (h *CouponUserHandler) ClaimCoupon(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())
	if userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "请先登录"))
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	var body struct {
		Source string `json:"source"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	uc, err := h.logic.ClaimCoupon(userID, id, body.Source)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, uc)
}

func (h *CouponUserHandler) ListMyCoupons(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())
	if userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "请先登录"))
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	list, total, err := h.logic.ListMyCoupons(userID, r.URL.Query().Get("status"), page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}

func (h *CouponInternalHandler) InternalMatchCoupons(w http.ResponseWriter, r *http.Request) {
	var req logic.MatchReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	resp, err := h.logic.MatchCoupons(req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, resp)
}

func (h *CouponInternalHandler) InternalLockCoupon(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserCouponID   uint64  `json:"user_coupon_id"`
		UserID         uint64  `json:"user_id"`
		OrderID        uint64  `json:"order_id"`
		DiscountAmount float64 `json:"discount_amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.LockCoupon(body.UserCouponID, body.UserID, body.OrderID, body.DiscountAmount); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CouponInternalHandler) InternalUnlockCoupon(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserCouponID uint64 `json:"user_coupon_id"`
		OrderID      uint64 `json:"order_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UnlockCoupon(body.UserCouponID, body.OrderID); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CouponInternalHandler) InternalRedeemCoupon(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserCouponID   uint64  `json:"user_coupon_id"`
		OrderID        uint64  `json:"order_id"`
		DiscountAmount float64 `json:"discount_amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.RedeemCoupon(body.UserCouponID, body.OrderID, body.DiscountAmount); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CouponInternalHandler) InternalReturnCoupon(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserCouponID uint64 `json:"user_coupon_id"`
		OrderID      uint64 `json:"order_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.ReturnCoupon(body.UserCouponID, body.OrderID); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *CouponInternalHandler) InternalOrderGift(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserID uint64 `json:"user_id"`
		ShopID uint64 `json:"shop_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	n, err := h.logic.OrderGiftCoupons(body.UserID, body.ShopID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"granted": n})
}
