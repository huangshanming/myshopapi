package handler

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"
	"strconv"

	"mymall/pkg/middleware"
	"mymall/services/catalog-service/internal/shopops/model"
	"mymall/services/catalog-service/internal/shopops/repository"
	"mymall/services/catalog-service/internal/shopops/types"
	"mymall/services/catalog-service/internal/svc"

	"mymall/pkg/xerr"
)

// ShopOpsHandler 商家店铺 RBAC / 员工
type ShopOpsHandler struct {
	svcCtx *svc.ServiceContext
}

func NewShopOpsHandler(svcCtx *svc.ServiceContext) *ShopOpsHandler {
	return &ShopOpsHandler{svcCtx: svcCtx}
}

func (h *ShopOpsHandler) shopUser(ctx context.Context) (shopID, userID uint64, ok bool) {
	shopID = middleware.GetShopID(ctx)
	userID, _ = middleware.GetUserID(ctx)
	return shopID, userID, shopID > 0 && userID > 0
}

func (h *ShopOpsHandler) AuthMe(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, uid, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(ctx, shopID, uid)
	perms, _ := h.svcCtx.ShopRBAC.ListPerms(ctx, shopID, uid)
	menus, _ := h.svcCtx.ShopRBAC.ListMenusForUser(ctx, shopID, uid)
	return map[string]interface{}{
		"perms": perms, "menus": menus, "menu_tree": repository.BuildShopMenuTree(menus),
		"is_owner": h.svcCtx.ShopRBAC.IsOwner(ctx, shopID, uid),
	}, nil
}

func (h *ShopOpsHandler) ListRoles(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	list, err := h.svcCtx.ShopRBAC.ListRoles(ctx, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *ShopOpsHandler) ListMenus(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, uid, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	_ = h.svcCtx.ShopRBAC.EnsureShopMenus(ctx)
	_ = h.svcCtx.ShopRBAC.EnsureOwnerRole(ctx, shopID, uid)
	menus, err := h.svcCtx.ShopRBAC.MenuTree(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return repository.BuildShopMenuTree(menus), nil
}

func (h *ShopOpsHandler) RoleMenus(ctx context.Context, in appinput.CallInput) (any, error) {
	_, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	ids, err := h.svcCtx.ShopRBAC.ListRoleMenuIDs(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return ids, nil
}

func (h *ShopOpsHandler) SaveRole(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, uid, ok := h.shopUser(ctx)
	if !ok || !h.svcCtx.ShopRBAC.IsOwner(ctx, shopID, uid) {
		return nil, xerr.New(http.StatusForbidden, "仅店主可操作")
	}
	var req types.ShopRoleReq
	_ = appinput.BindBody(in, &req)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	role := &model.ShopRole{ID: id, ShopID: shopID, Code: req.Code, Name: req.Name, Remark: req.Remark, Status: 1}
	if role.Code == "" {
		role.Code = "custom"
	}
	if err := h.svcCtx.ShopRBAC.SaveRole(ctx, role, req.MenuIDs); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return role, nil
}

func (h *ShopOpsHandler) ListStaff(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, _, ok := h.shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	list, err := h.svcCtx.ShopRBAC.ListStaff(ctx, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *ShopOpsHandler) BindStaff(ctx context.Context, in appinput.CallInput) (any, error) {
	shopID, uid, ok := h.shopUser(ctx)
	if !ok || !h.svcCtx.ShopRBAC.IsOwner(ctx, shopID, uid) {
		return nil, xerr.New(http.StatusForbidden, "仅店主可操作")
	}
	var req types.ShopStaffReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if req.Mobile == "" {
		return nil, xerr.New(http.StatusBadRequest, "请填写手机号")
	}
	if req.RoleID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "请选择角色")
	}
	mode := req.Mode
	if mode == "" {
		mode = "bind"
	}
	var userID uint64
	var err error
	switch mode {
	case "create":
		userID, err = h.svcCtx.ShopRBAC.CreateStaffUser(ctx, req.Mobile, req.Password, req.Nickname)
		if err != nil {
			return nil, xerr.New(http.StatusBadRequest, err.Error())
		}
	default:
		userID, err = h.svcCtx.ShopRBAC.FindUserIDByMobile(ctx, req.Mobile)
		if err != nil {
			return nil, xerr.New(http.StatusBadRequest, err.Error())
		}
	}
	if err := h.svcCtx.ShopRBAC.BindStaff(ctx, shopID, userID, req.RoleID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	msg := "已绑定"
	if mode == "create" {
		msg = "已创建账号并绑定店铺"
	}
	return map[string]interface{}{"user_id": userID, "msg": msg}, nil
}
