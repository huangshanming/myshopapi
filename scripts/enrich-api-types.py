#!/usr/bin/env python3
"""Enrich go-zero .api with Req/Resp types and (Req) returns (Resp) on every route.

Usage: enrich-api-types.py <path-to.api>
Idempotent: skips if file already contains 'type PageReq'.
"""
from __future__ import annotations

import re
import sys
from pathlib import Path

COMMON_TYPES = r'''
// ---- shared DTOs (goctl types) ----
type PageReq {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=10"`
}
type IdPathReq {
	Id uint64 `path:"id"`
}
type CodePathReq {
	Code string `path:"code"`
}
type FilePathReq {
	File string `path:"file"`
}
type EmptyResp {}
type PageListResp {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}
type CountResp {
	Count int64 `json:"count"`
}
type URLResp {
	Url string `json:"url"`
}
type TokenResp {
	Token string `json:"token"`
}
type PointsResp {
	Points int64 `json:"points"`
}
type LoginReq {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	ShopId   uint64 `json:"shop_id,optional"`
}
type LoginResp {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
type RegisterReq {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}
type AddressReq {
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	Province     string `json:"province"`
	City         string `json:"city"`
	District     string `json:"district"`
	Detail       string `json:"detail"`
	IsDefault    int    `json:"is_default,optional"`
}
type AddressUpdateReq {
	Id           uint64 `path:"id"`
	ContactName  string `json:"contact_name,optional"`
	ContactPhone string `json:"contact_phone,optional"`
	Province     string `json:"province,optional"`
	City         string `json:"city,optional"`
	District     string `json:"district,optional"`
	Detail       string `json:"detail,optional"`
	IsDefault    int    `json:"is_default,optional"`
}
type WalletOrderOpReq {
	UserId  uint64  `json:"user_id"`
	Amount  float64 `json:"amount"`
	OrderId uint64  `json:"order_id"`
	OrderNo string  `json:"order_no"`
}
type WalletAdjustReq {
	Id     uint64  `path:"id"`
	Amount float64 `json:"amount"`
	Remark string  `json:"remark,optional"`
}
type InternalAddressReq {
	Id     uint64 `form:"id"`
	UserId uint64 `form:"user_id"`
}
type MenuReq {
	ParentId  uint64 `json:"parent_id,optional"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Path      string `json:"path,optional"`
	Component string `json:"component,optional"`
	Icon      string `json:"icon,optional"`
	Perms     string `json:"perms,optional"`
	Sort      int    `json:"sort,optional"`
	Visible   int    `json:"visible,optional"`
	Status    int    `json:"status,optional"`
}
type MenuUpdateReq {
	Id        uint64 `path:"id"`
	ParentId  uint64 `json:"parent_id,optional"`
	Name      string `json:"name,optional"`
	Type      string `json:"type,optional"`
	Path      string `json:"path,optional"`
	Component string `json:"component,optional"`
	Icon      string `json:"icon,optional"`
	Perms     string `json:"perms,optional"`
	Sort      int    `json:"sort,optional"`
	Visible   int    `json:"visible,optional"`
	Status    int    `json:"status,optional"`
}
type RoleReq {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status int    `json:"status,optional"`
	Remark string `json:"remark,optional"`
}
type RoleUpdateReq {
	Id     uint64 `path:"id"`
	Code   string `json:"code,optional"`
	Name   string `json:"name,optional"`
	Status int    `json:"status,optional"`
	Remark string `json:"remark,optional"`
}
type RoleMenusReq {
	Id      uint64   `path:"id"`
	MenuIds []uint64 `json:"menu_ids"`
}
type UserStatusReq {
	Id     uint64 `path:"id"`
	Status int    `json:"status"`
}
type UserUpdateReq {
	Id       uint64 `path:"id"`
	Nickname string `json:"nickname,optional"`
	Avatar   string `json:"avatar,optional"`
	Gender   int    `json:"gender,optional"`
	Mobile   string `json:"mobile,optional"`
}
type UserResetPwdReq {
	Id       uint64 `path:"id"`
	Password string `json:"password"`
}
type AdminCreateReq {
	Mobile   string   `json:"mobile"`
	Password string   `json:"password"`
	Nickname string   `json:"nickname,optional"`
	RoleIds  []uint64 `json:"role_ids,optional"`
}
type AdminRolesReq {
	Id      uint64   `path:"id"`
	RoleIds []uint64 `json:"role_ids"`
}
type AdminResetPwdReq {
	Id       uint64 `path:"id"`
	Password string `json:"password"`
}
type ConfigItemReq {
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	Remark      string `json:"remark,optional"`
}
type ConfigBatchReq {
	Items []ConfigItemReq `json:"items"`
}
type ListUsersReq {
	Mobile   string `form:"mobile,optional"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}
type ListAdminsReq {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=10"`
}
type RegionListReq {
	ParentCode string `form:"parent_code,optional"`
}
type NotifyCreateReq {
	UserId  uint64 `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Type    string `json:"type,optional"`
}
type AdminSendReq {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	UserIds  []uint64 `json:"user_ids,optional"`
	SendAll  bool     `json:"send_all,optional"`
}
type TaskEventReq {
	Event string `json:"event"`
	RefId string `json:"ref_id,optional"`
}
type PointsLedgerReq {
	UserId uint64 `json:"user_id"`
	Points int64  `json:"points"`
	Reason string `json:"reason,optional"`
	RefNo  string `json:"ref_no,optional"`
}
type UpdateTaskReq {
	Id          uint64 `path:"id"`
	Name        string `json:"name,optional"`
	Points      int64  `json:"points,optional"`
	Status      string `json:"status,optional"`
	Description string `json:"description,optional"`
}
type ExchangeReq {
	ProductId uint64 `json:"product_id"`
	AddressId uint64 `json:"address_id,optional"`
	Remark    string `json:"remark,optional"`
}
type PointsProductSaveReq {
	Name         string `json:"name"`
	CoverUrl     string `json:"cover_url,optional"`
	Description  string `json:"description,optional"`
	PointsPrice  int    `json:"points_price,optional"`
	Stock        int    `json:"stock,optional"`
	PerUserLimit int    `json:"per_user_limit,optional"`
	Status       string `json:"status,optional"`
	Sort         int    `json:"sort,optional"`
}
type PointsProductUpdateReq {
	Id           uint64 `path:"id"`
	Name         string `json:"name,optional"`
	CoverUrl     string `json:"cover_url,optional"`
	Description  string `json:"description,optional"`
	PointsPrice  int    `json:"points_price,optional"`
	Stock        int    `json:"stock,optional"`
	PerUserLimit int    `json:"per_user_limit,optional"`
	Status       string `json:"status,optional"`
	Sort         int    `json:"sort,optional"`
}
type PointsProductStatusReq {
	Id     uint64 `path:"id"`
	Status string `json:"status"`
}
type ListPointsProductsReq {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status,optional"`
	Keyword  string `form:"keyword,optional"`
}
type ListPointsOrdersReq {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status,optional"`
	OrderNo  string `form:"order_no,optional"`
	Keyword  string `form:"keyword,optional"`
	UserId   uint64 `form:"user_id,optional"`
}
type ShipReq {
	Id             uint64 `path:"id"`
	ExpressCompany string `json:"express_company"`
	ExpressNo      string `json:"express_no"`
}
type RemarkReq {
	Id     uint64 `path:"id"`
	Remark string `json:"remark"`
}
type AdminWalletLogsReq {
	Id       uint64 `path:"id"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}
type NotificationRecipientsReq {
	Id       uint64 `path:"id"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}
type AnyResp {
	Data interface{} `json:"data,optional"`
}
'''

# handler name -> (req_type or None, resp_type or None)
# None req = no request type; None resp = no returns (httpx.Ok)
USER_HANDLER_TYPES: dict[str, tuple[str | None, str | None]] = {
    "Healthz": (None, "EmptyResp"),
    "Readyz": (None, "EmptyResp"),
    "Metrics": (None, None),  # prometheus raw
    "Login": ("LoginReq", "LoginResp"),
    "Register": ("RegisterReq", "AnyResp"),
    "ListRegions": ("RegionListReq", "PageListResp"),
    "RegionTree": (None, "AnyResp"),
    "InternalFreezeWallet": ("WalletOrderOpReq", None),
    "InternalUnfreezeWallet": ("WalletOrderOpReq", None),
    "InternalSettleWallet": ("WalletOrderOpReq", None),
    "InternalGet": ("InternalAddressReq", "AnyResp"),
    "InternalCreateNotification": ("NotifyCreateReq", "AnyResp"),
    "InternalEvent": ("TaskEventReq", None),
    "InternalDeductPoints": ("PointsLedgerReq", "PointsResp"),
    "InternalRefundPoints": ("PointsLedgerReq", "PointsResp"),
    "UserProfile": (None, "AnyResp"),
    "UserGetWallet": (None, "AnyResp"),
    "UserWalletLogs": ("PageReq", "PageListResp"),
    "UserListAddresses": (None, "PageListResp"),
    "UserCreateAddress": ("AddressReq", "AnyResp"),
    "UserUpdateAddress": ("AddressUpdateReq", None),
    "UserDeleteAddress": ("IdPathReq", None),
    "SetDefault": ("IdPathReq", None),
    "ListNotifications": ("PageReq", "PageListResp"),
    "UnreadNotificationCount": (None, "CountResp"),
    "MarkNotificationRead": ("IdPathReq", None),
    "MarkAllNotificationsRead": (None, None),
    "UserPoints": (None, "PointsResp"),
    "UserPointLogs": ("PageReq", "PageListResp"),
    "UserListTasks": (None, "PageListResp"),
    "UserCheckin": (None, "AnyResp"),
    "UserClaim": ("CodePathReq", "PointsResp"),
    "UserReportEvent": ("TaskEventReq", None),
    "AuthMe": (None, "AnyResp"),
    "MenuTree": (None, "AnyResp"),
    "CreateMenu": ("MenuReq", "AnyResp"),
    "UpdateMenu": ("MenuUpdateReq", None),
    "DeleteMenu": ("IdPathReq", None),
    "ListRoles": (None, "PageListResp"),
    "CreateRole": ("RoleReq", "AnyResp"),
    "UpdateRole": ("RoleUpdateReq", None),
    "DeleteRole": ("IdPathReq", None),
    "GetRoleMenus": ("IdPathReq", "AnyResp"),
    "AssignRoleMenus": ("RoleMenusReq", None),
    "ListUsers": ("ListUsersReq", "PageListResp"),
    "GetUser": ("IdPathReq", "AnyResp"),
    "UpdateUser": ("UserUpdateReq", None),
    "SetUserStatus": ("UserStatusReq", None),
    "ResetUserPassword": ("UserResetPwdReq", None),
    "GenerateUserToken": ("IdPathReq", "AnyResp"),
    "AdminGetWallet": ("IdPathReq", "AnyResp"),
    "AdminAdjustWallet": ("WalletAdjustReq", "AnyResp"),
    "AdminWalletLogs": ("AdminWalletLogsReq", "PageListResp"),
    "AdminListUserAddresses": ("IdPathReq", "PageListResp"),
    "ListAdmins": ("ListAdminsReq", "PageListResp"),
    "CreateAdmin": ("AdminCreateReq", "AnyResp"),
    "GetAdminRoles": ("IdPathReq", "AnyResp"),
    "AssignAdminRoles": ("AdminRolesReq", None),
    "ResetAdminPassword": ("AdminResetPwdReq", None),
    "ListConfigs": (None, "PageListResp"),
    "SaveConfigs": ("ConfigBatchReq", None),
    "AdminSendNotification": ("AdminSendReq", "AnyResp"),
    "AdminListNotificationSends": ("PageReq", "PageListResp"),
    "AdminListNotificationRecipients": ("NotificationRecipientsReq", "AnyResp"),
    "AdminListTasks": (None, "PageListResp"),
    "AdminUpdateTask": ("UpdateTaskReq", "AnyResp"),
    "ServePointsMallUpload": ("FilePathReq", None),
    "Exchange": ("ExchangeReq", "AnyResp"),
    "ListUserPointsOrders": ("PageReq", "PageListResp"),
    "DetailUserPointsOrder": ("IdPathReq", "AnyResp"),
    "ListPointsProducts": ("ListPointsProductsReq", "PageListResp"),
    "CreatePointsProduct": ("PointsProductSaveReq", "AnyResp"),
    "UploadPointsProduct": (None, "URLResp"),
    "DetailPointsProduct": ("IdPathReq", "AnyResp"),
    "UpdatePointsProduct": ("PointsProductUpdateReq", "AnyResp"),
    "SetPointsProductStatus": ("PointsProductStatusReq", "AnyResp"),
    "DeletePointsProduct": ("IdPathReq", None),
    "ListPointsOrders": ("ListPointsOrdersReq", "PageListResp"),
    "DetailPointsOrder": ("IdPathReq", "AnyResp"),
    "ShipPointsOrder": ("ShipReq", "AnyResp"),
    "CompletePointsOrder": ("IdPathReq", "AnyResp"),
    "CancelPointsOrder": ("RemarkReq", "AnyResp"),
    "RemarkPointsOrder": ("RemarkReq", "AnyResp"),
}


def enrich(api_path: Path, mapping: dict[str, tuple[str | None, str | None]]) -> None:
    text = api_path.read_text()
    if "type PageReq {" in text or "type PageReq\n" in text:
        print(f"already enriched: {api_path}")
        return

    # Insert types after syntax line
    m = re.match(r"(syntax\s*=\s*\"v1\"\s*\n)", text)
    if not m:
        raise SystemExit(f"no syntax=v1 in {api_path}")
    text = text[: m.end()] + "\n" + COMMON_TYPES + "\n" + text[m.end() :]

    # Rewrite each route: @handler X \n METHOD path  -> with (Req) returns (Resp)
    def repl_route(mm: re.Match) -> str:
        indent, handler, method, path = mm.group(1), mm.group(2), mm.group(3), mm.group(4)
        req, resp = mapping.get(handler, (None, "AnyResp"))
        line = f"{indent}{method} {path}"
        if req:
            line += f" ({req})"
        if resp:
            line += f" returns ({resp})"
        return f"{indent}@handler {handler}\n{line}"

    text2, n = re.subn(
        r"^([ \t]*)@handler\s+(\w+)\s*\n[ \t]*(get|post|put|delete|patch)\s+(\S+)",
        repl_route,
        text,
        flags=re.M | re.I,
    )
    api_path.write_text(text2)
    print(f"enriched {api_path}: {n} routes")


def main() -> None:
    if len(sys.argv) < 2:
        raise SystemExit("usage: enrich-api-types.py <api-file> [user|generic]")
    path = Path(sys.argv[1])
    kind = sys.argv[2] if len(sys.argv) > 2 else "user"
    if kind == "user":
        enrich(path, USER_HANDLER_TYPES)
    else:
        # generic: PageReq for GET with query-ish, Empty/Any for others
        text = path.read_text()
        handlers = re.findall(r"@handler\s+(\w+)", text)
        mapping = {}
        for h in handlers:
            # heuristic filled later by service-specific maps
            mapping[h] = (None, "AnyResp")
        enrich(path, mapping)


if __name__ == "__main__":
    main()
