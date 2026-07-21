package internalapi

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"
)

func (h *SeckillHandler) SeckillConsume(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.SeckillConsumeReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	data, err := h.logic.ConsumeSeckill(req.EntryID, req.ProductID, req.Quantity)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return data, nil
}

func (h *SeckillHandler) SeckillRestore(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.SeckillRestoreReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.RestoreSeckill(req.EntryID, req.Quantity); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}
