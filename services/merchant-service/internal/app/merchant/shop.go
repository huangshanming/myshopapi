package merchant

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"
	"strconv"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"
)

func (h *ShopHandler) Apply(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	var req types.ApplyReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	app, err := h.logic.Apply(ctx, userID, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return app, nil
}

func (h *ShopHandler) MyShops(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	shops, err := h.logic.MyShops(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return shops, nil
}

func (h *ShopHandler) UpdateMyShop(ctx context.Context, in appinput.CallInput) (any, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	shopID, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	var req types.UpdateShopReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UpdateMyShop(ctx, shopID, userID, req); err != nil {
		return nil, xerr.New(http.StatusForbidden, err.Error())
	}
	return nil, nil
}
