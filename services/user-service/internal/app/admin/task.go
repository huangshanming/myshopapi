package admin

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"net/http"
	"strconv"
)

func (h *TaskHandler) AdminList(ctx context.Context, in appinput.CallInput) (any, error) {
	list, err := h.logic.AdminListTasks(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list}, nil
}

func (h *TaskHandler) AdminUpdate(ctx context.Context, in appinput.CallInput) (any, error) {
	id, err := strconv.ParseUint(in.Path("id"), 10, 64)
	if err != nil || id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "任务ID无效")
	}
	var req biz.UpdateTaskReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	t, err := h.logic.AdminUpdateTask(ctx, id, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return t, nil
}
