package internalapi

import (
	"encoding/json"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *TaskHandler) InternalEvent(w http.ResponseWriter, r *http.Request) {
	var req biz.TaskEventReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.HandleEvent(req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *TaskHandler) InternalDeductPoints(w http.ResponseWriter, r *http.Request) {
	var req biz.PointsLedgerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	p, err := h.logic.DeductPoints(req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, p)
}

func (h *TaskHandler) InternalRefundPoints(w http.ResponseWriter, r *http.Request) {
	var req biz.PointsLedgerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	p, err := h.logic.RefundPoints(req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, p)
}
