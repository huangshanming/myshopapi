package admin

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/repository"
	"mymall/services/order-service/internal/types"
	"net/http"
	"strconv"
)

func (h *LogisticsHandler) AdminList(ctx context.Context, in appinput.CallInput) (any, error) {
	p, ps := in.Page()
	f := repository.LogisticsListFilter{
		Name: in.QueryGet("name"), Code: in.QueryGet("code"),
		Page: p, PageSize: ps,
	}
	if s := in.QueryGet("status"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err == nil {
			st := int8(v)
			f.Status = &st
		}
	}
	list, total, err := h.logic.List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}

func (h *LogisticsHandler) Options(ctx context.Context, in appinput.CallInput) (any, error) {
	list, err := h.logic.Options(ctx, in.QueryGet("keyword"))
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *LogisticsHandler) Create(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.LogisticsSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	c, err := h.logic.Create(ctx, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return c, nil
}

func (h *LogisticsHandler) Update(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	var req types.LogisticsSaveReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.Update(ctx, id, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *LogisticsHandler) UpdateStatus(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	var req types.LogisticsStatusReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UpdateStatus(ctx, id, req.Status); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *LogisticsHandler) Delete(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	if err := h.logic.Delete(ctx, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}
