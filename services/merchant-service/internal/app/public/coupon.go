package public

import (
	"context"
	"net/http"
	"strconv"

	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
)

func (h *CouponHandler) PublicCouponCenter(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, _ := middleware.GetUserID(ctx)
	if userID == 0 && in.Request != nil {
		if raw := in.Request.Header.Get("X-User-Id"); raw != "" {
			userID, _ = strconv.ParseUint(raw, 10, 64)
		}
	}
	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	list, err := h.logic.ListCenter(userID, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list}, nil
}

func (h *CouponHandler) PublicCouponPopup(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, _ := middleware.GetUserID(ctx)
	list, err := h.logic.ListPopup(userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list}, nil
}
