package merchant

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"strconv"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
)

func (h *CouponHandler) MerchantListCoupons(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	list, total, err := h.logic.ListCoupons("shop", shopID, in.QueryGet("status"), in.QueryGet("keyword"), page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *CouponHandler) MerchantCreateCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	userID, _ := middleware.GetUserID(ctx)
	var req biz.CouponSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	c, err := h.logic.MerchantCreateCoupon(shopID, userID, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return c, nil
}

func (h *CouponHandler) MerchantUpdateCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	shopID := middleware.GetShopID(ctx)
	var req biz.CouponSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UpdateCoupon(id, shopID, false, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *CouponHandler) MerchantOffCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	shopID := middleware.GetShopID(ctx)
	if err := h.logic.OffCoupon(id, shopID, false); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *CouponHandler) MerchantCopyCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	shopID := middleware.GetShopID(ctx)
	userID, _ := middleware.GetUserID(ctx)
	c, err := h.logic.CopyCoupon(id, shopID, userID, false)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return c, nil
}

func (h *CouponHandler) MerchantGrantCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	userID, _ := middleware.GetUserID(ctx)
	var body struct {
		CouponID uint64   `json:"coupon_id"`
		UserIDs  []uint64 `json:"user_ids"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	g, err := h.logic.GrantCoupon(userID, body.CouponID, body.UserIDs, shopID, false)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return g, nil
}

func (h *CouponHandler) MerchantCouponClaims(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	list, total, err := h.logic.CouponClaims(id, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *CouponHandler) MerchantCouponRedeems(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	list, total, err := h.logic.CouponRedeems(id, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *CouponHandler) MerchantCouponStats(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	st, err := h.logic.CouponStats(id)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return st, nil
}
