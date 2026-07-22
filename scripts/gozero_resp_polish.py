#!/usr/bin/env python3
"""Polish: drop JSONBody; empty ops → EmptyResp; list maps → PageListResp; typed data wrappers."""
from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SERVICES = ["catalog-service", "merchant-service", "order-service"]


def rm_jsonbody() -> None:
    for svc in SERVICES:
        api = next((ROOT / "services" / svc / "api").glob("*.api"), None)
        if api:
            t = api.read_text()
            t2 = re.sub(r"\ntype JSONBody \{\}\n?", "\n", t)
            if t2 != t:
                api.write_text(t2)
                print("api -JSONBody", api.relative_to(ROOT))
        types_go = ROOT / "services" / svc / "internal/types/types.go"
        if types_go.exists():
            t = types_go.read_text()
            t2 = re.sub(r"\ntype JSONBody struct \{\n\}\n?", "\n", t)
            if t2 != t:
                types_go.write_text(t2)
                print("types -JSONBody", types_go.relative_to(ROOT))


def convert_empty_logic(logic_root: Path) -> set[str]:
    """Return method names converted to EmptyResp."""
    methods: set[str] = set()
    empty_pat = re.compile(r"return &types\.AnyResp\{(?:Data: &types\.AnyResp\{\})?\}, nil")
    any_data_pat = re.compile(r"return &types\.AnyResp\{Data: (?!&types\.AnyResp\{\})")

    for p in logic_root.rglob("*_logic.go"):
        t = p.read_text()
        if "return &types.AnyResp{" not in t:
            continue
        parts = re.split(r"(?=func \(l \*\w+\) )", t)
        out = [parts[0]]
        changed = False
        for part in parts[1:]:
            m = re.match(r"func \(l \*(\w+)\) (\w+)\((.*?)\) \((.*?)\) \{", part, re.S)
            if not m:
                out.append(part)
                continue
            method = m.group(2)
            brace = part.find("{")
            depth = 0
            end = brace
            for i, c in enumerate(part[brace:], brace):
                if c == "{":
                    depth += 1
                elif c == "}":
                    depth -= 1
                    if depth == 0:
                        end = i
                        break
            body = part[brace + 1 : end]
            rest = part[end + 1 :]
            rets = m.group(4)
            has_empty = bool(empty_pat.search(body))
            has_data = bool(any_data_pat.search(body))
            if has_empty and not has_data and "*types.AnyResp" in rets:
                new_body = empty_pat.sub("return &types.EmptyResp{}, nil", body)
                new_rets = rets.replace("*types.AnyResp", "*types.EmptyResp")
                out.append(f"func (l *{m.group(1)}) {method}({m.group(3)}) ({new_rets}) {{{new_body}}}" + rest)
                methods.add(method)
                changed = True
            else:
                out.append(part)
        if changed:
            p.write_text("".join(out))
            print("empty→EmptyResp", p.relative_to(ROOT))
    return methods


def patch_api_empty(api: Path, methods: set[str]) -> None:
    if not methods:
        return
    t = api.read_text()
    o = t
    for method in sorted(methods):
        # @handler Method\n\tVERB path ... returns (AnyResp)
        pat = rf"(@handler {re.escape(method)}\n\t[^\n]+ returns )\(AnyResp\)"
        t2, n = re.subn(pat, r"\1(EmptyResp)", t)
        if n:
            t = t2
        else:
            print("api miss EmptyResp", method)
    if t != o:
        api.write_text(t)
        print("api EmptyResp", api.relative_to(ROOT))


def fix_merchant_page_lists() -> None:
    mapping = {
        "MerchantCouponClaims": ("IdPageReq", "PageListResp"),
        "MerchantCouponRedeems": ("IdPageReq", "PageListResp"),
        "AdminCouponClaims": ("IdPageReq", "PageListResp"),
        "AdminCouponRedeems": ("IdPageReq", "PageListResp"),
    }
    api = next((ROOT / "services/merchant-service/api").glob("*.api"))
    t = api.read_text()
    for method, (_, resp) in mapping.items():
        pat = rf"(@handler {method}\n\t[^\n]+ returns )\(AnyResp\)"
        t = re.sub(pat, rf"\1({resp})", t)
    api.write_text(t)

    for rel, method in [
        ("internal/logic/merchant/coupon/merchant_coupon_claims_logic.go", "MerchantCouponClaims"),
        ("internal/logic/merchant/coupon/merchant_coupon_redeems_logic.go", "MerchantCouponRedeems"),
        ("internal/logic/admin/coupon/admin_coupon_claims_logic.go", "AdminCouponClaims"),
        ("internal/logic/admin/coupon/admin_coupon_redeems_logic.go", "AdminCouponRedeems"),
    ]:
        p = ROOT / "services/merchant-service" / rel
        text = p.read_text()
        text = text.replace(
            "func (l *%sLogic) %s(ctx context.Context, req *types.IdPageReq) (resp *types.AnyResp, err error)"
            % (method, method),
            "func (l *%sLogic) %s(ctx context.Context, req *types.IdPageReq) (resp *types.PageListResp, err error)"
            % (method, method),
        )
        # struct name may differ - do generic
        text = re.sub(
            rf"(func \(l \*\w+\) {method}\(ctx context\.Context, req \*types\.IdPageReq\) )\(resp \*types\.AnyResp, err error\)",
            r"\1(resp *types.PageListResp, err error)",
            text,
        )
        text = re.sub(
            r'return &types\.AnyResp\{Data: map\[string\]interface\{\}\{"list": list, "total": total\}\}, nil',
            "return &types.PageListResp{List: list, Total: total}, nil",
            text,
        )
        p.write_text(text)
        print("PageListResp", p.relative_to(ROOT))


def add_typed_data_helpers() -> None:
    """Add common response data types; rewrite map literals to typed structs (still in AnyResp)."""
    cat_biz = ROOT / "services/catalog-service/internal/types/biz_types.go"
    t = cat_biz.read_text()
    if "type URLData struct" not in t:
        t += """
// ---- typed response payloads (used inside AnyResp.Data) ----

type URLData struct {
	Url string `json:"url"`
}

type CountData struct {
	Count int64 `json:"count"`
}

type FavoriteStatusData struct {
	Favorited bool `json:"favorited"`
}

type EngagementData struct {
	Liked     bool `json:"liked"`
	Favorited bool `json:"favorited"`
}
"""
        cat_biz.write_text(t)
        print("+typed data catalog")

    # rewrite catalog map returns
    repls = [
        (
            r'return &types\.AnyResp\{Data: map\[string\]string\{"url": url\}\}, nil',
            "return &types.AnyResp{Data: types.URLData{Url: url}}, nil",
        ),
        (
            r'return &types\.AnyResp\{Data: map\[string\]int64\{"count": n\}\}, nil',
            "return &types.AnyResp{Data: types.CountData{Count: n}}, nil",
        ),
        (
            r'return &types\.AnyResp\{Data: map\[string\]bool\{"favorited": okFav\}\}, nil',
            "return &types.AnyResp{Data: types.FavoriteStatusData{Favorited: okFav}}, nil",
        ),
        (
            r'return &types\.AnyResp\{Data: map\[string\]bool\{"liked": liked, "favorited": favorited\}\}, nil',
            "return &types.AnyResp{Data: types.EngagementData{Liked: liked, Favorited: favorited}}, nil",
        ),
    ]
    for p in (ROOT / "services/catalog-service/internal/logic").rglob("*_logic.go"):
        text = p.read_text()
        o = text
        for a, b in repls:
            text = re.sub(a, b, text)
        if text != o:
            p.write_text(text)
            print("typed data", p.relative_to(ROOT))

    # merchant coupon center list
    mer_biz = ROOT / "services/merchant-service/internal/types/biz_types.go"
    mt = mer_biz.read_text()
    if "type ListData struct" not in mt:
        mt += """
type ListData struct {
	List interface{} `json:"list"`
}
"""
        mer_biz.write_text(mt)
    p = ROOT / "services/merchant-service/internal/logic/public/coupon/public_coupon_center_logic.go"
    if p.exists():
        text = p.read_text()
        text2 = text.replace(
            'return &types.AnyResp{Data: map[string]interface{}{"list": list}}, nil',
            "return &types.AnyResp{Data: types.ListData{List: list}}, nil",
        )
        if text2 != text:
            p.write_text(text2)
            print("typed ListData", p.relative_to(ROOT))


def expand_catalog_api_body_stubs() -> None:
    """Align a few skinny .api body types with biz_types field names (documentation sync)."""
    api = next((ROOT / "services/catalog-service/api").glob("*.api"))
    t = api.read_text()
    # ProductUpdateBodyReq already stubby — expand if minimal
    old = """type ProductUpdateBodyReq {
	Id   uint64 `path:"id"`
	Name string `json:"name"`
}"""
    new = """type ProductUpdateBodyReq {
	Id          uint64 `path:"id"`
	Name        string `json:"name"`
	CategoryId  uint64 `json:"category_id,optional"`
	ProductType string `json:"product_type,optional"`
	Status      string `json:"status,optional"`
}"""
    if old in t:
        t = t.replace(old, new)
        api.write_text(t)
        print("expanded ProductUpdateBodyReq api")


def main() -> None:
    rm_jsonbody()
    cat_methods = convert_empty_logic(ROOT / "services/catalog-service/internal/logic")
    mer_methods = convert_empty_logic(ROOT / "services/merchant-service/internal/logic")
    patch_api_empty(next((ROOT / "services/catalog-service/api").glob("*.api")), cat_methods)
    patch_api_empty(next((ROOT / "services/merchant-service/api").glob("*.api")), mer_methods)
    fix_merchant_page_lists()
    add_typed_data_helpers()
    expand_catalog_api_body_stubs()


if __name__ == "__main__":
    main()
