package merchant

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"net/http"
	"strconv"
)

func (h *ReviewHandler) MerchantList(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "店铺未绑定")
	}
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	level := in.QueryGet("rating_level")
	list, total, err := h.logic.MerchantList(ctx, shopID, level, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *ReviewHandler) MerchantReply(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "店铺未绑定")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "评价ID无效")
	}
	var req struct {
		Reply string `json:"reply"`
	}
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.Reply(ctx, shopID, id, req.Reply); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *ReviewHandler) MerchantDelete(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "店铺未绑定")
	}
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "评价ID无效")
	}
	if err := h.logic.SoftDelete(ctx, id, shopID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}
