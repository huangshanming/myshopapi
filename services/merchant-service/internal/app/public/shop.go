package public

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"
	"strconv"

	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"
)

func (h *ShopHandler) PublicListShops(ctx context.Context, in appinput.CallInput) (any, error) {
	p, ps := in.Page()
	list, total, err := h.logic.ListPublicShops(ctx, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}

func (h *ShopHandler) PublicGetShop(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "店铺ID无效")
	}
	shop, err := h.logic.GetPublicShop(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return shop, nil
}
