package admin

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/types"
	"net/http"
	"strconv"
)

func (h *AdminHandler) AdminSendNotification(ctx context.Context, in appinput.CallInput) (any, error) {
	adminID, _ := middleware.GetUserID(ctx)
	var req biz.AdminSendReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	userLogic := biz.NewUserLogic(h.svcCtx)
	batch, err := userLogic.AdminSend(ctx, adminID, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return batch, nil
}

func (h *AdminHandler) AdminListNotificationSends(ctx context.Context, in appinput.CallInput) (any, error) {
	page, pageSize := in.Page()
	userLogic := biz.NewUserLogic(h.svcCtx)
	list, total, err := userLogic.ListSendBatches(ctx, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list, "total": total}, nil
}

func (h *AdminHandler) AdminListNotificationRecipients(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	if id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	page, pageSize := in.Page()
	userLogic := biz.NewUserLogic(h.svcCtx)
	batch, err := userLogic.GetSendBatch(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "发送记录不存在")
	}
	list, total, err := userLogic.ListBatchRecipients(ctx, id, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{
		"batch": batch,
		"list":  list,
		"total": total,
	}, nil
}

func (h *AdminHandler) IsSuperAdmin(userID uint64) bool {
	return h.logic.IsSuperAdmin(context.Background(), userID)
}

func (h *AdminHandler) HasPerm(userID uint64, code string) bool {
	return h.logic.HasPerm(context.Background(), userID, code)
}

func (h *AdminHandler) AuthMe(ctx context.Context, in appinput.CallInput) (any, error) {
	uid, _ := middleware.GetUserID(ctx)
	data, err := h.logic.AuthMe(ctx, uid)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return data, nil
}

func (h *AdminHandler) MenuTree(ctx context.Context, in appinput.CallInput) (any, error) {
	tree, err := h.logic.MenuTreeAll(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return tree, nil
}

func (h *AdminHandler) CreateMenu(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.MenuReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	m, err := h.logic.CreateMenu(ctx, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return m, nil
}

func (h *AdminHandler) UpdateMenu(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.MenuReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UpdateMenu(ctx, id, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) DeleteMenu(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	if err := h.logic.DeleteMenu(ctx, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) ListRoles(ctx context.Context, in appinput.CallInput) (any, error) {
	list, err := h.logic.ListRoles(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *AdminHandler) CreateRole(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.RoleReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	role, err := h.logic.CreateRole(ctx, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return role, nil
}

func (h *AdminHandler) UpdateRole(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.RoleReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UpdateRole(ctx, id, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) DeleteRole(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	if err := h.logic.DeleteRole(ctx, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) GetRoleMenus(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	ids, err := h.logic.GetRoleMenus(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return ids, nil
}

func (h *AdminHandler) AssignRoleMenus(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.RoleMenusReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.AssignRoleMenus(ctx, id, req.MenuIDs); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) ListUsers(ctx context.Context, in appinput.CallInput) (any, error) {
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	mobile := in.QueryGet("mobile")
	list, total, err := h.logic.ListUsers(ctx, page, pageSize, mobile)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}

func (h *AdminHandler) SetUserStatus(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.UserStatusReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.SetUserStatus(ctx, id, req.Status); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) GetUser(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	user, err := h.logic.GetUser(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return user, nil
}

func (h *AdminHandler) UpdateUser(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.UserUpdateReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.UpdateUser(ctx, id, req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) ResetUserPassword(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.UserResetPwdReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.ResetUserPassword(ctx, id, req.Password); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) GenerateUserToken(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	data, err := h.logic.GenerateUserToken(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return data, nil
}

func (h *AdminHandler) ListAdmins(ctx context.Context, in appinput.CallInput) (any, error) {
	page, _ := strconv.Atoi(in.QueryGet("page"))
	pageSize, _ := strconv.Atoi(in.QueryGet("page_size"))
	list, total, err := h.logic.ListAdmins(ctx, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return types.PageListResp{Total: total, List: list}, nil
}

func (h *AdminHandler) CreateAdmin(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.AdminCreateReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	user, err := h.logic.CreateAdmin(ctx, req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return user, nil
}

func (h *AdminHandler) GetAdminRoles(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	ids, err := h.logic.AdminRoleIDs(ctx, id)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return ids, nil
}

func (h *AdminHandler) AssignAdminRoles(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.AdminRolesReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.AssignAdminRoles(ctx, id, req.RoleIDs); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) ResetAdminPassword(ctx context.Context, in appinput.CallInput) (any, error) {
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var req types.AdminResetPwdReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.ResetAdminPassword(ctx, id, req.Password); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}

func (h *AdminHandler) ListConfigs(ctx context.Context, in appinput.CallInput) (any, error) {
	list, err := h.logic.ListConfigs(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return list, nil
}

func (h *AdminHandler) SaveConfigs(ctx context.Context, in appinput.CallInput) (any, error) {
	var req types.ConfigBatchReq
	if err := appinput.BindBody(in, &req); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := h.logic.SaveConfigs(ctx, req.Items); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil, nil
}
