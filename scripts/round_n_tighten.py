#!/usr/bin/env python3
"""Round-N: catalog path+body + CallInput cleanup; merchant CallInput cleanup; uploads/purge."""
from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
CAT = ROOT / "services/catalog-service"
MER = ROOT / "services/merchant-service"


def find_matching(text: str, open_idx: int) -> int:
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
        if c == "{":
            depth += 1
        elif c == "}":
            depth -= 1
            if depth == 0:
                return i
        i += 1
    raise RuntimeError("unbalanced")


def strip_callinput_path(body: str) -> str:
    body = re.sub(
        r'\tin := appinput\.CallInput\{PathVars: map\[string\]string\{"id": fmt\.Sprintf\("%d", req\.Id\)\}(?:, Body: req)?\}\n\n?',
        "",
        body,
    )
    body = re.sub(
        r'\tid, (?:err|_) := strconv\.ParseUint\(in\.Path\("id"\), 10, 64\)\n(?:\tif err != nil(?: \|\| id == 0)? \{\n\t\treturn nil, xerr\.New\(http\.StatusBadRequest, "[^"]+"\)\n\t\}\n)?',
        "\tid := req.Id\n",
        body,
    )
    body = re.sub(
        r'\tid, _ := strconv\.ParseUint\(in\.Path\("id"\), 10, 64\)\n',
        "\tid := req.Id\n",
        body,
    )
    body = re.sub(
        r'\tappID, err := strconv\.ParseUint\(in\.Path\("id"\), 10, 64\)\n\tif err != nil \{\n\t\treturn nil, xerr\.New\(http\.StatusBadRequest, "[^"]+"\)\n\t\}\n',
        "\tappID := req.Id\n",
        body,
    )
    body = re.sub(
        r'\tshopID, err := strconv\.ParseUint\(in\.Path\("id"\), 10, 64\)\n\tif err != nil \{\n\t\treturn nil, xerr\.New\(http\.StatusBadRequest, "[^"]+"\)\n\t\}\n',
        "\tshopID := req.Id\n",
        body,
    )
    return body


def strip_callinput_page(body: str) -> str:
    body = re.sub(
        r'\tin := appinput\.CallInput\{Query: url\.Values\{"page": \{fmt\.Sprintf\("%d", req\.Page\)\}, "page_size": \{fmt\.Sprintf\("%d", req\.PageSize\)\}\}\}\n\n?',
        "",
        body,
    )
    body = re.sub(r"\tp, ps := in\.Page\(\)\n", "\tp, ps := req.Page, req.PageSize\n", body)
    body = re.sub(r"\tpage, pageSize := in\.Page\(\)\n", "\tpage, pageSize := req.Page, req.PageSize\n", body)
    body = re.sub(
        r'\tpage, _ := strconv\.Atoi\(in\.QueryGet\("page"\)\)\n\tpageSize, _ := strconv\.Atoi\(in\.QueryGet\("page_size"\)\)\n',
        "\tpage, pageSize := int(req.Page), int(req.PageSize)\n",
        body,
    )
    # leftover QueryGet on empty/missing query → "" (status/keyword etc. already lost in CallInput bridge)
    body = re.sub(r'in\.QueryGet\("([^"]+)"\)', r'"" /* was query:\1 */', body)
    return body


def strip_callinput_empty_query(body: str) -> str:
    """Empty CallInput that only read query params — keep behavior via typed req when available."""
    body = re.sub(r"\tin := appinput\.CallInput\{\}\n\n?", "", body)
    body = re.sub(r'slotType := in\.QueryGet\("slot_type"\)', "slotType := req.SlotType", body)
    body = re.sub(
        r'shopID, _ := strconv\.ParseUint\(in\.QueryGet\("shop_id"\), 10, 64\)',
        "shopID := req.ShopId",
        body,
    )
    body = re.sub(r'in\.QueryGet\("([^"]+)"\)', r'"" /* was query:\1 */', body)
    return body


def strip_unused_imports(text: str) -> str:
    if "appinput." not in text:
        text = re.sub(r'\n\t"mymall/pkg/appinput"', "", text)
    if "fmt." not in text:
        text = re.sub(r'\n\t"fmt"', "", text)
    if "strconv." not in text:
        text = re.sub(r'\n\t"strconv"', "", text)
    if "url." not in text and "net/url" in text:
        text = re.sub(r'\n\t"net/url"', "", text)
    if "ptypes." not in text:
        text = re.sub(r'\n\tptypes "mymall/services/catalog-service/internal/product/types"', "", text)
    if "ctypes." not in text:
        text = re.sub(r'\n\tctypes "mymall/services/catalog-service/internal/content/types"', "", text)
    if "sotypes." not in text:
        text = re.sub(r'\n\tsotypes "mymall/services/catalog-service/internal/shopops/types"', "", text)
    return text


# method -> (new_req_type, kind)
CAT_BODY = {
    "MerchantUpdateProduct": ("ProductUpdateBodyReq", "product_update"),
    "MerchantSetProductStatus": ("SetStatusBodyReq", "set_status"),
    "MerchantUpdateTag": ("TagUpdateBodyReq", "tag_update"),
    "MerchantUpdateAttrTemplate": ("AttrTemplateUpdateBodyReq", "attr_update"),
    "AdjustStock": ("StockAdjustBodyReq", "stock"),
    "MerchantScheduleProduct": ("ScheduleBodyReq", "schedule"),
    "MerchantUpdateRole": ("ShopRoleUpdateBodyReq", "role_update"),
    "MerchantUpdateArticle": ("ArticleUpdateBodyReq", "article_update"),
    "MerchantPatchArticleComment": ("ArticleCommentPatchBodyReq", "comment_patch"),
    "AdminUpdateArticleCategory": ("ArticleCategoryUpdateBodyReq", "article_cat_update"),
    "AdminTopArticle": ("ArticleTopBodyReq", "article_top"),
    "AdminSoftDeleteArticle": ("ArticleRemarkBodyReq", "article_remark"),
    "AdminOfflineArticle": ("ArticleRemarkBodyReq", "article_remark"),
    "AdminUpdateArticle": ("ArticleUpdateBodyReq", "article_update"),
    "AdminAuditArticle": ("ArticleAuditBodyReq", "article_audit"),
    "UpdateMine": ("UserArticleUpdateBodyReq", "user_article_update"),
    "CreateComment": ("CreateCommentBodyReq", "create_comment"),
    "UpdateBanner": ("BannerUpdateBodyReq", "banner_update"),
    "AdminDeleteProduct": ("PlatformProductRemarkBodyReq", "platform_remark"),
    "AdminOffSaleProduct": ("PlatformProductRemarkBodyReq", "platform_remark"),
    "AdminUpdateCategory": ("CategoryUpdateBodyReq", "category_update"),
    "EmojiUpdate": ("EmojiUpdateBodyReq", "emoji_update"),
    "AdminPatchArticleComment": ("ArticleCommentPatchBodyReq", "comment_patch"),
    "MerchantPurgeProducts": ("RecycleReq", "purge_products"),
    "AdminPurgeArticleRecycle": ("ArticleIdListReq", "purge_article"),
}


def rewrite_cat_body(body: str, kind: str) -> str:
    body = strip_callinput_path(body)
    body = re.sub(r"\tin := appinput\.CallInput\{\}\n\n?", "", body)
    # remove BindBody blocks
    body = re.sub(
        r"\tvar body [^\n]+\n\t(?:if err := |_ = )appinput\.BindBody\(in, &body\)[^\n]*\n(?:\t\treturn nil, xerr\.New\(http\.StatusBadRequest, \"[^\"]+\"\)\n\t\}\n)?",
        "",
        body,
        count=1,
    )
    body = re.sub(
        r"\tvar body struct \{[^}]+\}\n\t(?:if err := |_ = )appinput\.BindBody\(in, &body\)[^\n]*\n(?:\t\treturn nil, xerr\.New\(http\.StatusBadRequest, \"[^\"]+\"\)\n\t\}\n)?",
        "",
        body,
        count=1,
        flags=re.S,
    )

    if kind == "product_update":
        body = re.sub(r"Save\(ctx, shopID, uid, id, body\)", "Save(ctx, shopID, uid, req.Id, req.ToProduct())", body)
        body = body.replace("Save(ctx, shopID, uid, id, body)", "Save(ctx, shopID, uid, req.Id, req.ToProduct())")
    elif kind == "set_status":
        body = re.sub(r"body\.Status", "req.Status", body)
        body = re.sub(r", body\)", ", req.ToProduct())", body)
    elif kind == "tag_update":
        body = re.sub(r"body\.Name", "req.Name", body)
        body = re.sub(r"body\.Color", "req.Color", body)
        body = re.sub(r"ID: id,", "ID: req.Id,", body)
    elif kind == "attr_update":
        body = re.sub(r"body\.Name", "req.Name", body)
        body = re.sub(r"body\.AttrsJSON", "req.AttrsJSON", body)
        body = re.sub(r"ID: id,", "ID: req.Id,", body)
    elif kind == "stock":
        body = re.sub(r"AdjustSkuStock\(ctx, shopID, body\)", "AdjustSkuStock(ctx, shopID, req.ToProduct())", body)
        body = re.sub(r", body\)", ", req.ToProduct())", body)
        body = re.sub(r"body\.SkuID", "req.SkuID", body)
    elif kind == "schedule":
        body = re.sub(r", body\)", ", req.ToProduct())", body)
        body = re.sub(r"body\.Action", "req.Action", body)
        body = re.sub(r"body\.RunAt", "req.RunAt", body)
    elif kind == "role_update":
        body = re.sub(r"body\.Code", "req.Code", body)
        body = re.sub(r"body\.Name", "req.Name", body)
        body = re.sub(r"body\.Remark", "req.Remark", body)
        body = re.sub(r"body\.MenuIDs", "req.MenuIDs", body)
        body = re.sub(r"ID: id,", "ID: req.Id,", body)
        body = re.sub(r"SaveRole\(ctx, role, req\.MenuIDs\)", "SaveRole(ctx, role, req.MenuIDs)", body)
    elif kind == "article_update":
        body = re.sub(r", body\)", ", req.ToContent())", body)
        body = re.sub(r"Update\(ctx, id,", "Update(ctx, req.Id,", body)
        body = re.sub(r"MerchantUpdate\(ctx, shopID, id,", "MerchantUpdate(ctx, shopID, req.Id,", body)
    elif kind == "comment_patch":
        body = re.sub(r", body\.Status\)", ", req.Status)", body)
        body = re.sub(r", body\)", ", req.ToContent())", body)
        body = re.sub(r"body\.Status", "req.Status", body)
    elif kind == "article_cat_update":
        body = re.sub(r", body\)", ", req.ToContent())", body)
    elif kind == "article_top":
        body = re.sub(r"body\.IsTop", "req.IsTop", body)
        body = re.sub(r", body\)", ", req.ToContent())", body)
    elif kind == "article_remark":
        body = re.sub(r"body\.Remark", "req.Remark", body)
    elif kind == "article_audit":
        body = re.sub(r", body\)", ", req.ToContent())", body)
        body = re.sub(r"body\.Pass", "req.Pass", body)
        body = re.sub(r"body\.RejectReason", "req.RejectReason", body)
    elif kind == "user_article_update":
        body = re.sub(
            r"ctypes\.ArticleSaveReq\{[^}]+\}",
            "req.ToContent()",
            body,
            flags=re.S,
        )
        body = re.sub(
            r"ArticleSaveReq\{\n\t\tCategoryID: (?:body|req)\.[^}]+\}",
            "req.ToContent()",
            body,
            flags=re.S,
        )
    elif kind == "create_comment":
        body = re.sub(
            r"CreatePublicComment\(ctx, userID, id, body\)",
            "CreatePublicComment(ctx, userID, req.Id, clogic.CreateCommentReq{Content: req.Content, ParentID: req.ParentID})",
            body,
        )
        body = re.sub(
            r"CreateComment\(ctx, userID, id, body\)",
            "CreateComment(ctx, userID, req.Id, clogic.CreateCommentReq{Content: req.Content, ParentID: req.ParentID})",
            body,
        )
        body = re.sub(r", body\)", ", clogic.CreateCommentReq{Content: req.Content, ParentID: req.ParentID})", body)
        body = re.sub(r"\tvar body logic\.CreateCommentReq\n", "", body)
    elif kind == "banner_update":
        body = re.sub(
            r"AdminUpdateBanner\(id, body\)",
            "AdminUpdateBanner(req.Id, clogic.BannerSaveReq{Title: req.Title, ImageURL: req.ImageURL, LinkType: req.LinkType, LinkID: req.LinkID, Sort: req.Sort, Status: req.Status, StartAt: req.StartAt, EndAt: req.EndAt})",
            body,
        )
        body = re.sub(r", body\)", ", clogic.BannerSaveReq{Title: req.Title, ImageURL: req.ImageURL, LinkType: req.LinkType, LinkID: req.LinkID, Sort: req.Sort, Status: req.Status, StartAt: req.StartAt, EndAt: req.EndAt})", body)
    elif kind == "platform_remark":
        body = re.sub(r"body\.Remark", "req.Remark", body)
    elif kind == "category_update":
        body = re.sub(r", body\)", ", req.ToProduct())", body)
        body = re.sub(r"body\.Name", "req.Name", body)
    elif kind == "emoji_update":
        body = re.sub(r"body\.Name", "req.Name", body)
        body = re.sub(r"body\.ImageURL", "req.ImageURL", body)
        body = re.sub(r"body\.Sort", "req.Sort", body)
        body = re.sub(r"body\.Status", "req.Status", body)
    elif kind == "purge_products":
        body = re.sub(r"body\.ProductIDs", "req.ProductIDs", body)
        body = re.sub(
            r"PermanentDelete\(ctx, shopID, uid, body\.ProductIDs\)",
            "PermanentDelete(ctx, shopID, uid, req.ProductIDs)",
            body,
        )
    elif kind == "purge_article":
        body = re.sub(r"body\.ID", "req.Id", body)
        body = re.sub(r"body\.Id", "req.Id", body)
    return body


def process_logic_file(path: Path, mapping: dict, is_catalog: bool) -> bool:
    text = path.read_text()
    parts = re.split(r"(?=func \(l \*\w+\) )", text)
    out = [parts[0]]
    changed = False
    for part in parts[1:]:
        m = re.match(r"func \(l \*(\w+)\) (\w+)\((.*?)\) \((.*?)\) \{", part, re.S)
        if not m:
            out.append(part)
            continue
        method = m.group(2)
        brace = part.find("{")
        end = find_matching(part, brace)
        body = part[brace + 1 : end]
        rest = part[end + 1 :]
        params = m.group(3)
        rets = m.group(4)

        if method in mapping:
            req_type, kind = mapping[method]
            if kind in ("purge_products", "purge_article"):
                new_params = f"ctx context.Context, req *types.{req_type}"
            else:
                new_params = re.sub(r"req \*types\.(?:IdPathReq|JSONBody)", f"req *types.{req_type}", params)
                if params.strip() == "ctx context.Context":
                    new_params = f"ctx context.Context, req *types.{req_type}"
                elif "req *types." + req_type not in new_params and "req *types." not in new_params:
                    # failed replace
                    new_params = re.sub(r"req \*types\.\w+", f"req *types.{req_type}", params)
            new_body = rewrite_cat_body(body, kind) if is_catalog else body
            out.append(f"func (l *{m.group(1)}) {method}({new_params}) ({rets}) {{{new_body}}}" + rest)
            changed = True
            continue

        # generic CallInput cleanup
        new_body = body
        if "appinput.CallInput" in body:
            if 'Query: url.Values{"page"' in body:
                new_body = strip_callinput_page(body)
            elif "PathVars:" in body and "BindBody" not in body:
                new_body = strip_callinput_path(body)
            elif "CallInput{}" in body and "BindBody" not in body and "Request" not in body:
                # empty CallInput with QueryGet - leave for typed query later
                pass
            elif "Request:" in body:
                # upload: replace in.Request with r
                new_body = re.sub(r"\tin := appinput\.CallInput\{Request: r\}\n\n?", "", body)
                new_body = new_body.replace("in.Request", "r")
        if new_body != body:
            out.append(f"func (l *{m.group(1)}) {method}({params}) ({rets}) {{{new_body}}}" + rest)
            changed = True
        else:
            out.append(part)

    if not changed:
        return False
    path.write_text(strip_unused_imports("".join(out)))
    return True


def update_handlers(handler_root: Path, mapping: dict) -> int:
    n = 0
    for p in handler_root.rglob("*_handler.go"):
        t = p.read_text()
        o = t
        for method, (req_type, kind) in mapping.items():
            hname = method + "Handler"
            if kind in ("purge_products", "purge_article"):
                # ensure parse RecycleReq / ArticleIdListReq
                pat = (
                    r"(func " + hname + r"\(svcCtx \*svc\.ServiceContext\) http\.HandlerFunc \{\n"
                    r"\treturn func\(w http\.ResponseWriter, r \*http\.Request\) \{\n)"
                    r"(?:\t\tvar req types\.\w+\n\t\tif err := httpx\.Parse\(r, &req\); err != nil \{\n\t\t\thttpx\.ErrorCtx\(r\.Context\(\), w, err\)\n\t\t\treturn\n\t\t\}\n\n)?"
                    r"(\t\tl := [^\n]+\n)"
                    r"\t\tresp, err := l\." + method + r"\(r\.Context\(\)(?:, &req)?\)"
                )
                repl = (
                    r"\1\t\tvar req types." + req_type + r"\n"
                    r"\t\tif err := httpx.Parse(r, &req); err != nil {\n"
                    r"\t\t\thttpx.ErrorCtx(r.Context(), w, err)\n"
                    r"\t\t\treturn\n"
                    r"\t\t}\n\n"
                    r"\2\t\tresp, err := l." + method + r"(r.Context(), &req)"
                )
                t2, c = re.subn(pat, repl, t, count=1)
                if c:
                    t = t2
                continue
            pat = (
                r"(func " + hname + r"\(svcCtx \*svc\.ServiceContext\) http\.HandlerFunc \{\n"
                r"\treturn func\(w http\.ResponseWriter, r \*http\.Request\) \{\n"
                r"\t\tvar req types\.)IdPathReq"
            )
            t2, c = re.subn(pat, r"\1" + req_type, t, count=1)
            if c:
                t = t2
        # upload handlers: drop JSONBody parse if still present
        for uh in (
            "UploadMine",
            "MerchantUploadImage",
            "MerchantUploadArticle",
            "AdminUploadArticle",
            "UploadBanner",
            "AdminUploadShop",
            "MerchantImportProducts",
        ):
            hname = uh + "Handler"
            pat = (
                r"(func " + hname + r"\(svcCtx \*svc\.ServiceContext\) http\.HandlerFunc \{\n"
                r"\treturn func\(w http\.ResponseWriter, r \*http\.Request\) \{\n)"
                r"\t\tvar req types\.JSONBody\n"
                r"\t\tif err := httpx\.Parse\(r, &req\); err != nil \{\n"
                r"\t\t\thttpx\.ErrorCtx\(r\.Context\(\), w, err\)\n"
                r"\t\t\treturn\n"
                r"\t\t\}\n\n"
            )
            t2, c = re.subn(pat, r"\1", t, count=1)
            if c:
                t = t2
        if t != o:
            p.write_text(t)
            n += 1
            print("handler", p.relative_to(ROOT))
    return n


def patch_catalog_api() -> None:
    api = next((CAT / "api").glob("*.api"))
    t = api.read_text()
    # ensure path+body type stubs exist in .api
    stubs = """
type ProductUpdateBodyReq {
	Id   uint64 `path:"id"`
	Name string `json:"name"`
}
type SetStatusBodyReq {
	Id     uint64 `path:"id"`
	Status string `json:"status"`
}
type TagUpdateBodyReq {
	Id    uint64 `path:"id"`
	Name  string `json:"name"`
	Color string `json:"color,optional"`
}
type AttrTemplateUpdateBodyReq {
	Id        uint64 `path:"id"`
	Name      string `json:"name"`
	AttrsJSON string `json:"attrs_json,optional"`
}
type StockAdjustBodyReq {
	Id    uint64 `path:"id"`
	SkuID uint64 `json:"sku_id"`
}
type ScheduleBodyReq {
	Id     uint64 `path:"id"`
	Action string `json:"action"`
	RunAt  string `json:"run_at"`
}
type ShopRoleUpdateBodyReq {
	Id   uint64 `path:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
type ArticleUpdateBodyReq {
	Id    uint64 `path:"id"`
	Title string `json:"title"`
}
type ArticleCommentPatchBodyReq {
	Id     uint64 `path:"id"`
	Status string `json:"status"`
}
type ArticleCategoryUpdateBodyReq {
	Id   uint64 `path:"id"`
	Name string `json:"name"`
}
type ArticleTopBodyReq {
	Id    uint64 `path:"id"`
	IsTop int8   `json:"is_top"`
}
type ArticleRemarkBodyReq {
	Id     uint64 `path:"id"`
	Remark string `json:"remark,optional"`
}
type ArticleAuditBodyReq {
	Id   uint64 `path:"id"`
	Pass bool   `json:"pass"`
}
type UserArticleUpdateBodyReq {
	Id    uint64 `path:"id"`
	Title string `json:"title"`
}
type CreateCommentBodyReq {
	Id      uint64 `path:"id"`
	Content string `json:"content"`
}
type BannerUpdateBodyReq {
	Id       uint64 `path:"id"`
	ImageURL string `json:"image_url"`
}
type PlatformProductRemarkBodyReq {
	Id     uint64 `path:"id"`
	Remark string `json:"remark,optional"`
}
type CategoryUpdateBodyReq {
	Id   uint64 `path:"id"`
	Name string `json:"name"`
}
type EmojiUpdateBodyReq {
	Id       uint64 `path:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}
type ArticleIdListReq {
	Id uint64 `json:"id"`
}
"""
    if "type ProductUpdateBodyReq" not in t:
        t = t.replace("type JSONBody {}", "type JSONBody {}\n" + stubs)
    # uploads: remove (JSONBody)
    for route in [
        "post /api/v1/user/article-uploads (JSONBody) returns (AnyResp)",
        "post /api/v1/merchant/products/import (JSONBody) returns (AnyResp)",
        "post /api/v1/merchant/uploads/images (JSONBody) returns (AnyResp)",
        "post /api/v1/merchant/article-uploads (JSONBody) returns (AnyResp)",
        "post /api/v1/admin/article-uploads (JSONBody) returns (AnyResp)",
        "post /api/v1/admin/banners/upload (JSONBody) returns (AnyResp)",
        "post /api/v1/admin/shop-uploads (JSONBody) returns (AnyResp)",
    ]:
        t = t.replace(route, route.replace(" (JSONBody)", ""))
    # purge
    t = t.replace(
        "delete /api/v1/merchant/products/recycle returns (AnyResp)",
        "delete /api/v1/merchant/products/recycle (RecycleReq) returns (AnyResp)",
    )
    t = t.replace(
        "delete /api/v1/admin/articles/recycle returns (AnyResp)",
        "delete /api/v1/admin/articles/recycle (ArticleIdListReq) returns (AnyResp)",
    )
    # path+body route replacements (IdPathReq -> BodyReq) — match by handler proximity via unique paths
    repls = [
        ("put /api/v1/merchant/products/:id (IdPathReq)", "put /api/v1/merchant/products/:id (ProductUpdateBodyReq)"),
        ("put /api/v1/merchant/products/:id/status (IdPathReq)", "put /api/v1/merchant/products/:id/status (SetStatusBodyReq)"),
        ("put /api/v1/merchant/tags/:id (IdPathReq)", "put /api/v1/merchant/tags/:id (TagUpdateBodyReq)"),
        ("put /api/v1/merchant/attr-templates/:id (IdPathReq)", "put /api/v1/merchant/attr-templates/:id (AttrTemplateUpdateBodyReq)"),
        ("put /api/v1/merchant/skus/:id/stock (IdPathReq)", "put /api/v1/merchant/skus/:id/stock (StockAdjustBodyReq)"),
        ("post /api/v1/merchant/products/:id/schedules (IdPathReq)", "post /api/v1/merchant/products/:id/schedules (ScheduleBodyReq)"),
        ("put /api/v1/merchant/shop/roles/:id (IdPathReq)", "put /api/v1/merchant/shop/roles/:id (ShopRoleUpdateBodyReq)"),
        ("put /api/v1/merchant/articles/:id (IdPathReq)", "put /api/v1/merchant/articles/:id (ArticleUpdateBodyReq)"),
        ("patch /api/v1/merchant/article-comments/:id (IdPathReq)", "patch /api/v1/merchant/article-comments/:id (ArticleCommentPatchBodyReq)"),
        ("put /api/v1/admin/article-categories/:id (IdPathReq)", "put /api/v1/admin/article-categories/:id (ArticleCategoryUpdateBodyReq)"),
        ("post /api/v1/admin/articles/:id/top (IdPathReq)", "post /api/v1/admin/articles/:id/top (ArticleTopBodyReq)"),
        ("delete /api/v1/admin/articles/:id (IdPathReq)", "delete /api/v1/admin/articles/:id (ArticleRemarkBodyReq)"),
        ("post /api/v1/admin/articles/:id/offline (IdPathReq)", "post /api/v1/admin/articles/:id/offline (ArticleRemarkBodyReq)"),
        ("put /api/v1/admin/articles/:id (IdPathReq)", "put /api/v1/admin/articles/:id (ArticleUpdateBodyReq)"),
        ("post /api/v1/admin/articles/:id/audit (IdPathReq)", "post /api/v1/admin/articles/:id/audit (ArticleAuditBodyReq)"),
        ("put /api/v1/user/articles/:id (IdPathReq)", "put /api/v1/user/articles/:id (UserArticleUpdateBodyReq)"),
        ("post /api/v1/articles/:id/comments (IdPathReq)", "post /api/v1/articles/:id/comments (CreateCommentBodyReq)"),
        ("put /api/v1/admin/banners/:id (IdPathReq)", "put /api/v1/admin/banners/:id (BannerUpdateBodyReq)"),
        ("delete /api/v1/admin/products/:id (IdPathReq)", "delete /api/v1/admin/products/:id (PlatformProductRemarkBodyReq)"),
        ("put /api/v1/admin/products/:id/off_sale (IdPathReq)", "put /api/v1/admin/products/:id/off_sale (PlatformProductRemarkBodyReq)"),
        ("put /api/v1/admin/categories/:id (IdPathReq)", "put /api/v1/admin/categories/:id (CategoryUpdateBodyReq)"),
        ("put /api/v1/admin/comment-emojis/:id (IdPathReq)", "put /api/v1/admin/comment-emojis/:id (EmojiUpdateBodyReq)"),
        ("patch /api/v1/admin/article-comments/:id (IdPathReq)", "patch /api/v1/admin/article-comments/:id (ArticleCommentPatchBodyReq)"),
    ]
    # fix home-slots - may be different path in catalog
    for a, b in repls:
        if a in t:
            t = t.replace(a, b)
        else:
            print("api miss", a)
    api.write_text(t)
    print("catalog api patched", api)


def process_merchant_callinput() -> int:
    n = 0
    for p in (MER / "internal/logic").rglob("*_logic.go"):
        text = p.read_text()
        if "appinput." not in text:
            continue
        parts = re.split(r"(?=func \(l \*\w+\) )", text)
        out = [parts[0]]
        changed = False
        for part in parts[1:]:
            m = re.match(r"func \(l \*(\w+)\) (\w+)\((.*?)\) \((.*?)\) \{", part, re.S)
            if not m:
                out.append(part)
                continue
            method = m.group(2)
            brace = part.find("{")
            end = find_matching(part, brace)
            body = part[brace + 1 : end]
            rest = part[end + 1 :]
            params = m.group(3)
            new_body = body
            if 'Query: url.Values{"page"' in body:
                new_body = strip_callinput_page(body)
            elif "PathVars:" in body and "BindBody" not in body:
                new_body = strip_callinput_path(body)
            elif re.search(r"CallInput\{\}\s*\n", body) and "BindBody" not in body:
                new_body = strip_callinput_empty_query(body)
                if method == "PublicHomeSlots":
                    params = "ctx context.Context, req *types.SlotTypeQueryReq"
                elif method == "PublicCouponCenter":
                    params = "ctx context.Context, req *types.ShopIdQueryReq"
            if new_body != body:
                out.append(
                    f"func (l *{m.group(1)}) {method}({params}) ({m.group(4)}) {{{new_body}}}" + rest
                )
                changed = True
            else:
                out.append(part)
        if changed:
            p.write_text(strip_unused_imports("".join(out)))
            n += 1
            print("merchant", p.relative_to(ROOT))
    return n


def patch_merchant_query_types() -> None:
    biz = MER / "internal/types/biz_types.go"
    t = biz.read_text()
    if "type SlotTypeQueryReq struct" not in t:
        t += """
type SlotTypeQueryReq struct {
	SlotType string `form:"slot_type,optional"`
}

type ShopIdQueryReq struct {
	ShopId uint64 `form:"shop_id,optional"`
}
"""
        biz.write_text(t)
        print("merchant biz_types query DTOs")
    api = next((MER / "api").glob("*.api"))
    at = api.read_text()
    at = at.replace(
        "get /api/v1/shops/home-slots returns (AnyResp)",
        "get /api/v1/shops/home-slots (SlotTypeQueryReq) returns (AnyResp)",
    )
    at = at.replace(
        "get /api/v1/coupons/center returns (AnyResp)",
        "get /api/v1/coupons/center (ShopIdQueryReq) returns (AnyResp)",
    )
    # find actual coupon center path
    if "ShopIdQueryReq" not in at and "coupons/center" in at:
        pass
    api.write_text(at)
    # handlers
    for p in (MER / "internal/handler").rglob("*_handler.go"):
        t = p.read_text()
        o = t
        if "PublicHomeSlotsHandler" in t and "SlotTypeQueryReq" not in t:
            t = re.sub(
                r"(func PublicHomeSlotsHandler\(svcCtx \*svc\.ServiceContext\) http\.HandlerFunc \{\n"
                r"\treturn func\(w http\.ResponseWriter, r \*http\.Request\) \{\n)"
                r"(\t\tl := [^\n]+\n)"
                r"\t\tresp, err := l\.PublicHomeSlots\(r\.Context\(\)\)",
                r"\1\t\tvar req types.SlotTypeQueryReq\n"
                r"\t\tif err := httpx.Parse(r, &req); err != nil {\n"
                r"\t\t\thttpx.ErrorCtx(r.Context(), w, err)\n"
                r"\t\t\treturn\n"
                r"\t\t}\n\n"
                r"\2\t\tresp, err := l.PublicHomeSlots(r.Context(), &req)",
                t,
                count=1,
            )
        if "PublicCouponCenterHandler" in t and "ShopIdQueryReq" not in t:
            t = re.sub(
                r"(func PublicCouponCenterHandler\(svcCtx \*svc\.ServiceContext\) http\.HandlerFunc \{\n"
                r"\treturn func\(w http\.ResponseWriter, r \*http\.Request\) \{\n)"
                r"(\t\tl := [^\n]+\n)"
                r"\t\tresp, err := l\.PublicCouponCenter\(r\.Context\(\)\)",
                r"\1\t\tvar req types.ShopIdQueryReq\n"
                r"\t\tif err := httpx.Parse(r, &req); err != nil {\n"
                r"\t\t\thttpx.ErrorCtx(r.Context(), w, err)\n"
                r"\t\t\treturn\n"
                r"\t\t}\n\n"
                r"\2\t\tresp, err := l.PublicCouponCenter(r.Context(), &req)",
                t,
                count=1,
            )
            if "mymall/services/merchant-service/internal/types" not in t:
                t = t.replace(
                    '\n\t"mymall/services/merchant-service/internal/svc"\n',
                    '\n\t"mymall/services/merchant-service/internal/svc"\n\t"mymall/services/merchant-service/internal/types"\n',
                )
            if "httpx" not in t:
                t = t.replace(
                    '"github.com/zeromicro/go-zero/rest/httpx"\n',
                    '"github.com/zeromicro/go-zero/rest/httpx"\n',
                )
                if '"github.com/zeromicro/go-zero/rest/httpx"' not in t:
                    t = t.replace(
                        "\n\t\"net/http\"\n",
                        '\n\t"net/http"\n\n\t"github.com/zeromicro/go-zero/rest/httpx"\n',
                    )
        if t != o:
            p.write_text(t)
            print("merchant handler", p.relative_to(ROOT))


def main() -> None:
    patch_catalog_api()
    cn = 0
    for p in sorted((CAT / "internal/logic").rglob("*_logic.go")):
        if process_logic_file(p, CAT_BODY, True):
            cn += 1
            print("catalog logic", p.relative_to(ROOT))
    print("catalog logic files", cn)
    hn = update_handlers(CAT / "internal/handler", CAT_BODY)
    print("catalog handlers", hn)
    patch_merchant_query_types()
    mn = process_merchant_callinput()
    print("merchant callinput files", mn)


if __name__ == "__main__":
    main()
