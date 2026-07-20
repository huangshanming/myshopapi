package internalapi

import (
	"encoding/json"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/logic"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *NotificationHandler) InternalCreateNotification(w http.ResponseWriter, r *http.Request) {
	var req logic.NotifyCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	n, err := h.logic.CreateNotification(req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, n)
}
