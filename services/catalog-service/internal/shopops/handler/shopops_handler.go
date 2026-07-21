package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mymall/pkg/httpserver"
	"mymall/pkg/middleware"
	"mymall/services/catalog-service/internal/shopops/model"
	"mymall/services/catalog-service/internal/shopops/repository"
	"mymall/services/catalog-service/internal/shopops/types"
	"mymall/services/catalog-service/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mymall/pkg/xerr")

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
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(r.Context(), shopID, uid)
	perms, _ := h.svcCtx.ShopRBAC.ListPerms(r.Context(), shopID, uid)
	menus, _ := h.svcCtx.ShopRBAC.ListMenusForUser(r.Context(), shopID, uid)
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
		"perms": perms, "menus": menus, "menu_tree": repository.BuildShopMenuTree(menus),
		"is_owner": h.svcCtx.ShopRBAC.IsOwner(r.Context(), shopID, uid),
	})
}

func (h *ShopOpsHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	list, err := h.svcCtx.ShopRBAC.ListRoles(r.Context(), shopID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *ShopOpsHandler) ListMenus(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	_ = h.svcCtx.ShopRBAC.EnsureShopMenus(r.Context(), )
	_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(r.Context(), shopID, uid)
	menus, err := h.svcCtx.ShopRBAC.MenuTree(r.Context(), )
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, repository.BuildShopMenuTree(menus))
}

func (h *ShopOpsHandler) RoleMenus(w http.ResponseWriter, r *http.Request) {
	_, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	ids, err := h.svcCtx.ShopRBAC.ListRoleMenuIDs(r.Context(), id)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, ids)
}

func (h *ShopOpsHandler) SaveRole(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok || !h.svcCtx.ShopRBAC.IsOwner(r.Context(), shopID, uid) {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "仅店主可操作"))
		return
	}
	var req types.ShopRoleReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	id, _ := strconv.ParseUint(httpserver.PathParam(r, "id"), 10, 64)
	role := &model.ShopRole{ID: id, ShopID: shopID, Code: req.Code, Name: req.Name, Remark: req.Remark, Status: 1}
	if role.Code == "" {
		role.Code = "custom"
	}
	if err := h.svcCtx.ShopRBAC.SaveRole(r.Context(), role, req.MenuIDs); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, role)
}

func (h *ShopOpsHandler) ListStaff(w http.ResponseWriter, r *http.Request) {
	shopID, _, ok := h.shopUser(r)
	if !ok {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "缺少店铺上下文"))
		return
	}
	list, err := h.svcCtx.ShopRBAC.ListStaff(r.Context(), shopID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusInternalServerError, err.Error()))
		return
	}
	httpx.OkJsonCtx(r.Context(), w, list)
}

func (h *ShopOpsHandler) BindStaff(w http.ResponseWriter, r *http.Request) {
	shopID, uid, ok := h.shopUser(r)
	if !ok || !h.svcCtx.ShopRBAC.IsOwner(r.Context(), shopID, uid) {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusForbidden, "仅店主可操作"))
		return
	}
	var req types.ShopStaffReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "参数错误"))
		return
	}
	if req.Mobile == "" {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "请填写手机号"))
		return
	}
	if req.RoleID == 0 {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, "请选择角色"))
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
		userID, err = h.svcCtx.ShopRBAC.CreateStaffUser(r.Context(), req.Mobile, req.Password, req.Nickname)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
	default:
		userID, err = h.svcCtx.ShopRBAC.FindUserIDByMobile(r.Context(), req.Mobile)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
			return
		}
	}
	if err := h.svcCtx.ShopRBAC.BindStaff(r.Context(), shopID, userID, req.RoleID); err != nil {
		httpx.ErrorCtx(r.Context(), w, xerr.New(http.StatusBadRequest, err.Error()))
		return
	}
	msg := "已绑定"
	if mode == "create" {
		msg = "已创建账号并绑定店铺"
	}
	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"user_id": userID, "msg": msg})
}
