package admin

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"
	"strconv"

	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"
)

func (h *SeckillHandler) AdminGetSeckillRule(ctx context.Context, in appinput.CallInput) (any, error) {
	rule, err := h.logic.GetSeckillRule()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return rule, nil
}

func (h *SeckillHandler) AdminUpdateSeckillRule(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.SeckillRuleReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	rule, err := h.logic.UpdateSeckillRule(req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return rule, nil
}

func (h *SeckillHandler) AdminListSeckillSessions(ctx context.Context, in appinput.CallInput) (any, error) {
	p, ps := in.Page()
	list, total, err := h.logic.ListSeckillSessions(p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}

func (h *SeckillHandler) AdminListSeckillEntries(ctx context.Context, in appinput.CallInput) (any, error) {
	p, ps := in.Page()
	sid, _ := strconv.ParseUint(in.QueryGet("session_id"), 10, 64)
	list, total, err := h.logic.ListAdminSeckillEntries(sid, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}
