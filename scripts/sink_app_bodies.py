#!/usr/bin/env python3
"""Sink merchant/catalog logic: inline app method bodies as domain logic calls.

Catalog keys methods as import|Recv|Method to avoid List/Create collisions.
Inlines shopUser/requirePerm helpers when referenced.
"""
from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

SHOP_USER_HELPER = """
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}
""".rstrip(
    "\n"
)

REQUIRE_PERM_HELPER = """
	requirePerm := func(ctx context.Context, code string) error {
		shopID, uid, ok := shopUser(ctx)
		if !ok {
			return xerr.New(http.StatusForbidden, "缺少店铺上下文")
		}
		if middleware.GetUserRole(ctx) == jwt.RoleMerchantOwner {
			_ = l.svcCtx.ShopRBAC.EnsureOwnerRole(ctx, shopID, uid)
		}
		if !l.svcCtx.ShopRBAC.HasPerm(ctx, shopID, uid, code) {
			return xerr.New(http.StatusForbidden, "无权限: "+code)
		}
		return nil
	}
""".rstrip(
    "\n"
)


def find_matching(text: str, open_idx: int, open_ch="(", close_ch=")") -> int:
    depth = 0
    i = open_idx
    in_str = None
    while i < len(text):
        c = text[i]
        if in_str:
            if c == "\\" and in_str != "`":
                i += 2
                continue
            if c == in_str:
                in_str = None
            i += 1
            continue
        if c in "\"'`":
            in_str = c
            i += 1
            continue
        if c == open_ch:
            depth += 1
        elif c == close_ch:
            depth -= 1
            if depth == 0:
                return i
        i += 1
    raise RuntimeError("unbalanced")


def extract_methods(go_text: str) -> dict[str, tuple[str, str]]:
    out = {}
    for m in re.finditer(
        r"func \(h \*(\w+)\) (\w+)\(ctx context\.Context, in appinput\.CallInput\) \(any, error\) ",
        go_text,
    ):
        recv, name = m.group(1), m.group(2)
        brace = go_text.find("{", m.end() - 1)
        end = find_matching(go_text, brace, "{", "}")
        out[name] = (recv, go_text[brace : end + 1])
    return out


def parse_imports(text: str) -> dict[str, str]:
    """alias -> import path"""
    out = {}
    m = re.search(r"import \((.*?)\)", text, re.S)
    if not m:
        return out
    for line in m.group(1).splitlines():
        line = line.strip()
        if not line or line.startswith("//"):
            continue
        mm = re.match(r'(\w+)\s+"([^"]+)"', line)
        if mm:
            out[mm.group(1)] = mm.group(2)
            continue
        mm = re.match(r'"([^"]+)"', line)
        if mm:
            path = mm.group(1)
            out[path.rsplit("/", 1)[-1]] = path
    return out


def transform_app_body(
    body: str,
    call_input_expr: str,
    biz_expr: str,
    rets: str,
    *,
    inject_shop_helpers: bool,
) -> str:
    s = body[1:-1]
    s = re.sub(r"\bh\.logic\.", biz_expr + ".", s) if biz_expr else s
    s = s.replace("h.svcCtx.", "l.svcCtx.")
    s = s.replace("h.shopUser(", "shopUser(")
    s = s.replace("h.requirePerm(", "requirePerm(")
    s = s.replace("saveShopUpload(", "upload.SaveShopImage(")

    ci = call_input_expr.strip()
    preamble = f"\tin := {ci}\n"
    if inject_shop_helpers and ("shopUser(" in s or "requirePerm(" in s):
        preamble += SHOP_USER_HELPER + "\n"
        if "requirePerm(" in s:
            preamble += REQUIRE_PERM_HELPER + "\n"

    s = re.sub(r"\bvar req ", "var body ", s)
    s = re.sub(r"BindBody\(in, &req\)", "BindBody(in, &body)", s)
    s = re.sub(r", req\)", ", body)", s)
    s = re.sub(r"\(ctx, req\)", "(ctx, body)", s)
    s = re.sub(r"\(req\)", "(body)", s)
    s = re.sub(
        r"\breq\.(Reason|Password|UserIDs|Amount|Remark|Field|Action|Name|Code|EntryID|ProductID|Quantity|AutoRenew|PackageID|SlotID|ThemeID|ShopID|ProductIds|Stock|Price|StartAt|EndAt|Items|MenuIDs|Mobile|Nickname|RoleID)\b",
        r"body.\1",
        s,
    )

    inner = preamble + s

    # Multiline-safe return wraps
    if "*types.AnyResp" in rets:
        inner = re.sub(
            r"return nil, nil\b",
            "return &types.AnyResp{}, nil",
            inner,
        )
        # return expr, nil  where expr may contain braces — use balanced-ish heuristic
        def wrap_any_ret(mm: re.Match) -> str:
            expr = mm.group(1).strip()
            if expr == "nil":
                return "return &types.AnyResp{}, nil"
            return f"return &types.AnyResp{{Data: {expr}}}, nil"

        # Only simple single-line returns; multiline fixed in post
        inner = re.sub(r"return ([^,\n]+), nil\s*$", wrap_any_ret, inner, flags=re.M)
    elif "*types.PageListResp" in rets:
        inner = re.sub(
            r"return types\.PageListResp\{",
            "return &types.PageListResp{",
            inner,
        )
        inner = re.sub(
            r"return (map\[string\]interface\{\}\{.*?\}), nil",
            r"return &types.PageListResp{List: \1}, nil",
            inner,
            flags=re.S,
        )

    return inner


def detect_biz_expr(svc: str, app_import: str, recv: str) -> str:
    if svc == "merchant-service":
        return "biz.NewMerchantLogic(l.svcCtx)"
    if "product/app" in app_import:
        if recv == "FavoriteHandler":
            return "plogic.NewFavoriteLogic(l.svcCtx)"
        if recv == "ProductHandler":
            return "plogic.NewProductAdminLogic(l.svcCtx)"
        if recv == "PlatformProductHandler":
            return "plogic.NewPlatformProductLogic(l.svcCtx)"
        if recv == "ShopUploadHandler":
            return ""  # no domain logic
        return "plogic.NewCatalogLogic(l.svcCtx)"
    if "content/app" in app_import:
        return "clogic.NewArticleLogic(l.svcCtx)"
    if "notify/handler" in app_import:
        return "nlogic.NewNotificationLogic(l.svcCtx)"
    if "shopops/handler" in app_import:
        return ""  # uses l.svcCtx.ShopRBAC directly
    return ""


def fix_imports(new_text: str, svc: str) -> str:
    new_text = re.sub(r'\n\t\w+ "mymall/services/[^"]+/internal/app/[^"]+"', "", new_text)
    new_text = re.sub(
        r'\n\t\w+ "mymall/services/catalog-service/internal/product/app/[^"]+"', "", new_text
    )
    new_text = re.sub(
        r'\n\t\w+ "mymall/services/catalog-service/internal/content/app/[^"]+"', "", new_text
    )
    new_text = re.sub(
        r'\n\t\w+ "mymall/services/catalog-service/internal/notify/handler"', "", new_text
    )
    new_text = re.sub(
        r'\n\t\w+ "mymall/services/catalog-service/internal/shopops/handler"', "", new_text
    )

    def ensure(imp: str) -> None:
        nonlocal new_text
        if imp not in new_text:
            new_text = new_text.replace("import (", f"import (\n\t{imp}", 1)

    if svc == "merchant-service" and "biz.NewMerchantLogic" in new_text:
        ensure('"mymall/services/merchant-service/internal/biz"')
    if "plogic." in new_text:
        ensure('plogic "mymall/services/catalog-service/internal/product/logic"')
    if "clogic." in new_text:
        ensure('clogic "mymall/services/catalog-service/internal/content/logic"')
    if "nlogic." in new_text:
        ensure('nlogic "mymall/services/catalog-service/internal/notify/logic"')
    if "upload.SaveShopImage" in new_text:
        ensure('"mymall/services/catalog-service/internal/product/upload"')
    if "middleware." in new_text:
        ensure('"mymall/pkg/middleware"')
    if "jwt." in new_text:
        ensure('"mymall/pkg/jwt"')
    if "xerr." in new_text:
        ensure('"mymall/pkg/xerr"')
    if "http.Status" in new_text or "*http.Request" in new_text:
        ensure('"net/http"')
    if "appinput." in new_text:
        ensure('"mymall/pkg/appinput"')
    if "strconv." in new_text:
        ensure('"strconv"')
    if "json." in new_text:
        ensure('"encoding/json"')
    if re.search(r"\bfmt\.", new_text):
        ensure('"fmt"')
    if "url.Values" in new_text:
        ensure('"net/url"')
    if "io." in new_text:
        ensure('"io"')
    if "repository." in new_text and "product/repository" not in new_text:
        # may be product or notify or shopops — goimports will fix; add product as default if plogic
        if "plogic." in new_text or "ProductListFilter" in new_text:
            ensure('"mymall/services/catalog-service/internal/product/repository"')
        if "NotificationListFilter" in new_text:
            ensure('"mymall/services/catalog-service/internal/notify/repository"')
    if "model." in new_text and "product/model" not in new_text and "plogic." in new_text:
        ensure('"mymall/services/catalog-service/internal/product/model"')
    if "types.MerchantProductSaveReq" in new_text or "product/types" in new_text:
        # app used product/types — may conflict with go-zero types package alias
        if 'ptypes "' not in new_text and "MerchantProductSaveReq" in new_text:
            # bodies reference types.X from product/types — need alias
            pass

    new_text = re.sub(r"\n\t_ = fmt\.Sprintf\n", "\n", new_text)
    new_text = re.sub(r"\n\t_ = url\.Values\{\}\n", "\n", new_text)
    return new_text


def rewrite_logic_file(path: Path, app_index: dict[str, tuple[str, str, str]], svc: str) -> bool:
    text = path.read_text()
    if "appinput.CallInput{" not in text and ".CallInput{" not in text:
        return False

    imports = parse_imports(text)

    # Alias.NewCtor(l.svcCtx).Method / Alias.NewCtor().Method
    call_re = re.compile(
        r"(\w+)\.New(\w+)\((?:l\.svcCtx)?\)\.(\w+)\(\s*ctx\s*,\s*(appinput\.CallInput\{.*?\})\s*\)",
        re.S,
    )

    parts = re.split(r"(?=func \(l \*\w+\) )", text)
    new_parts = [parts[0]]
    changed = False
    for part in parts[1:]:
        m = re.match(
            r"func \(l \*(\w+)\) (\w+)\((.*?)\) \((.*?)\) \{",
            part,
            re.S,
        )
        if not m:
            new_parts.append(part)
            continue
        struct, method, params, rets = m.group(1), m.group(2), m.group(3), m.group(4)
        brace = part.find("{")
        end = find_matching(part, brace, "{", "}")
        old_body = part[brace : end + 1]
        rest = part[end + 1 :]

        cm = call_re.search(old_body)
        if not cm:
            new_parts.append(part)
            continue
        alias, ctor, app_method, cin = cm.group(1), cm.group(2), cm.group(3), cm.group(4)
        imp = imports.get(alias, "")
        recv = ctor  # NewProductHandler -> wait, group is without New: New(\w+) so ProductHandler
        # Actually New(\w+) captures ProductHandler from NewProductHandler
        key = f"{imp}|{recv}|{app_method}"
        # fallback: method-only for merchant
        if key not in app_index:
            key2 = app_method
            if key2 in app_index:
                key = key2
            else:
                new_parts.append(part)
                continue

        imp2, recv2, app_body = app_index[key]
        biz_expr = detect_biz_expr(svc, imp2, recv2)
        inject = "shopUser(" in app_body or "requirePerm(" in app_body or "h.shopUser" in app_body

        # product/types vs go-zero types: rewrite product types refs after transform
        new_inner = transform_app_body(
            app_body, cin, biz_expr, rets, inject_shop_helpers=inject
        )
        # Fix product domain types package collision: types.Merchant* from product/types
        if "product/app" in imp2 and "types." in new_inner:
            # Alias product/types as ptypes in bodies that bind MerchantProductSaveReq etc.
            for name in (
                "MerchantProductSaveReq",
                "ProductStatusReq",
                "BatchProductsReq",
                "AdjustStockReq",
                "BatchStockReq",
                "ScheduleReq",
                "ProductTagSaveReq",
                "AttrTemplateSaveReq",
            ):
                new_inner = new_inner.replace(f"types.{name}", f"ptypes.{name}")

        if "content/app" in imp2:
            for name in (
                "ArticleSaveReq",
                "ArticleCategorySaveReq",
                "CommentCreateReq",
                "BannerSaveReq",
                "EmojiSaveReq",
                "AuditReq",
                "BatchAuditReq",
                "CommentPatchReq",
            ):
                new_inner = new_inner.replace(f"types.{name}", f"ctypes.{name}")

        if "shopops/handler" in imp2:
            new_inner = new_inner.replace("types.", "sotypes.")

        new_body = "{\n" + new_inner + "\n}"
        new_parts.append(f"func (l *{struct}) {method}({params}) ({rets}) " + new_body + rest)
        changed = True

    if not changed:
        return False

    new_text = fix_imports("".join(new_parts), svc)
    if "ptypes." in new_text:
        if 'ptypes "' not in new_text:
            new_text = new_text.replace(
                "import (",
                'import (\n\tptypes "mymall/services/catalog-service/internal/product/types"',
                1,
            )
    if "ctypes." in new_text:
        if 'ctypes "' not in new_text:
            new_text = new_text.replace(
                "import (",
                'import (\n\tctypes "mymall/services/catalog-service/internal/content/types"',
                1,
            )
    if "sotypes." in new_text:
        if 'sotypes "' not in new_text:
            new_text = new_text.replace(
                "import (",
                'import (\n\tsotypes "mymall/services/catalog-service/internal/shopops/types"',
                1,
            )
    # shopops model
    if "model.ShopRole" in new_text or "model.ShopMenu" in new_text:
        if "shopops/model" not in new_text:
            new_text = new_text.replace(
                "import (",
                'import (\n\t"mymall/services/catalog-service/internal/shopops/model"',
                1,
            )
    if "repository.BuildShopMenuTree" in new_text:
        if "shopops/repository" not in new_text:
            new_text = new_text.replace(
                "import (",
                'import (\n\t"mymall/services/catalog-service/internal/shopops/repository"',
                1,
            )
    if "notify/repository" not in new_text and "NotificationListFilter" in new_text:
        new_text = new_text.replace(
            "import (",
            'import (\n\t"mymall/services/catalog-service/internal/notify/repository"',
            1,
        )

    path.write_text(new_text)
    return True


def index_app(svc: str) -> dict[str, tuple[str, str, str]]:
    base = ROOT / "services" / svc
    dirs: list[Path] = []
    if svc == "merchant-service":
        dirs = [base / "internal" / "app"]
    elif svc == "catalog-service":
        dirs = [
            base / "internal" / "product" / "app",
            base / "internal" / "content" / "app",
            base / "internal" / "notify" / "handler",
            base / "internal" / "shopops" / "handler",
        ]
    idx: dict[str, tuple[str, str, str]] = {}
    for d in dirs:
        if not d.is_dir():
            continue
        for go in d.rglob("*.go"):
            text = go.read_text()
            rel = go.parent.relative_to(base)
            imp = f"mymall/services/{svc}/{rel.as_posix()}"
            for name, (recv, body) in extract_methods(text).items():
                # unique key
                idx[f"{imp}|{recv}|{name}"] = (imp, recv, body)
                # also plain name for merchant uniqueness / fallback (last wins)
                idx[name] = (imp, recv, body)
    return idx


def main() -> int:
    services = sys.argv[1:] or ["merchant-service"]
    for svc in services:
        idx = index_app(svc)
        print(f"{svc}: indexed {len(idx)} keys")
        n = 0
        logic = ROOT / "services" / svc / "internal" / "logic"
        for go in sorted(logic.rglob("*.go")):
            if rewrite_logic_file(go, idx, svc):
                n += 1
                print(" ", go.relative_to(ROOT))
        print(f"{svc}: rewrote {n} logic files")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
