package internalapi

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"net/http"
)

func (h *NotificationHandler) InternalCreateNotification(ctx context.Context, in appinput.CallInput) (any, error) {
	var req biz.NotifyCreateReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	n, err := h.logic.CreateNotification(ctx, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return n, nil
}
