package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/response"
	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/repository"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"
)

type LogisticsHandler struct {
	logic *logic.LogisticsLogic
}

func NewLogisticsHandler(svcCtx *svc.ServiceContext) *LogisticsHandler {
	return &LogisticsHandler{logic: logic.NewLogisticsLogic(svcCtx)}
}

func (h *LogisticsHandler) AdminList(w http.ResponseWriter, r *http.Request) {
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
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "查询成功")
}

func (h *LogisticsHandler) Options(w http.ResponseWriter, r *http.Request) {
	list, err := h.logic.Options(r.URL.Query().Get("keyword"))
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, list, "查询成功")
}

func (h *LogisticsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req types.LogisticsSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	c, err := h.logic.Create(req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, c, "创建成功")
}

func (h *LogisticsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "ID无效", http.StatusBadRequest)
		return
	}
	var req types.LogisticsSaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.Update(id, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "已更新")
}

func (h *LogisticsHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "ID无效", http.StatusBadRequest)
		return
	}
	var req types.LogisticsStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.UpdateStatus(id, req.Status); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "已更新")
}

func (h *LogisticsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "ID无效", http.StatusBadRequest)
		return
	}
	if err := h.logic.Delete(id); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "已删除")
}
