package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/response"
	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
)

type MerchantHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.MerchantLogic
}

func NewMerchantHandler(svcCtx *svc.ServiceContext) *MerchantHandler {
	return &MerchantHandler{
		svcCtx: svcCtx,
		logic:  logic.NewMerchantLogic(svcCtx),
	}
}

func (h *MerchantHandler) Apply(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, "未授权", http.StatusUnauthorized)
		return
	}
	var in types.ApplyReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	app, err := h.logic.Apply(userID, in)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, app, "提交成功")
}

func (h *MerchantHandler) MyShops(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, "未授权", http.StatusUnauthorized)
		return
	}
	shops, err := h.logic.MyShops(userID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, shops, "查询成功")
}

func (h *MerchantHandler) UpdateMyShop(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, "未授权", http.StatusUnauthorized)
		return
	}
	shopID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "店铺ID无效", http.StatusBadRequest)
		return
	}
	var req types.UpdateShopReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.UpdateMyShop(shopID, userID, req); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	response.Success(w, nil, "更新成功")
}

func (h *MerchantHandler) AdminListApplications(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListApplications(r.URL.Query().Get("status"), p, ps)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "查询成功")
}

func (h *MerchantHandler) AdminApprove(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, "未授权", http.StatusUnauthorized)
		return
	}
	appID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "申请ID无效", http.StatusBadRequest)
		return
	}
	shop, err := h.logic.Approve(appID, adminID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, shop, "审核通过")
}

func (h *MerchantHandler) AdminReject(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, "未授权", http.StatusUnauthorized)
		return
	}
	appID, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "申请ID无效", http.StatusBadRequest)
		return
	}
	var req types.RejectReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.Reject(appID, adminID, req.Reason); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "已拒绝")
}

func (h *MerchantHandler) AdminListShops(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListShops(r.URL.Query().Get("status"), r.URL.Query().Get("name"), p, ps)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "查询成功")
}

// PublicListShops C 端公开商户列表（无需登录）
func (h *MerchantHandler) PublicListShops(w http.ResponseWriter, r *http.Request) {
	p, ps := middleware.ParsePage(r)
	list, total, err := h.logic.ListPublicShops(p, ps)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, types.PageListResp{Total: total, List: list}, "查询成功")
}

func (h *MerchantHandler) PublicGetShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(w, "店铺ID无效", http.StatusBadRequest)
		return
	}
	shop, err := h.logic.GetPublicShop(id)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, shop, "ok")
}

func (h *MerchantHandler) AdminCreateShop(w http.ResponseWriter, r *http.Request) {
	var req types.AdminCreateShopReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	shop, err := h.logic.CreateShop(req)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, shop, "创建成功")
}

func (h *MerchantHandler) AdminUpdateShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "店铺ID无效", http.StatusBadRequest)
		return
	}
	var req types.AdminUpdateShopReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.AdminUpdateShop(id, req); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "更新成功")
}

func (h *MerchantHandler) AdminResetOwnerPassword(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "店铺ID无效", http.StatusBadRequest)
		return
	}
	var req types.OwnerPasswordReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if err := h.logic.ResetOwnerPassword(id, req.Password); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "重置成功")
}

func (h *MerchantHandler) AdminGetShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "店铺ID无效", http.StatusBadRequest)
		return
	}
	shop, err := h.logic.GetShop(id)
	if err != nil {
		response.Error(w, "店铺不存在", http.StatusNotFound)
		return
	}
	response.Success(w, shop, "查询成功")
}

func (h *MerchantHandler) AdminDisableShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "店铺ID无效", http.StatusBadRequest)
		return
	}
	var req types.RejectReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.logic.DisableShop(id, req.Reason); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "已禁用")
}

func (h *MerchantHandler) AdminEnableShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, "店铺ID无效", http.StatusBadRequest)
		return
	}
	if err := h.logic.EnableShop(id); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, nil, "已启用")
}
