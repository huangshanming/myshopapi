package internalapi

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"net/http"
)

func (h *TaskHandler) InternalEvent(ctx context.Context, in appinput.CallInput) (any, error) {
	var req biz.TaskEventReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.HandleEvent(ctx, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *TaskHandler) InternalDeductPoints(ctx context.Context, in appinput.CallInput) (any, error) {
	var req biz.PointsLedgerReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	p, err := h.logic.DeductPoints(ctx, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return p, nil
}

func (h *TaskHandler) InternalRefundPoints(ctx context.Context, in appinput.CallInput) (any, error) {
	var req biz.PointsLedgerReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	p, err := h.logic.RefundPoints(ctx, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return p, nil
}
