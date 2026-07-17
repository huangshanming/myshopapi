package handler

import (
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/response"
	"mymall/services/catalog-service/internal/notify/logic"
	"mymall/services/catalog-service/internal/notify/repository"
	"mymall/services/catalog-service/internal/svc"
)

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
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
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
	data, err := h.logic.List(f)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "ok")
}

func (h *NotificationHandler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	data, err := h.logic.UnreadCount(shopID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, data, "ok")
}

func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err := h.logic.MarkRead(id, shopID); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}

func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	shopID := middleware.GetShopID(r.Context())
	if shopID == 0 {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	if err := h.logic.MarkAllRead(shopID); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "ok")
}
