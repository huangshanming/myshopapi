package internalapi

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/pkg/xerr"
)

func (h *CouponHandler) InternalMatchCoupons(ctx context.Context, in appinput.CallInput) (any, error) {
	var req biz.MatchReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	resp, err := h.logic.MatchCoupons(req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return resp, nil
}

func (h *CouponHandler) InternalLockCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	var body struct {
		UserCouponID   uint64  `json:"user_coupon_id"`
		UserID         uint64  `json:"user_id"`
		OrderID        uint64  `json:"order_id"`
		DiscountAmount float64 `json:"discount_amount"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.LockCoupon(body.UserCouponID, body.UserID, body.OrderID, body.DiscountAmount); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *CouponHandler) InternalUnlockCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	var body struct {
		UserCouponID uint64 `json:"user_coupon_id"`
		OrderID      uint64 `json:"order_id"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UnlockCoupon(body.UserCouponID, body.OrderID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *CouponHandler) InternalRedeemCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	var body struct {
		UserCouponID   uint64  `json:"user_coupon_id"`
		OrderID        uint64  `json:"order_id"`
		DiscountAmount float64 `json:"discount_amount"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.RedeemCoupon(body.UserCouponID, body.OrderID, body.DiscountAmount); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *CouponHandler) InternalReturnCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	var body struct {
		UserCouponID uint64 `json:"user_coupon_id"`
		OrderID      uint64 `json:"order_id"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.ReturnCoupon(body.UserCouponID, body.OrderID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *CouponHandler) InternalOrderGift(ctx context.Context, in appinput.CallInput) (any, error) {
	var body struct {
		UserID uint64 `json:"user_id"`
		ShopID uint64 `json:"shop_id"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	n, err := h.logic.OrderGiftCoupons(body.UserID, body.ShopID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return map[string]interface{}{"granted": n}, nil
}
