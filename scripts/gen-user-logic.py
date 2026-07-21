#!/usr/bin/env python3
"""Generate user-service logic files that invoke legacy app handlers via httpinvoke."""
from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
LOGIC = ROOT / "services/user-service/internal/logic"
TYPES_IMPORT = '"mymall/services/user-service/internal/types"'
SVC_IMPORT = '"mymall/services/user-service/internal/svc"'
INVOKE = '"mymall/pkg/httpinvoke"'

# handler -> (app_import_alias, app_constructor_or_func, method, http_method, path_template, how_to_build_query/body)
# path_template uses {Id} etc from req fields

SPECS: dict[str, dict] = {}


def spec(handler, alias, ctor, method, http, path, *, body="req", query=None, path_vars=None, wrap="raw", perm=None):
    SPECS[handler] = {
        "alias": alias,
        "ctor": ctor,
        "method": method,
        "http": http,
        "path": path,
        "body": body,
        "query": query,
        "path_vars": path_vars,
        "wrap": wrap,
        "perm": perm,
    }


# public
spec("Login", "huser", "huser.NewUserHandler(l.svcCtx).Login", "Login", "POST", "/api/v1/user/login", wrap="login")
spec("Register", "huser", "huser.NewUserHandler(l.svcCtx).Register", "Register", "POST", "/api/v1/user/register")
spec("ListRegions", "hpublic", "hpublic.NewRegionHandler(l.svcCtx).List", "List", "GET", "/api/v1/regions", body=None, query="region")
spec("RegionTree", "hpublic", "hpublic.NewRegionHandler(l.svcCtx).Tree", "Tree", "GET", "/api/v1/regions/tree", body=None)

# internal wallet
for h, m in [("InternalFreezeWallet", "Freeze"), ("InternalUnfreezeWallet", "Unfreeze"), ("InternalSettleWallet", "Settle")]:
    spec(h, "hinternal", f"hinternal.NewWalletHandler(l.svcCtx).{m}", m, "POST", f"/api/v1/user/wallet/{m.lower()}", wrap="empty")

spec("InternalGet", "hinternal", "hinternal.NewAddressHandler(l.svcCtx).InternalGet", "InternalGet", "GET", "/api/v1/user/addresses/internal", body=None, query="internal_addr")
spec("InternalCreateNotification", "hinternal", "hinternal.NewNotificationHandler(l.svcCtx).InternalCreateNotification", "InternalCreateNotification", "POST", "/api/v1/internal/notifications")
spec("InternalEvent", "hinternal", "hinternal.NewTaskHandler(l.svcCtx).InternalEvent", "InternalEvent", "POST", "/api/v1/internal/tasks/events", wrap="empty")
spec("InternalDeductPoints", "hinternal", "hinternal.NewTaskHandler(l.svcCtx).InternalDeductPoints", "InternalDeductPoints", "POST", "/api/v1/internal/points/deduct", wrap="points")
spec("InternalRefundPoints", "hinternal", "hinternal.NewTaskHandler(l.svcCtx).InternalRefundPoints", "InternalRefundPoints", "POST", "/api/v1/internal/points/refund", wrap="points")

# user
spec("UserProfile", "huser", "huser.NewUserHandler(l.svcCtx).Profile", "Profile", "GET", "/api/v1/user/profile", body=None)
spec("UserGetWallet", "huser", "huser.NewWalletHandler(l.svcCtx).UserGetWallet", "UserGetWallet", "GET", "/api/v1/user/wallet", body=None)
spec("UserWalletLogs", "huser", "huser.NewWalletHandler(l.svcCtx).UserWalletLogs", "UserWalletLogs", "GET", "/api/v1/user/wallet/logs", body=None, query="page", wrap="page")
spec("UserListAddresses", "huser", "huser.NewAddressHandler(l.svcCtx).List", "List", "GET", "/api/v1/user/addresses", body=None, wrap="list_as_page")
spec("UserCreateAddress", "huser", "huser.NewAddressHandler(l.svcCtx).Create", "Create", "POST", "/api/v1/user/addresses")
spec("UserUpdateAddress", "huser", "huser.NewAddressHandler(l.svcCtx).Update", "Update", "PUT", "/api/v1/user/addresses/{Id}", path_vars=["Id"], wrap="empty")
spec("UserDeleteAddress", "huser", "huser.NewAddressHandler(l.svcCtx).Delete", "Delete", "DELETE", "/api/v1/user/addresses/{Id}", body=None, path_vars=["Id"], wrap="empty")
spec("SetDefault", "huser", "huser.NewAddressHandler(l.svcCtx).SetDefault", "SetDefault", "PUT", "/api/v1/user/addresses/{Id}/default", body=None, path_vars=["Id"], wrap="empty")
spec("ListNotifications", "huser", "huser.NewUserHandler(l.svcCtx).ListNotifications", "ListNotifications", "GET", "/api/v1/user/notifications", body=None, query="page", wrap="page")
spec("UnreadNotificationCount", "huser", "huser.NewUserHandler(l.svcCtx).UnreadNotificationCount", "UnreadNotificationCount", "GET", "/api/v1/user/notifications/unread-count", body=None, wrap="count")
spec("MarkNotificationRead", "huser", "huser.NewUserHandler(l.svcCtx).MarkNotificationRead", "MarkNotificationRead", "POST", "/api/v1/user/notifications/{Id}/read", body=None, path_vars=["Id"], wrap="empty")
spec("MarkAllNotificationsRead", "huser", "huser.NewUserHandler(l.svcCtx).MarkAllNotificationsRead", "MarkAllNotificationsRead", "POST", "/api/v1/user/notifications/read-all", body=None, wrap="empty")
spec("UserPoints", "huser", "huser.NewTaskHandler(l.svcCtx).UserPoints", "UserPoints", "GET", "/api/v1/user/points", body=None, wrap="points")
spec("UserPointLogs", "huser", "huser.NewTaskHandler(l.svcCtx).UserPointLogs", "UserPointLogs", "GET", "/api/v1/user/points/logs", body=None, query="page", wrap="page")
spec("UserListTasks", "huser", "huser.NewTaskHandler(l.svcCtx).UserListTasks", "UserListTasks", "GET", "/api/v1/user/tasks", body=None, wrap="list_as_page")
spec("UserCheckin", "huser", "huser.NewTaskHandler(l.svcCtx).UserCheckin", "UserCheckin", "POST", "/api/v1/user/tasks/checkin", body=None)
spec("UserClaim", "huser", "huser.NewTaskHandler(l.svcCtx).UserClaim", "UserClaim", "POST", "/api/v1/user/tasks/{Code}/claim", body=None, path_vars=["Code"], wrap="points")
spec("UserReportEvent", "huser", "huser.NewTaskHandler(l.svcCtx).UserReportEvent", "UserReportEvent", "POST", "/api/v1/user/tasks/events", wrap="empty")

# points mall user
spec("Exchange", "huser", "huser.NewPointsOrderHandler(l.svcCtx).Exchange", "Exchange", "POST", "/api/v1/user/points-mall/exchange")
spec("ListUserPointsOrders", "huser", "huser.NewPointsOrderHandler(l.svcCtx).List", "List", "GET", "/api/v1/user/points-mall/orders", body=None, query="page", wrap="page")
spec("DetailUserPointsOrder", "huser", "huser.NewPointsOrderHandler(l.svcCtx).Detail", "Detail", "GET", "/api/v1/user/points-mall/orders/{Id}", body=None, path_vars=["Id"])

# admin auth
spec("AuthMe", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).AuthMe", "AuthMe", "GET", "/api/v1/admin/auth/me", body=None)

# admin menus — with perms applied inside app via RequirePermission in old logic; skip perm wrapping for now (middleware RequirePlatformAdmin covers admin)
spec("MenuTree", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).MenuTree", "MenuTree", "GET", "/api/v1/admin/menus", body=None)
spec("CreateMenu", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).CreateMenu", "CreateMenu", "POST", "/api/v1/admin/menus")
spec("UpdateMenu", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).UpdateMenu", "UpdateMenu", "PUT", "/api/v1/admin/menus/{Id}", path_vars=["Id"], wrap="empty")
spec("DeleteMenu", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).DeleteMenu", "DeleteMenu", "DELETE", "/api/v1/admin/menus/{Id}", body=None, path_vars=["Id"], wrap="empty")
spec("ListRoles", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).ListRoles", "ListRoles", "GET", "/api/v1/admin/roles", body=None, wrap="list_as_page")
spec("CreateRole", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).CreateRole", "CreateRole", "POST", "/api/v1/admin/roles")
spec("UpdateRole", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).UpdateRole", "UpdateRole", "PUT", "/api/v1/admin/roles/{Id}", path_vars=["Id"], wrap="empty")
spec("DeleteRole", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).DeleteRole", "DeleteRole", "DELETE", "/api/v1/admin/roles/{Id}", body=None, path_vars=["Id"], wrap="empty")
spec("GetRoleMenus", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).GetRoleMenus", "GetRoleMenus", "GET", "/api/v1/admin/roles/{Id}/menus", body=None, path_vars=["Id"])
spec("AssignRoleMenus", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).AssignRoleMenus", "AssignRoleMenus", "PUT", "/api/v1/admin/roles/{Id}/menus", path_vars=["Id"], wrap="empty")
spec("ListUsers", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).ListUsers", "ListUsers", "GET", "/api/v1/admin/users", body=None, query="list_users", wrap="page")
spec("GetUser", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).GetUser", "GetUser", "GET", "/api/v1/admin/users/{Id}", body=None, path_vars=["Id"])
spec("UpdateUser", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).UpdateUser", "UpdateUser", "PUT", "/api/v1/admin/users/{Id}", path_vars=["Id"], wrap="empty")
spec("SetUserStatus", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).SetUserStatus", "SetUserStatus", "PUT", "/api/v1/admin/users/{Id}/status", path_vars=["Id"], wrap="empty")
spec("ResetUserPassword", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).ResetUserPassword", "ResetUserPassword", "PUT", "/api/v1/admin/users/{Id}/password", path_vars=["Id"], wrap="empty")
spec("GenerateUserToken", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).GenerateUserToken", "GenerateUserToken", "POST", "/api/v1/admin/users/{Id}/token", body=None, path_vars=["Id"])
spec("AdminGetWallet", "hadmin", "hadmin.NewWalletHandler(l.svcCtx).AdminGetWallet", "AdminGetWallet", "GET", "/api/v1/admin/users/{Id}/wallet", body=None, path_vars=["Id"])
spec("AdminAdjustWallet", "hadmin", "hadmin.NewWalletHandler(l.svcCtx).AdminAdjustWallet", "AdminAdjustWallet", "POST", "/api/v1/admin/users/{Id}/wallet/adjust", path_vars=["Id"])
spec("AdminWalletLogs", "hadmin", "hadmin.NewWalletHandler(l.svcCtx).AdminWalletLogs", "AdminWalletLogs", "GET", "/api/v1/admin/users/{Id}/wallet/logs", body=None, path_vars=["Id"], query="page", wrap="page")
spec("AdminListUserAddresses", "hadmin", "hadmin.NewAddressHandler(l.svcCtx).AdminList", "AdminList", "GET", "/api/v1/admin/users/{Id}/addresses", body=None, path_vars=["Id"], wrap="list_as_page")
spec("ListAdmins", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).ListAdmins", "ListAdmins", "GET", "/api/v1/admin/admins", body=None, query="page", wrap="page")
spec("CreateAdmin", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).CreateAdmin", "CreateAdmin", "POST", "/api/v1/admin/admins")
spec("GetAdminRoles", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).GetAdminRoles", "GetAdminRoles", "GET", "/api/v1/admin/admins/{Id}/roles", body=None, path_vars=["Id"])
spec("AssignAdminRoles", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).AssignAdminRoles", "AssignAdminRoles", "PUT", "/api/v1/admin/admins/{Id}/roles", path_vars=["Id"], wrap="empty")
spec("ResetAdminPassword", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).ResetAdminPassword", "ResetAdminPassword", "PUT", "/api/v1/admin/admins/{Id}/password", path_vars=["Id"], wrap="empty")
spec("ListConfigs", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).ListConfigs", "ListConfigs", "GET", "/api/v1/admin/configs", body=None, wrap="list_as_page")
spec("SaveConfigs", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).SaveConfigs", "SaveConfigs", "PUT", "/api/v1/admin/configs", wrap="empty")
spec("AdminSendNotification", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).AdminSendNotification", "AdminSendNotification", "POST", "/api/v1/admin/notifications/send")
spec("AdminListNotificationSends", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).AdminListNotificationSends", "AdminListNotificationSends", "GET", "/api/v1/admin/notifications/sends", body=None, query="page", wrap="page")
spec("AdminListNotificationRecipients", "hadmin", "hadmin.NewAdminHandler(l.svcCtx).AdminListNotificationRecipients", "AdminListNotificationRecipients", "GET", "/api/v1/admin/notifications/sends/{Id}/recipients", body=None, path_vars=["Id"], query="page")
spec("AdminListTasks", "hadmin", "hadmin.NewTaskHandler(l.svcCtx).AdminList", "AdminList", "GET", "/api/v1/admin/tasks", body=None, wrap="list_as_page")
spec("AdminUpdateTask", "hadmin", "hadmin.NewTaskHandler(l.svcCtx).AdminUpdate", "AdminUpdate", "PUT", "/api/v1/admin/tasks/{Id}", path_vars=["Id"])

# admin points mall
spec("ListPointsProducts", "hadmin", "hadmin.NewPointsProductHandler(l.svcCtx).List", "List", "GET", "/api/v1/admin/points-products", body=None, query="points_products", wrap="page")
spec("CreatePointsProduct", "hadmin", "hadmin.NewPointsProductHandler(l.svcCtx).Create", "Create", "POST", "/api/v1/admin/points-products")
spec("DetailPointsProduct", "hadmin", "hadmin.NewPointsProductHandler(l.svcCtx).Detail", "Detail", "GET", "/api/v1/admin/points-products/{Id}", body=None, path_vars=["Id"])
spec("UpdatePointsProduct", "hadmin", "hadmin.NewPointsProductHandler(l.svcCtx).Update", "Update", "PUT", "/api/v1/admin/points-products/{Id}", path_vars=["Id"])
spec("SetPointsProductStatus", "hadmin", "hadmin.NewPointsProductHandler(l.svcCtx).SetStatus", "SetStatus", "PUT", "/api/v1/admin/points-products/{Id}/status", path_vars=["Id"])
spec("DeletePointsProduct", "hadmin", "hadmin.NewPointsProductHandler(l.svcCtx).Delete", "Delete", "DELETE", "/api/v1/admin/points-products/{Id}", body=None, path_vars=["Id"], wrap="empty")
spec("ListPointsOrders", "hadmin", "hadmin.NewPointsOrderHandler(l.svcCtx).List", "List", "GET", "/api/v1/admin/points-orders", body=None, query="points_orders", wrap="page")
spec("DetailPointsOrder", "hadmin", "hadmin.NewPointsOrderHandler(l.svcCtx).Detail", "Detail", "GET", "/api/v1/admin/points-orders/{Id}", body=None, path_vars=["Id"])
spec("ShipPointsOrder", "hadmin", "hadmin.NewPointsOrderHandler(l.svcCtx).Ship", "Ship", "POST", "/api/v1/admin/points-orders/{Id}/ship", path_vars=["Id"])
spec("CompletePointsOrder", "hadmin", "hadmin.NewPointsOrderHandler(l.svcCtx).Complete", "Complete", "POST", "/api/v1/admin/points-orders/{Id}/complete", body=None, path_vars=["Id"])
spec("CancelPointsOrder", "hadmin", "hadmin.NewPointsOrderHandler(l.svcCtx).Cancel", "Cancel", "POST", "/api/v1/admin/points-orders/{Id}/cancel", path_vars=["Id"])
spec("RemarkPointsOrder", "hadmin", "hadmin.NewPointsOrderHandler(l.svcCtx).Remark", "Remark", "PUT", "/api/v1/admin/points-orders/{Id}/remark", path_vars=["Id"])

ALIAS_IMPORTS = {
    "huser": '"mymall/services/user-service/internal/app/user"',
    "hadmin": '"mymall/services/user-service/internal/app/admin"',
    "hpublic": '"mymall/services/user-service/internal/app/public"',
    "hinternal": '"mymall/services/user-service/internal/app/internalapi"',
}


def find_logic_files():
    return list(LOGIC.rglob("*_logic.go"))


def parse_logic_meta(text: str):
    pkg = re.search(r"^package (\w+)", text, re.M).group(1)
    typ = re.search(r"type (\w+) struct", text).group(1)
    fn = re.search(r"func \(l \*%s\) (\w+)\(" % typ, text).group(1)
    has_req = "req *types." in text
    req_type = None
    if has_req:
        req_type = re.search(r"req \*types\.(\w+)", text).group(1)
    has_resp = "resp *types." in text or "(resp *types." in text
    resp_type = None
    if has_resp:
        resp_type = re.search(r"resp \*types\.(\w+)", text).group(1)
    return pkg, typ, fn, req_type, resp_type


def query_code(kind, req_type):
    if kind == "page":
        return 'url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}'
    if kind == "region":
        return 'url.Values{"parent_code": {req.ParentCode}}'
    if kind == "internal_addr":
        return 'url.Values{"id": {fmt.Sprintf("%d", req.Id)}, "user_id": {fmt.Sprintf("%d", req.UserID)}}'
    if kind == "list_users":
        return 'url.Values{"mobile": {req.Mobile}, "page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}'
    if kind == "points_products":
        return 'url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}, "status": {req.Status}, "keyword": {req.Keyword}}'
    if kind == "points_orders":
        return 'url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}, "status": {req.Status}, "order_no": {req.OrderNo}, "keyword": {req.Keyword}, "user_id": {fmt.Sprintf("%d", req.UserID)}}'
    return "nil"


def gen_body(fn, req_type, resp_type, sp):
    imports_extra = set()
    alias = sp["alias"]
    imports_extra.add(ALIAS_IMPORTS[alias])
    imports_extra.add(INVOKE)
    imports_extra.add('"net/url"')
    imports_extra.add('"fmt"')
    imports_extra.add('"encoding/json"')

    path = sp["path"]
    path_vars = sp.get("path_vars") or []
    pv = "nil"
    if path_vars:
        parts = []
        for v in path_vars:
            parts.append(f'"{v.lower()}": fmt.Sprintf("%v", req.{v})')
            # path keys in PathParam are "id", "code", "file"
            key = v.lower()
            if v == "Id":
                key = "id"
            elif v == "Code":
                key = "code"
            elif v == "File":
                key = "file"
            parts[-1] = f'"{key}": fmt.Sprintf("%v", req.{v})'
        pv = "map[string]string{" + ", ".join(parts) + "}"

    q = "nil"
    if sp.get("query"):
        q = query_code(sp["query"], req_type)
        if req_type is None and sp["query"] == "page":
            q = 'url.Values{}'

    body = "nil"
    if sp.get("body") == "req" and req_type:
        body = "req"

    wrap = sp.get("wrap", "raw")
    call = f'''raw, err := httpinvoke.Run(ctx, "{sp["http"]}", "{path}", {pv}, {q}, {body}, {sp["ctor"]})
	if err != nil {{
		return {"nil, " if resp_type else ""}err
	}}'''

    if wrap == "empty" or resp_type is None:
        call = call.replace("raw, err :=", "_, err :=")
        ret = call + "\n\treturn nil"
        if resp_type:
            ret = call + "\n\treturn nil, nil"
    elif wrap == "page":
        ret = call + f'''
	var out types.{resp_type}
	if err := httpinvoke.Decode(raw, &out); err != nil {{
		// raw may already be {{list,total}}
		var m map[string]json.RawMessage
		if err2 := json.Unmarshal(raw, &m); err2 == nil {{
			_ = json.Unmarshal(m["list"], &out.List)
			_ = json.Unmarshal(m["total"], &out.Total)
			return &out, nil
		}}
		return nil, err
	}}
	return &out, nil'''
    elif wrap == "list_as_page":
        ret = call + f'''
	var list interface{{}}
	if err := httpinvoke.Decode(raw, &list); err != nil {{
		return nil, err
	}}
	return &types.{resp_type}{{List: list}}, nil'''
    elif wrap == "count":
        ret = call + f'''
	var out types.{resp_type}
	if err := httpinvoke.Decode(raw, &out); err != nil {{
		return nil, err
	}}
	return &out, nil'''
    elif wrap == "points":
        ret = call + f'''
	var out types.{resp_type}
	if err := httpinvoke.Decode(raw, &out); err != nil {{
		// may be bare number or {{points:n}}
		var n int64
		if err2 := json.Unmarshal(raw, &n); err2 == nil {{
			out.Points = n
			return &out, nil
		}}
		return nil, err
	}}
	return &out, nil'''
    elif wrap == "login":
        ret = call + f'''
	var out types.{resp_type}
	if err := httpinvoke.Decode(raw, &out); err != nil {{
		return nil, err
	}}
	return &out, nil'''
    elif resp_type == "AnyResp":
        ret = call + '''
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil'''
    elif resp_type == "EmptyResp":
        ret = call + "\n\treturn &types.EmptyResp{}, nil"
    else:
        ret = call + f'''
	var out types.{resp_type}
	if err := httpinvoke.Decode(raw, &out); err != nil {{
		return nil, err
	}}
	return &out, nil'''

    return ret, imports_extra


def write_logic(path: Path):
    text = path.read_text()
    pkg, typ, fn, req_type, resp_type = parse_logic_meta(text)

    # special cases
    if fn in ("Healthz", "Readyz", "Metrics", "ServePointsMallUpload", "UploadPointsProduct"):
        return write_special(path, pkg, typ, fn, req_type, resp_type)

    sp = SPECS.get(fn)
    if not sp:
        print("NO SPEC", fn, path)
        return

    body, extra = gen_body(fn, req_type, resp_type, sp)
    imports = [
        '"context"',
        SVC_IMPORT,
    ]
    if req_type or resp_type:
        imports.append(TYPES_IMPORT)
    imports.append('"github.com/zeromicro/go-zero/core/logx"')
    for e in sorted(extra):
        if e not in imports:
            imports.append(e)

    # alias imports need alias
    import_block = []
    for im in imports:
        if "internal/app/user" in im:
            import_block.append(f"\thuser {im}")
        elif "internal/app/admin" in im:
            import_block.append(f"\thadmin {im}")
        elif "internal/app/public" in im:
            import_block.append(f"\thpublic {im}")
        elif "internal/app/internalapi" in im:
            import_block.append(f"\thinternal {im}")
        else:
            import_block.append(f"\t{im}")

    req_param = f", req *types.{req_type}" if req_type else ""
    if resp_type:
        sig = f"func (l *{typ}) {fn}(ctx context.Context{req_param}) (resp *types.{resp_type}, err error)"
    else:
        sig = f"func (l *{typ}) {fn}(ctx context.Context{req_param}) error"

    out = f'''package {pkg}

import (
{chr(10).join(import_block)}
)

type {typ} struct {{
	logx.Logger
	svcCtx *svc.ServiceContext
}}

func New{typ}(svcCtx *svc.ServiceContext) *{typ} {{
	return &{typ}{{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}}
}}

{sig} {{
{body}
}}
'''
    path.write_text(out)
    print("wrote", path.relative_to(ROOT))


def write_special(path, pkg, typ, fn, req_type, resp_type):
    if fn == "Healthz":
        body = '\treturn &types.EmptyResp{}, nil'
        # actual write happens in patched handler
        imports = f'''
import (
	"context"
	{TYPES_IMPORT}
	{SVC_IMPORT}
	"github.com/zeromicro/go-zero/core/logx"
)
'''
    elif fn == "Readyz":
        body = '''\tif err := l.svcCtx.Health.Ready(ctx); err != nil {
		return nil, err
	}
	return &types.EmptyResp{}, nil'''
        imports = f'''
import (
	"context"
	{TYPES_IMPORT}
	{SVC_IMPORT}
	"github.com/zeromicro/go-zero/core/logx"
)
'''
    elif fn == "Metrics":
        body = "\treturn nil"
        imports = f'''
import (
	"context"
	{SVC_IMPORT}
	"github.com/zeromicro/go-zero/core/logx"
)
'''
    elif fn == "ServePointsMallUpload":
        body = '''\t_ = req
	_ = ctx
	return nil'''
        imports = f'''
import (
	"context"
	{TYPES_IMPORT}
	{SVC_IMPORT}
	"github.com/zeromicro/go-zero/core/logx"
)
'''
    elif fn == "UploadPointsProduct":
        # multipart — call app via raw is hard; return error asking handler — implement properly later
        body = '''\treturn nil, xerr.New(400, "请使用 multipart 上传")'''
        imports = f'''
import (
	"context"
	{TYPES_IMPORT}
	{SVC_IMPORT}
	"mymall/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)
'''
    else:
        return

    req_param = f", req *types.{req_type}" if req_type else ""
    if resp_type:
        sig = f"func (l *{typ}) {fn}(ctx context.Context{req_param}) (resp *types.{resp_type}, err error)"
    else:
        sig = f"func (l *{typ}) {fn}(ctx context.Context{req_param}) error"

    out = f'''package {pkg}
{imports}
type {typ} struct {{
	logx.Logger
	svcCtx *svc.ServiceContext
}}

func New{typ}(svcCtx *svc.ServiceContext) *{typ} {{
	return &{typ}{{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}}
}}

{sig} {{
{body}
}}
'''
    path.write_text(out)
    print("special", path.relative_to(ROOT))


def main():
    for path in find_logic_files():
        write_logic(path)


if __name__ == "__main__":
    main()
