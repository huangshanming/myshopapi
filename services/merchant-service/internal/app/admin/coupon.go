package admin

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/services/merchant-service/internal/biz"
	"net/http"
	"strconv"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
)

func (h *CouponHandler) AdminListCoupons(ctx context.Context, in appinput.CallInput) (any, error) {
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	list, total, err := h.logic.ListCoupons("platform", 0, in.QueryGet("status"), in.QueryGet("keyword"), page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *CouponHandler) AdminCreateCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	adminID, _ := middleware.GetUserID(ctx)
	var req biz.CouponSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	c, err := h.logic.AdminCreateCoupon(adminID, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return c, nil
}

func (h *CouponHandler) AdminUpdateCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req biz.CouponSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UpdateCoupon(id, 0, true, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *CouponHandler) AdminOffCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	if err := h.logic.OffCoupon(id, 0, true); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *CouponHandler) AdminCopyCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	adminID, _ := middleware.GetUserID(ctx)
	c, err := h.logic.CopyCoupon(id, 0, adminID, true)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return c, nil
}

func (h *CouponHandler) AdminGrantCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	adminID, _ := middleware.GetUserID(ctx)
	var body struct {
		CouponID uint64   `json:"coupon_id"`
		UserIDs  []uint64 `json:"user_ids"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	g, err := h.logic.GrantCoupon(adminID, body.CouponID, body.UserIDs, 0, true)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return g, nil
}

func (h *CouponHandler) AdminCouponClaims(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	list, total, err := h.logic.CouponClaims(id, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *CouponHandler) AdminCouponRedeems(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	list, total, err := h.logic.CouponRedeems(id, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *CouponHandler) AdminCouponStats(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	st, err := h.logic.CouponStats(id)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return st, nil
}
