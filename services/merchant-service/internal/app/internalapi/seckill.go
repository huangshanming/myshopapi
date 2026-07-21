package internalapi

import (
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/types"
)

func (h *SeckillHandler) SeckillConsume(w http.ResponseWriter, r *http.Request) {
	var req types.SeckillConsumeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	data, err := h.logic.ConsumeSeckill(req.EntryID, req.ProductID, req.Quantity)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *SeckillHandler) SeckillRestore(w http.ResponseWriter, r *http.Request) {
	var req types.SeckillRestoreReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.RestoreSeckill(req.EntryID, req.Quantity); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
