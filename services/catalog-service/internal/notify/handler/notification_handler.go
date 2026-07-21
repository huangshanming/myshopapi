package handler

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"
	"strconv"

	"mymall/pkg/middleware"
	"mymall/services/catalog-service/internal/notify/logic"
	"mymall/services/catalog-service/internal/notify/repository"
	"mymall/services/catalog-service/internal/svc"

	"mymall/pkg/xerr"
)

type NotificationHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.NotificationLogic
}

func NewNotificationHandler(svcCtx *svc.ServiceContext) *NotificationHandler {
	return &NotificationHandler{svcCtx: svcCtx, logic: logic.NewNotificationLogic(svcCtx)}
}

func (h *NotificationHandler) List(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	page, pageSize := in.Page()
	f := repository.NotificationListFilter{ShopID: shopID, Page: page, PageSize: pageSize}
	if s := in.QueryGet("is_read"); s == "0" || s == "1" {
		v := int8(0)
		if s == "1" {
			v = 1
		}
		f.IsRead = &v
	}
	data, err := h.logic.List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *NotificationHandler) UnreadCount(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	data, err := h.logic.UnreadCount(ctx, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *NotificationHandler) MarkRead(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	if err := h.logic.MarkRead(ctx, id, shopID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *NotificationHandler) MarkAllRead(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	if err := h.logic.MarkAllRead(ctx, shopID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}
