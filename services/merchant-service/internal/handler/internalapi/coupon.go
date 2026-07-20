package internalapi

import (
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/logic"
)

func (h *CouponHandler) InternalMatchCoupons(w http.ResponseWriter, r *http.Request) {
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

func (h *CouponHandler) InternalLockCoupon(w http.ResponseWriter, r *http.Request) {
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

func (h *CouponHandler) InternalUnlockCoupon(w http.ResponseWriter, r *http.Request) {
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

func (h *CouponHandler) InternalRedeemCoupon(w http.ResponseWriter, r *http.Request) {
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

func (h *CouponHandler) InternalReturnCoupon(w http.ResponseWriter, r *http.Request) {
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

func (h *CouponHandler) InternalOrderGift(w http.ResponseWriter, r *http.Request) {
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
