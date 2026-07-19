package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/response"
	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type AddressHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.AddressLogic
}

func NewAddressHandler(svcCtx *svc.ServiceContext) *AddressHandler {
	return &AddressHandler{
		svcCtx: svcCtx,
		logic:  logic.NewAddressLogic(svcCtx),
	}
}

func (h *AddressHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		response.Error(w, "未登录", http.StatusUnauthorized)
		return
	}
	list, err := h.logic.List(userID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, list, "ok")
}

func (h *AddressHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		response.Error(w, "未登录", http.StatusUnauthorized)
		return
	}
	var req types.AddressReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	a, err := h.logic.Create(userID, req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, a, "创建成功")
}

func (h *AddressHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		response.Error(w, "未登录", http.StatusUnauthorized)
		return
	}
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "地址ID无效", http.StatusBadRequest)
		return
	}
	var req types.AddressReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.Update(userID, id, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "更新成功")
}

func (h *AddressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		response.Error(w, "未登录", http.StatusUnauthorized)
		return
	}
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "地址ID无效", http.StatusBadRequest)
		return
	}
	if err := h.logic.Delete(userID, id); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "已删除")
}

func (h *AddressHandler) SetDefault(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID == 0 {
		response.Error(w, "未登录", http.StatusUnauthorized)
		return
	}
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "地址ID无效", http.StatusBadRequest)
		return
	}
	if err := h.logic.SetDefault(userID, id); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "已设为默认")
}

func (h *AddressHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "用户ID无效", http.StatusBadRequest)
		return
	}
	list, err := h.logic.List(userID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, list, "ok")
}

func (h *AddressHandler) InternalGet(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)
	id, _ := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	if userID == 0 || id == 0 {
		response.Error(w, "参数无效", http.StatusBadRequest)
		return
	}
	a, err := h.logic.Get(userID, id)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, a, "ok")
}
