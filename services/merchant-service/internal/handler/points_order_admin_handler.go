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

type PointsOrderAdminHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.PointsOrderLogic
}

func NewPointsOrderAdminHandler(svcCtx *svc.ServiceContext) *PointsOrderAdminHandler {
	return &PointsOrderAdminHandler{
		svcCtx: svcCtx,
		logic:  logic.NewPointsOrderLogic(context.Background(), svcCtx),
	}
}

func (h *PointsOrderAdminHandler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize := middleware.ParsePage(r)
	var userID uint64
	if v := r.URL.Query().Get("user_id"); v != "" {
		userID, _ = strconv.ParseUint(v, 10, 64)
	}
	list, total, err := h.logic.AdminList(page, pageSize,
		r.URL.Query().Get("status"),
		r.URL.Query().Get("order_no"),
		r.URL.Query().Get("keyword"),
		userID,
	)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"list": list, "total": total})
}

func (h *PointsOrderAdminHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	o, err := h.logic.AdminGet(id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, o)
}

func (h *PointsOrderAdminHandler) Ship(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	var req logic.ShipReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	o, err := h.logic.AdminShip(id, req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, o)
}

func (h *PointsOrderAdminHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	o, err := h.logic.AdminComplete(id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, o)
}

func (h *PointsOrderAdminHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	var req logic.RemarkReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	o, err := h.logic.AdminCancel(id, req.AdminRemark)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, o)
}

func (h *PointsOrderAdminHandler) Remark(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "订单ID无效"))
		return
	}
	var req logic.RemarkReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	o, err := h.logic.AdminRemark(id, req.AdminRemark)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, o)
}
