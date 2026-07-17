package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/pkg/response"
	"mymall/services/catalog-service/internal/shopops/model"
	"mymall/services/catalog-service/internal/shopops/repository"
	"mymall/services/catalog-service/internal/shopops/types"
	"mymall/services/catalog-service/internal/svc"
)

// ShopOpsHandler 商家店铺 RBAC / 员工
type ShopOpsHandler struct {
	svcCtx *svc.ServiceContext
}

func NewShopOpsHandler(svcCtx *svc.ServiceContext) *ShopOpsHandler {
	return &ShopOpsHandler{svcCtx: svcCtx}
}

func (h *ShopOpsHandler) shopUser(r *http.Request) (shopID, userID uint64, ok bool) {
	shopID = middleware.GetShopID(r.Context())
	userID, _ = middleware.GetUserID(r.Context())
	return shopID, userID, shopID > 0 && userID > 0
}

func (h *ShopOpsHandler) AuthMe(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(shopID, uid)
	perms, _ := h.svcCtx.ShopRBAC.ListPerms(shopID, uid)
	menus, _ := h.svcCtx.ShopRBAC.ListMenusForUser(shopID, uid)
	response.Success(w, map[string]interface{}{
		"perms": perms, "menus": menus, "menu_tree": repository.BuildShopMenuTree(menus),
		"is_owner": h.svcCtx.ShopRBAC.IsOwner(shopID, uid),
	}, "ok")
}

func (h *ShopOpsHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	list, err := h.svcCtx.ShopRBAC.ListRoles(shopID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, list, "ok")
}

func (h *ShopOpsHandler) ListMenus(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	_ = h.svcCtx.ShopRBAC.EnsureShopMenus()
	_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(shopID, uid)
	menus, err := h.svcCtx.ShopRBAC.MenuTree()
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, repository.BuildShopMenuTree(menus), "ok")
}

func (h *ShopOpsHandler) RoleMenus(w http.ResponseWriter, r *http.Request) {
	_, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	ids, err := h.svcCtx.ShopRBAC.ListRoleMenuIDs(id)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, ids, "ok")
}

func (h *ShopOpsHandler) SaveRole(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok || !h.svcCtx.ShopRBAC.IsOwner(shopID, uid) {
		response.Error(w, "仅店主可操作", http.StatusForbidden)
		return
	}
	var req types.ShopRoleReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	role := &model.ShopRole{ID: id, ShopID: shopID, Code: req.Code, Name: req.Name, Remark: req.Remark, Status: 1}
	if role.Code == "" {
		role.Code = "custom"
	}
	if err := h.svcCtx.ShopRBAC.SaveRole(role, req.MenuIDs); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, role, "ok")
}

func (h *ShopOpsHandler) ListStaff(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		response.Error(w, "缺少店铺上下文", http.StatusForbidden)
		return
	}
	list, err := h.svcCtx.ShopRBAC.ListStaff(shopID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(w, list, "ok")
}

func (h *ShopOpsHandler) BindStaff(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok || !h.svcCtx.ShopRBAC.IsOwner(shopID, uid) {
		response.Error(w, "仅店主可操作", http.StatusForbidden)
		return
	}
	var req types.ShopStaffReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if req.Mobile == "" {
		response.Error(w, "请填写手机号", http.StatusBadRequest)
		return
	}
	if req.RoleID == 0 {
		response.Error(w, "请选择角色", http.StatusBadRequest)
		return
	}
	mode := req.Mode
	if mode == "" {
		mode = "bind"
	}
	var userID uint64
	var err error
	switch mode {
	case "create":
		userID, err = h.svcCtx.ShopRBAC.CreateStaffUser(req.Mobile, req.Password, req.Nickname)
		if err != nil {
			response.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		userID, err = h.svcCtx.ShopRBAC.FindUserIDByMobile(req.Mobile)
		if err != nil {
			response.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if err := h.svcCtx.ShopRBAC.BindStaff(shopID, userID, req.RoleID); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg := "已绑定"
	if mode == "create" {
		msg = "已创建账号并绑定店铺"
	}
	response.Success(w, map[string]interface{}{"user_id": userID}, msg)
}
