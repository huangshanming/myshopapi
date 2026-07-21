package user

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"
	"strconv"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
)

func (h *CouponHandler) ClaimCoupon(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, _ := middleware.GetUserID(ctx)
	if userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "请先登录")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var body struct {
		Source string `json:"source"`
	}
	_ = appinput.BindBody(in, &body)
	uc, err := h.logic.ClaimCoupon(userID, id, body.Source)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return uc, nil
}

func (h *CouponHandler) ListMyCoupons(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, _ := middleware.GetUserID(ctx)
	if userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "请先登录")
	}
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	list, total, err := h.logic.ListMyCoupons(userID, in.QueryGet("status"), page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}
