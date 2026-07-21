#!/usr/bin/env python3
"""Rewrite user-service logic files to call biz directly (no httpinvoke/app)."""
from __future__ import annotations

from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
LOGIC = ROOT / "services/user-service/internal/logic"

# Common boilerplate pieces
HDR = '''package {pkg}

import (
	"context"
{imports}

	"github.com/zeromicro/go-zero/core/logx"
)

type {typ} struct {{
	logx.Logger
	svcCtx *svc.ServiceContext
}}

func New{typ}(ctx context.Context, svcCtx *svc.ServiceContext) *{typ} {{
	return &{typ}{{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}}
}}
'''

def write(path: Path, pkg: str, typ: str, imports: list[str], body: str):
    imps = "\n".join("\t" + i for i in imports)
    path.write_text(HDR.format(pkg=pkg, typ=typ, imports=imps) + "\n" + body + "\n")
    print("wrote", path.relative_to(ROOT))


def need_user(body_lines: str) -> str:
    return f'''	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {{
		return {body_lines}
	}}
'''


# ---------- generate files ----------

# Login
write(LOGIC/"public/auth/login_logic.go", "auth", "LoginLogic", [
    '"net/http"',
    '"mymall/pkg/xerr"',
    '"mymall/services/user-service/internal/biz"',
    '"mymall/services/user-service/internal/svc"',
    '"mymall/services/user-service/internal/types"',
], '''func (l *LoginLogic) Login(ctx context.Context, req *types.LoginReq) (resp *types.LoginResp, err error) {
	token, user, err := biz.NewUserLogic(l.svcCtx).LoginWithShop(ctx, req.Mobile, req.Password, req.ShopID)
	if err != nil {
		return nil, xerr.New(http.StatusUnauthorized, err.Error())
	}
	return &types.LoginResp{
		Token: token,
		User: map[string]interface{}{
			"id": user.ID, "mobile": user.Mobile, "nickname": user.Nickname,
			"avatar": user.Avatar, "role": user.Role, "status": user.Status,
		},
	}, nil
}''')

write(LOGIC/"public/auth/register_logic.go", "auth", "RegisterLogic", [
    '"net/http"',
    '"mymall/pkg/xerr"',
    '"mymall/services/user-service/internal/biz"',
    '"mymall/services/user-service/internal/svc"',
    '"mymall/services/user-service/internal/types"',
], '''func (l *RegisterLogic) Register(ctx context.Context, req *types.RegisterReq) (resp *types.AnyResp, err error) {
	u, err := biz.NewUserLogic(l.svcCtx).Register(ctx, req.Mobile, req.Password)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: u}, nil
}''')

write(LOGIC/"public/region/list_regions_logic.go", "region", "ListRegionsLogic", [
    '"net/http"',
    '"mymall/pkg/xerr"',
    '"mymall/services/user-service/internal/biz"',
    '"mymall/services/user-service/internal/svc"',
    '"mymall/services/user-service/internal/types"',
], '''func (l *ListRegionsLogic) ListRegions(ctx context.Context, req *types.RegionListReq) (resp *types.PageListResp, err error) {
	list, err := biz.NewRegionLogic(l.svcCtx).ListChildren(ctx, req.ParentCode)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: int64(len(list))}, nil
}''')

write(LOGIC/"public/region/region_tree_logic.go", "region", "RegionTreeLogic", [
    '"net/http"',
    '"mymall/pkg/xerr"',
    '"mymall/services/user-service/internal/biz"',
    '"mymall/services/user-service/internal/svc"',
    '"mymall/services/user-service/internal/types"',
], '''func (l *RegionTreeLogic) RegionTree(ctx context.Context) (resp *types.AnyResp, err error) {
	tree, err := biz.NewRegionLogic(l.svcCtx).Tree(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: tree}, nil
}''')

# Profile
write(LOGIC/"user/profile/user_profile_logic.go", "profile", "UserProfileLogic", [
    '"net/http"',
    '"mymall/pkg/middleware"',
    '"mymall/pkg/xerr"',
    '"mymall/services/user-service/internal/biz"',
    '"mymall/services/user-service/internal/svc"',
    '"mymall/services/user-service/internal/types"',
], '''func (l *UserProfileLogic) UserProfile(ctx context.Context) (resp *types.AnyResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	u, err := biz.NewUserLogic(l.svcCtx).GetProfile(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: u}, nil
}''')

# Notifications
write(LOGIC/"user/notification/list_notifications_logic.go", "notification", "ListNotificationsLogic", [
    '"net/http"',
    '"mymall/pkg/middleware"',
    '"mymall/pkg/xerr"',
    '"mymall/services/user-service/internal/biz"',
    '"mymall/services/user-service/internal/svc"',
    '"mymall/services/user-service/internal/types"',
], '''func (l *ListNotificationsLogic) ListNotifications(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	list, total, err := biz.NewUserLogic(l.svcCtx).ListMyNotifications(ctx, userID, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: total}, nil
}''')

write(LOGIC/"user/notification/unread_notification_count_logic.go", "notification", "UnreadNotificationCountLogic", [
    '"net/http"',
    '"mymall/pkg/middleware"',
    '"mymall/pkg/xerr"',
    '"mymall/services/user-service/internal/biz"',
    '"mymall/services/user-service/internal/svc"',
    '"mymall/services/user-service/internal/types"',
], '''func (l *UnreadNotificationCountLogic) UnreadNotificationCount(ctx context.Context) (resp *types.CountResp, err error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未登录")
	}
	n, err := biz.NewUserLogic(l.svcCtx).UnreadCount(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.CountResp{Count: n}, nil
}''')

write(LOGIC/"user/notification/mark_notification_read_logic.go", "notification", "MarkNotificationReadLogic", [
    '"net/http"',
    '"mymall/pkg/middleware"',
    '"mymall/pkg/xerr"',
    '"mymall/services/user-service/internal/biz"',
    '"mymall/services/user-service/internal/svc"',
    '"mymall/services/user-service/internal/types"',
], '''func (l *MarkNotificationReadLogic) MarkNotificationRead(ctx context.Context, req *types.IdPathReq) error {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return xerr.New(http.StatusUnauthorized, "未登录")
	}
	if err := biz.NewUserLogic(l.svcCtx).MarkRead(ctx, userID, req.Id); err != nil {
		return xerr.New(http.StatusBadRequest, err.Error())
	}
	return nil
}''')

write(LOGIC/"user/notification/mark_all_notifications_read_logic.go", "notification", "MarkAllNotificationsReadLogic", [
    '"net/http"',
    '"mymall/pkg/middleware"',
    '"mymall/pkg/xerr"',
    '"mymall/services/user-service/internal/biz"',
    '"mymall/services/user-service/internal/svc"',
], '''func (l *MarkAllNotificationsReadLogic) MarkAllNotificationsRead(ctx context.Context) error {
	userID, ok := middleware.GetUserID(ctx)
	if !ok || userID == 0 {
		return xerr.New(http.StatusUnauthorized, "未登录")
	}
	if err := biz.NewUserLogic(l.svcCtx).MarkAllRead(ctx, userID); err != nil {
		return xerr.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}''')

print("batch1 done")
