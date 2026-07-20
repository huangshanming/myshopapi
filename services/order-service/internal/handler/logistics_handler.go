package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/repository"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/pkg/xerr"
)

type LogisticsAdminHandler struct {
	logic *logic.LogisticsLogic
}

func NewLogisticsAdminHandler(svcCtx *svc.ServiceContext) *LogisticsAdminHandler {
	return &LogisticsAdminHandler{logic: logic.NewLogisticsLogic(context.Background(), svcCtx)}
}

func (h *LogisticsAdminHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	f := repository.LogisticsListFilter{
		Name: r.URL.Query().Get("name"), Code: r.URL.Query().Get("code"),
		Page: p, PageSize: ps,
	}
	if s := r.URL.Query().Get("status"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err == nil {
			st := int8(v)
			f.Status = &st
		}
	}
	list, total, err := h.logic.List(f)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, types.PageListResp{Total: total, List: list})
}

func (h *LogisticsAdminHandler) Options(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.Options(r.URL.Query().Get("keyword"))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *LogisticsAdminHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req types.LogisticsSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	c, err := h.logic.Create(req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, c)
}

func (h *LogisticsAdminHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "ID无效"))
		return
	}
	var req types.LogisticsSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.Update(id, req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *LogisticsAdminHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "ID无效"))
		return
	}
	var req types.LogisticsStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if err := h.logic.UpdateStatus(id, req.Status); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}

func (h *LogisticsAdminHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "ID无效"))
		return
	}
	if err := h.logic.Delete(id); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, nil)
}
