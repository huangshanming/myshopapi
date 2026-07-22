#!/usr/bin/env python3
"""Continue go-zero resp polish: order empty/URL; user EmptyResp + typed lists."""
from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
USR = ROOT / "services/user-service"
ORD = ROOT / "services/order-service"


def patch_order() -> None:
    # types
    biz = ORD / "internal/types/biz_types.go"
    t = biz.read_text()
    if "type URLResp struct" not in t:
        t += """
type URLResp struct {
	Url string `json:"url"`
}

type OrderDetailData struct {
	Order      interface{} `json:"order"`
	AfterSales interface{} `json:"after_sales"`
}

type ListData struct {
	List interface{} `json:"list"`
}
"""
        biz.write_text(t)
        print("+order typed")

    types_go = ORD / "internal/types/types.go"
    tt = types_go.read_text()
    if "type URLResp struct" not in tt:
        tt = tt.replace(
            "type PageListResp struct {",
            "type URLResp struct {\n\tUrl string `json:\"url\"`\n}\n\ntype PageListResp struct {",
        )
        types_go.write_text(tt)

    api = next((ORD / "api").glob("*.api"))
    at = api.read_text()
    if "type URLResp {" not in at:
        at = at.replace("type EmptyResp {}", "type EmptyResp {}\n\ntype URLResp {\n\tUrl string `json:\"url\"`\n}")
    at = at.replace(
        "post /api/v1/user/review-uploads returns (AnyResp)",
        "post /api/v1/user/review-uploads returns (URLResp)",
    )
    at = at.replace(
        "delete /api/v1/admin/reviews/:id (IdPathReq) returns (AnyResp)",
        "delete /api/v1/admin/reviews/:id (IdPathReq) returns (EmptyResp)",
    )
    api.write_text(at)

    # empty AnyResp → EmptyResp
    for rel in [
        "internal/logic/merchant/review/merchant_delete_logic.go",
        "internal/logic/admin/review/admin_delete_logic.go",
        "internal/logic/admin/logistics/admin_delete_logistics_logic.go",
    ]:
        p = ORD / rel
        text = p.read_text()
        text = text.replace("*types.AnyResp, error)", "*types.EmptyResp, error)")
        text = text.replace("return &types.AnyResp{}, nil", "return &types.EmptyResp{}, nil")
        p.write_text(text)
        print("order EmptyResp", rel)

    # upload
    p = ORD / "internal/logic/user/review/user_upload_review_logic.go"
    text = p.read_text()
    text = text.replace(
        "func (l *UserUploadReviewLogic) UserUploadReview(ctx context.Context, r *http.Request) (*types.AnyResp, error) {",
        "func (l *UserUploadReviewLogic) UserUploadReview(ctx context.Context, r *http.Request) (*types.URLResp, error) {",
    )
    text = text.replace(
        'return &types.AnyResp{Data: map[string]string{"url": url}}, nil',
        "return &types.URLResp{Url: url}, nil",
    )
    p.write_text(text)
    print("order URLResp upload")

    # order details
    for rel in [
        "internal/logic/admin/order/admin_detail_logic.go",
        "internal/logic/merchant/order/merchant_detail_logic.go",
    ]:
        p = ORD / rel
        text = p.read_text()
        text = text.replace(
            'return &types.AnyResp{Data: map[string]interface{}{"order": order, "after_sales": as}}, nil',
            "return &types.AnyResp{Data: types.OrderDetailData{Order: order, AfterSales: as}}, nil",
        )
        p.write_text(text)
        print("order detail", rel)

    # logistics options
    p = ORD / "internal/logic/shared/logistics/logistics_options_logic.go"
    if p.exists():
        text = p.read_text()
        text2 = text.replace(
            "return &types.AnyResp{Data: list}, nil",
            "return &types.AnyResp{Data: types.ListData{List: list}}, nil",
        )
        if text2 != text:
            p.write_text(text2)
            print("order logistics options")


def convert_user_error_only() -> None:
    """error-only logic → (*EmptyResp, error); handlers OkJson EmptyResp; api returns EmptyResp."""
    methods = []
    for p in (USR / "internal/logic").rglob("*_logic.go"):
        # skip file serve / metrics that write to ResponseWriter
        if "serve_points_mall_upload" in p.name or "metrics_logic" in p.name:
            continue
        t = p.read_text()
        # match: func (l *X) Method(...) error {
        for m in re.finditer(
            r"func \(l \*(\w+)\) (\w+)\(([^)]*)\) error \{",
            t,
        ):
            method = m.group(2)
            methods.append((p, method, m.group(0), m.group(3)))

    converted = []
    for p, method, full, params in methods:
        t = p.read_text()
        m = re.search(
            rf"func \(l \*(\w+)\) {re.escape(method)}\(([^)]*)\) error \{{",
            t,
        )
        if not m:
            continue
        recv = m.group(1)
        params = m.group(2)
        new_sig = f"func (l *{recv}) {method}({params}) (*types.EmptyResp, error) {{"
        t2 = t.replace(
            f"func (l *{recv}) {method}({params}) error {{",
            new_sig,
            1,
        )
        # change bare `return nil` success at end — careful with return nil, err patterns
        # Replace `return nil\n}` that are success - typically `return nil` without error after successful op
        # Safer: replace final success returns that are exactly `return nil` (not return nil, ...)
        # After conversion, `return nil` for success should become `return &types.EmptyResp{}, nil`
        # and `return err` / `return xerr...` stay as returning error as second value... wait
        # Old: `return xerr.New(...)` means error return
        # Old: `return nil` means success
        # New: `return nil, xerr.New(...)` and `return &types.EmptyResp{}, nil`

        # Convert return statements inside the method body only
        brace = t2.find(new_sig) + len(new_sig) - 1  # points to {
        # find method end
        start = t2.find(new_sig)
        body_start = t2.find("{", start)
        depth = 0
        end = body_start
        for i, c in enumerate(t2[body_start:], body_start):
            if c == "{":
                depth += 1
            elif c == "}":
                depth -= 1
                if depth == 0:
                    end = i
                    break
        body = t2[body_start + 1 : end]
        new_body = body
        # return xerr... or return err → return nil, ...
        new_body = re.sub(
            r"\treturn (xerr\.New\([^;]+\)|err|fmt\.Errorf\([^;]+\))\n",
            r"\treturn nil, \1\n",
            new_body,
        )
        new_body = re.sub(r"\treturn nil\n", "\treturn &types.EmptyResp{}, nil\n", new_body)
        t2 = t2[: body_start + 1] + new_body + t2[end:]
        if t2 != t:
            p.write_text(t2)
            converted.append(method)
            print("user EmptyResp logic", method, p.relative_to(ROOT))

    # handlers: err := l.Method → resp, err := ; Ok → OkJson EmptyResp
    for p in (USR / "internal/handler").rglob("*_handler.go"):
        t = p.read_text()
        o = t
        for method in converted:
            pat = (
                r"(\t\tl := [^\n]+\n)"
                r"\t\terr := l\." + method + r"\(([^)]*)\)\n"
                r"\t\tif err != nil \{\n"
                r"\t\t\thttpx\.ErrorCtx\(r\.Context\(\), w, err\)\n"
                r"\t\t\} else \{\n"
                r"\t\t\thttpx\.Ok\(w\)\n"
                r"\t\t\}"
            )
            repl = (
                r"\1"
                r"\t\tresp, err := l." + method + r"(\2)\n"
                r"\t\tif err != nil {\n"
                r"\t\t\thttpx.ErrorCtx(r.Context(), w, err)\n"
                r"\t\t} else {\n"
                r"\t\t\thttpx.OkJsonCtx(r.Context(), w, resp)\n"
                r"\t\t}"
            )
            t2, n = re.subn(pat, repl, t)
            if n:
                t = t2
        if t != o:
            p.write_text(t)
            print("user handler", p.relative_to(ROOT))

    # api: add returns (EmptyResp) where missing for these handlers
    api = next((USR / "api").glob("*.api"))
    at = api.read_text()
    for method in converted:
        # @handler Method\n\tVERB path (Req?)   without returns
        pat = rf"(@handler {re.escape(method)}\n\t(?:get|post|put|delete|patch) [^\n]+?)(?: returns \([^\)]+\))?\n"
        def repl_api(m, method=method):
            line = m.group(1)
            if "returns (" in line:
                return m.group(0).replace("returns (AnyResp)", "returns (EmptyResp)")
            return line + " returns (EmptyResp)\n"
        at2, n = re.subn(pat, repl_api, at)
        if n:
            at = at2
        else:
            # try if already has no newline after
            pass
    api.write_text(at)
    print("user api EmptyResp patched")


def user_typed_helpers() -> None:
    types_go = USR / "internal/types/types.go"
    t = types_go.read_text()
    extras = """
type IdListResp struct {
	Ids []uint64 `json:"ids"`
}

type TreeResp struct {
	Tree interface{} `json:"tree"`
}

type NotificationRecipientsResp struct {
	Batch interface{} `json:"batch"`
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type UserTokenResp struct {
	Token    string `json:"token"`
	UserId   uint64 `json:"user_id"`
	Mobile   string `json:"mobile"`
	Nickname string `json:"nickname"`
	Role     string `json:"role,optional"`
	ShopId   uint64 `json:"shop_id,optional"`
}
"""
    if "type IdListResp struct" not in t:
        # insert before last or after TokenResp
        if "type TokenResp struct" in t:
            t = t.replace(
                "type TokenResp struct {\n\tToken string `json:\"token\"`\n}",
                "type TokenResp struct {\n\tToken string `json:\"token\"`\n}\n" + extras,
            )
        else:
            t += extras
        types_go.write_text(t)
        print("+user typed structs")

    api = next((USR / "api").glob("*.api"))
    at = api.read_text()
    if "type IdListResp {" not in at:
        at = at.replace(
            "type TokenResp {\n\tToken string `json:\"token\"`\n}",
            "type TokenResp {\n\tToken string `json:\"token\"`\n}\n\n"
            "type IdListResp {\n\tIds []uint64 `json:\"ids\"`\n}\n\n"
            "type TreeResp {\n\tTree interface{} `json:\"tree\"`\n}\n\n"
            "type NotificationRecipientsResp {\n\tBatch interface{} `json:\"batch\"`\n\tList  interface{} `json:\"list\"`\n\tTotal int64 `json:\"total\"`\n}\n",
        )
    # route updates
    repls = [
        ("get /api/v1/admin/roles/:id/menus (IdPathReq) returns (AnyResp)", "get /api/v1/admin/roles/:id/menus (IdPathReq) returns (IdListResp)"),
        ("get /api/v1/admin/admins/:id/roles (IdPathReq) returns (AnyResp)", "get /api/v1/admin/admins/:id/roles (IdPathReq) returns (IdListResp)"),
        ("get /api/v1/admin/menus returns (AnyResp)", "get /api/v1/admin/menus returns (TreeResp)"),
        ("get /api/v1/regions/tree returns (AnyResp)", "get /api/v1/regions/tree returns (TreeResp)"),
        (
            "get /api/v1/admin/notifications/sends/:id/recipients (NotificationRecipientsReq) returns (AnyResp)",
            "get /api/v1/admin/notifications/sends/:id/recipients (NotificationRecipientsReq) returns (NotificationRecipientsResp)",
        ),
        ("post /api/v1/user/tasks/checkin returns (AnyResp)", "post /api/v1/user/tasks/checkin returns (AnyResp)"),  # keep entity
    ]
    for a, b in repls:
        at = at.replace(a, b)
    api.write_text(at)

    # logic rewrites
    patches = [
        (
            "internal/logic/admin/role/get_role_menus_logic.go",
            "*types.AnyResp",
            "*types.IdListResp",
            "return &types.AnyResp{Data: ids}, nil",
            "return &types.IdListResp{Ids: ids}, nil",
        ),
        (
            "internal/logic/admin/staff/get_admin_roles_logic.go",
            "*types.AnyResp",
            "*types.IdListResp",
            "return &types.AnyResp{Data: ids}, nil",
            "return &types.IdListResp{Ids: ids}, nil",
        ),
        (
            "internal/logic/admin/menu/menu_tree_logic.go",
            "*types.AnyResp",
            "*types.TreeResp",
            "return &types.AnyResp{Data: tree}, nil",
            "return &types.TreeResp{Tree: tree}, nil",
        ),
        (
            "internal/logic/public/region/region_tree_logic.go",
            "*types.AnyResp",
            "*types.TreeResp",
            "return &types.AnyResp{Data: tree}, nil",
            "return &types.TreeResp{Tree: tree}, nil",
        ),
    ]
    for rel, old_t, new_t, old_r, new_r in patches:
        p = USR / rel
        text = p.read_text().replace(old_t, new_t).replace(old_r, new_r)
        # also fix signature form (resp *types.X
        p.write_text(text)
        print("typed", rel)

    p = USR / "internal/logic/admin/notification/admin_list_notification_recipients_logic.go"
    text = p.read_text()
    text = text.replace("*types.AnyResp", "*types.NotificationRecipientsResp")
    text = re.sub(
        r"return &types\.AnyResp\{Data: map\[string\]interface\{\}\{\n\t\t\"batch\": batch,\n\t\t\"list\":  list,\n\t\t\"total\": total,\n\t\}\}, nil",
        "return &types.NotificationRecipientsResp{Batch: batch, List: list, Total: total}, nil",
        text,
    )
    p.write_text(text)
    print("typed recipients")


def main() -> None:
    patch_order()
    convert_user_error_only()
    user_typed_helpers()


if __name__ == "__main__":
    main()
