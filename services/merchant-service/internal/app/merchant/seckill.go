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

func (h *SeckillHandler) MerchantSeckillSessions(ctx context.Context, in appinput.CallInput) (any, error) {
	data, err := h.logic.MerchantSeckillSessions()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *SeckillHandler) MerchantApplySeckill(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	userID, _ := middleware.GetUserID(ctx)
	var req types.SeckillApplyReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	entry, err := h.logic.ApplySeckill(shopID, userID, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return entry, nil
}

func (h *SeckillHandler) MerchantSetSeckillAutoRenew(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "报名ID无效")
	}
	var req types.SeckillAutoRenewReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	entry, err := h.logic.SetSeckillAutoRenew(shopID, id, req.AutoRenew)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return entry, nil
}

func (h *SeckillHandler) MerchantListSeckillEntries(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	p, ps := in.Page()
	list, total, err := h.logic.ListShopSeckillEntries(shopID, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}
