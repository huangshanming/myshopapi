package public

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"
	"strconv"

	"mymall/pkg/xerr"
)

func (h *SeckillHandler) PublicSeckillCurrent(ctx context.Context, in appinput.CallInput) (any, error) {
	data, err := h.logic.PublicSeckillCurrent()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *SeckillHandler) PublicSeckillList(ctx context.Context, in appinput.CallInput) (any, error) {
	p, ps := in.Page()
	data, err := h.logic.PublicSeckillList(p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *SeckillHandler) PublicSeckillEntry(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "ID无效")
	}
	data, err := h.logic.PublicSeckillEntry(id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return data, nil
}
