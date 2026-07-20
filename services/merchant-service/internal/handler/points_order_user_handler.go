package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type PointsOrderUserHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.PointsOrderLogic
}

func NewPointsOrderUserHandler(svcCtx *svc.ServiceContext) *PointsOrderUserHandler {
	return &PointsOrderUserHandler{
		svcCtx: svcCtx,
		logic:  logic.NewPointsOrderLogic(context.Background(), svcCtx),
	}
}

func (h *PointsOrderUserHandler) Exchange(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	var req logic.ExchangeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	o, err := h.logic.UserExchange(userID, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, o)
}

func (h *PointsOrderUserHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	page, pageSize := middleware.ParsePage(r)
	list, total, err := h.logic.UserList(userID, page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}

func (h *PointsOrderUserHandler) Detail(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusUnauthorized, "未登录"))
		return
	}
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	o, err := h.logic.UserGet(userID, id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, o)
}
