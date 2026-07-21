package admin

import (
	"encoding/json"
	"mymall/pkg/httpserver"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func (h *TaskHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.AdminListTasks(r.Context())
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list})
}

func (h *TaskHandler) AdminUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "任务ID无效"))
		return
	}
	var req biz.UpdateTaskReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	t, err := h.logic.AdminUpdateTask(r.Context(), id, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, t)
}
