package merchant

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"
)

func (h *HomepageSlotHandler) MerchantListSlotPackages(ctx context.Context, in appinput.CallInput) (any, error) {
	list, err := h.logic.ListSlotPackages(in.QueryGet("slot_type"), true)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *HomepageSlotHandler) MerchantBuySlot(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	userID, ok := middleware.GetUserID(ctx)
	if !ok || shopID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	var req biz.BuySlotReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	order, err := h.logic.BuySlot(shopID, userID, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return order, nil
}

func (h *HomepageSlotHandler) MerchantListSlotOrders(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺")
	}
	p, ps := in.Page()
	list, total, err := h.logic.ListSlotOrders(shopID, in.QueryGet("slot_type"), in.QueryGet("status"), p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}
