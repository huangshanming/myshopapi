package handler

import (
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/services/catalog-service/internal/notify/logic"
	"mymall/services/catalog-service/internal/notify/repository"
	"mymall/services/catalog-service/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/pkg/xerr")

type NotificationHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.NotificationLogic
}

func NewNotificationHandler(svcCtx *svc.ServiceContext) *NotificationHandler {
	return &NotificationHandler{svcCtx: svcCtx, logic: logic.NewNotificationLogic(svcCtx)}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	page, pageSize := middleware.ParsePage(r)
	f := repository.NotificationListFilter{ShopID: shopID, Page: page, PageSize: pageSize}
	if s := r.URL.Query().Get("is_read"); s == "0" || s == "1" {
		v := int8(0)
		if s == "1" {
			v = 1
		}
		f.IsRead = &v
	}
	data, err := h.logic.List(r.Context(), f)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *NotificationHandler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	data, err := h.logic.UnreadCount(r.Context(), shopID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, data)
}

func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.MarkRead(r.Context(), id, shopID); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	if err := h.logic.MarkAllRead(r.Context(), shopID); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
