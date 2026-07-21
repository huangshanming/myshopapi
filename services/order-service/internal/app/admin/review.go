package admin

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"net/http"
	"strconv"
)

func (h *ReviewHandler) AdminList(ctx context.Context, in appinput.CallInput) (any, error) {
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	level := in.QueryGet("rating_level")
	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	list, total, err := h.logic.AdminList(ctx, shopID, level, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *ReviewHandler) AdminDelete(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "评价ID无效")
	}
	if err := h.logic.SoftDelete(ctx, id, 0); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}
