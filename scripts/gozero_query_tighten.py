#!/usr/bin/env python3
"""One-shot: restore typed list/query DTOs + IdPageReq per go-zero conventions."""
from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
CAT = ROOT / "services/catalog-service"
MER = ROOT / "services/merchant-service"
ORD = ROOT / "services/order-service"

CAT_QUERY_TYPES = r'''
// ---- list / query (go-zero form tags) ----

type IdPageReq struct {
	Id       uint64 `path:"id"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

type PublicProductListReq struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
	ShopId     uint64 `form:"shop_id,optional"`
	CategoryId uint64 `form:"category_id,optional"`
	OrderBy    string `form:"order_by,optional"`
}

type PublicArticleListReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Home     string `form:"home,optional"`
}

type MerchantProductListReq struct {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	Name        string `form:"name,optional"`
	ProductNo   string `form:"product_no,optional"`
	CategoryId  uint64 `form:"category_id,optional"`
	Status      string `form:"status,optional"`
	ProductType string `form:"product_type,optional"`
	StockWarn   string `form:"stock_warn,optional"`
	OrderBy     string `form:"order_by,optional"`
	Recycle     string `form:"recycle,optional"`
}

type AdminProductListReq struct {
	Page         int    `form:"page,default=1"`
	PageSize     int    `form:"page_size,default=10"`
	ShopId       uint64 `form:"shop_id,optional"`
	Name         string `form:"name,optional"`
	ProductNo    string `form:"product_no,optional"`
	CategoryId   uint64 `form:"category_id,optional"`
	Status       string `form:"status,optional"`
	ProductType  string `form:"product_type,optional"`
	OrderBy      string `form:"order_by,optional"`
	CreatedFrom  string `form:"created_from,optional"`
	CreatedTo    string `form:"created_to,optional"`
	PublishFrom  string `form:"publish_from,optional"`
	PublishTo    string `form:"publish_to,optional"`
}

type MerchantArticleListReq struct {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	Title       string `form:"title,optional"`
	AuditStatus string `form:"audit_status,optional"`
	Status      string `form:"status,optional"`
}

type AdminArticleListReq struct {
	Page         int    `form:"page,default=1"`
	PageSize     int    `form:"page_size,default=10"`
	Title        string `form:"title,optional"`
	AuditStatus  string `form:"audit_status,optional"`
	Status       string `form:"status,optional"`
	ShopId       uint64 `form:"shop_id,optional"`
	HasSchedule  string `form:"has_schedule,optional"`
	CreatedFrom  string `form:"created_from,optional"`
	CreatedTo    string `form:"created_to,optional"`
}

type ArticleCommentListReq struct {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
	ArticleId uint64 `form:"article_id,optional"`
	ShopId    uint64 `form:"shop_id,optional"`
	Status    string `form:"status,optional"`
}

type NotificationListReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	IsRead   string `form:"is_read,optional"`
}

type OpLogsReq struct {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
	ProductId uint64 `form:"product_id,optional"`
}
'''

MER_QUERY_TYPES = r'''
// ---- list / query (go-zero form tags) ----

type IdPageReq struct {
	Id       uint64 `path:"id"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

type StatusPageReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status,optional"`
}

type CouponListReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status,optional"`
	Keyword  string `form:"keyword,optional"`
}

type SlotTypePageReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	SlotType string `form:"slot_type,optional"`
}

type SlotOrderListReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	ShopId   uint64 `form:"shop_id,optional"`
	SlotType string `form:"slot_type,optional"`
	Status   string `form:"status,optional"`
}

type ThemePackageListReq struct {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	ThemeSlotId uint64 `form:"theme_slot_id,optional"`
}

type ThemeOrderListReq struct {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	ShopId      uint64 `form:"shop_id,optional"`
	ThemeSlotId uint64 `form:"theme_slot_id,optional"`
}

type ShopListReq struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status,optional"`
	Name     string `form:"name,optional"`
}

type SeckillEntryListReq struct {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
	SessionId uint64 `form:"session_id,optional"`
}
'''

CAT_API_STUBS = '''
type IdPageReq {
	Id       uint64 `path:"id"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}
type PublicProductListReq {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
	ShopId     uint64 `form:"shop_id,optional"`
	CategoryId uint64 `form:"category_id,optional"`
	OrderBy    string `form:"order_by,optional"`
}
type PublicArticleListReq {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Home     string `form:"home,optional"`
}
type MerchantProductListReq {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	Name        string `form:"name,optional"`
	ProductNo   string `form:"product_no,optional"`
	CategoryId  uint64 `form:"category_id,optional"`
	Status      string `form:"status,optional"`
	ProductType string `form:"product_type,optional"`
	StockWarn   string `form:"stock_warn,optional"`
	OrderBy     string `form:"order_by,optional"`
	Recycle     string `form:"recycle,optional"`
}
type AdminProductListReq {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	ShopId      uint64 `form:"shop_id,optional"`
	Name        string `form:"name,optional"`
	ProductNo   string `form:"product_no,optional"`
	CategoryId  uint64 `form:"category_id,optional"`
	Status      string `form:"status,optional"`
	ProductType string `form:"product_type,optional"`
	OrderBy     string `form:"order_by,optional"`
	CreatedFrom string `form:"created_from,optional"`
	CreatedTo   string `form:"created_to,optional"`
	PublishFrom string `form:"publish_from,optional"`
	PublishTo   string `form:"publish_to,optional"`
}
type MerchantArticleListReq {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	Title       string `form:"title,optional"`
	AuditStatus string `form:"audit_status,optional"`
	Status      string `form:"status,optional"`
}
type AdminArticleListReq {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	Title       string `form:"title,optional"`
	AuditStatus string `form:"audit_status,optional"`
	Status      string `form:"status,optional"`
	ShopId      uint64 `form:"shop_id,optional"`
	HasSchedule string `form:"has_schedule,optional"`
	CreatedFrom string `form:"created_from,optional"`
	CreatedTo   string `form:"created_to,optional"`
}
type ArticleCommentListReq {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
	ArticleId uint64 `form:"article_id,optional"`
	ShopId    uint64 `form:"shop_id,optional"`
	Status    string `form:"status,optional"`
}
type NotificationListReq {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	IsRead   string `form:"is_read,optional"`
}
type OpLogsReq {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
	ProductId uint64 `form:"product_id,optional"`
}
'''

MER_API_STUBS = '''
type IdPageReq {
	Id       uint64 `path:"id"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}
type StatusPageReq {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status,optional"`
}
type CouponListReq {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status,optional"`
	Keyword  string `form:"keyword,optional"`
}
type SlotTypePageReq {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	SlotType string `form:"slot_type,optional"`
}
type SlotOrderListReq {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	ShopId   uint64 `form:"shop_id,optional"`
	SlotType string `form:"slot_type,optional"`
	Status   string `form:"status,optional"`
}
type ThemePackageListReq {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	ThemeSlotId uint64 `form:"theme_slot_id,optional"`
}
type ThemeOrderListReq {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	ShopId      uint64 `form:"shop_id,optional"`
	ThemeSlotId uint64 `form:"theme_slot_id,optional"`
}
type ShopListReq {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Status   string `form:"status,optional"`
	Name     string `form:"name,optional"`
}
type SeckillEntryListReq {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
	SessionId uint64 `form:"session_id,optional"`
}
'''

# method -> new req type
CAT_HANDLER_MAP = {
    "GetProductList": "PublicProductListReq",
    "PublicListArticles": "PublicArticleListReq",
    "ListComments": "IdPageReq",
    "MerchantListProducts": "MerchantProductListReq",
    "OpLogs": "OpLogsReq",
    "MerchantListArticles": "MerchantArticleListReq",
    "MerchantListArticleComments": "ArticleCommentListReq",
    "MerchantListNotifications": "NotificationListReq",
    "AdminListProducts": "AdminProductListReq",
    "AdminListArticles": "AdminArticleListReq",
    "AdminListArticleRecycle": "AdminArticleListReq",
    "AdminListArticleComments": "ArticleCommentListReq",
    "AdminListUserFavorites": "IdPageReq",
}

MER_HANDLER_MAP = {
    "ListMyCoupons": "StatusPageReq",
    "MerchantListCoupons": "CouponListReq",
    "AdminListCoupons": "CouponListReq",
    "MerchantListSlotPackages": "SlotTypePageReq",
    "AdminListSlotPackages": "SlotTypePageReq",
    "MerchantListSlotOrders": "SlotOrderListReq",
    "AdminListSlotOrders": "SlotOrderListReq",
    "MerchantListThemePackages": "ThemePackageListReq",
    "AdminListThemePackages": "ThemePackageListReq",
    "AdminListThemeOrders": "ThemeOrderListReq",
    "AdminListApplications": "StatusPageReq",
    "AdminListShops": "ShopListReq",
    "AdminListSeckillEntries": "SeckillEntryListReq",
    "MerchantCouponClaims": "IdPageReq",
    "MerchantCouponRedeems": "IdPageReq",
    "AdminCouponClaims": "IdPageReq",
    "AdminCouponRedeems": "IdPageReq",
    "AdminWalletLogs": "IdPageReq",
}


def ensure_types(path: Path, blob: str, marker: str) -> None:
    t = path.read_text()
    if marker in t:
        print("types ok", path.relative_to(ROOT))
        return
    path.write_text(t.rstrip() + "\n" + blob)
    print("types+", path.relative_to(ROOT))


def patch_api(api: Path, stubs: str, marker: str, route_repls: list[tuple[str, str]]) -> None:
    t = api.read_text()
    if marker not in t:
        # insert after JSONBody or PageReq block
        if "type JSONBody {}" in t:
            t = t.replace("type JSONBody {}", "type JSONBody {}\n" + stubs, 1)
        else:
            t = stubs + "\n" + t
    for a, b in route_repls:
        if a in t:
            t = t.replace(a, b)
        else:
            print("api miss", a)
    api.write_text(t)
    print("api", api.relative_to(ROOT))


def update_handlers(root: Path, mapping: dict[str, str]) -> None:
    for p in root.rglob("*_handler.go"):
        t = p.read_text()
        o = t
        for method, req in mapping.items():
            h = method + "Handler"
            # replace types.PageReq / IdPathReq in that handler only
            pat = (
                rf"(func {h}\(svcCtx \*svc\.ServiceContext\) http\.HandlerFunc \{{\n"
                rf"\treturn func\(w http\.ResponseWriter, r \*http\.Request\) \{{\n"
                rf"\t\tvar req types\.)(?:PageReq|IdPathReq)"
            )
            t2, n = re.subn(pat, rf"\1{req}", t, count=1)
            if n:
                t = t2
                continue
            # handler without parse (IdPath only previously missing) — inject parse for IdPageReq
            if req == "IdPageReq":
                pat2 = (
                    rf"(func {h}\(svcCtx \*svc\.ServiceContext\) http\.HandlerFunc \{{\n"
                    rf"\treturn func\(w http\.ResponseWriter, r \*http\.Request\) \{{\n)"
                    rf"(\t\tl := [^\n]+\n)"
                    rf"\t\tresp, err := l\.{method}\(r\.Context\(\)(?:, &req)?\)"
                )
                # only if no var req yet
                if f"func {h}" in t and f"l.{method}(r.Context())" in t and "var req" not in t[t.find(f"func {h}"):t.find(f"func {h}")+800]:
                    t2, n = re.subn(
                        pat2,
                        rf"\1\t\tvar req types.{req}\n"
                        rf"\t\tif err := httpx.Parse(r, &req); err != nil {{\n"
                        rf"\t\t\thttpx.ErrorCtx(r.Context(), w, err)\n"
                        rf"\t\t\treturn\n"
                        rf"\t\t}}\n\n"
                        rf"\2\t\tresp, err := l.{method}(r.Context(), &req)",
                        t,
                        count=1,
                    )
                    if n:
                        t = t2
        if t != o:
            p.write_text(t)
            print("handler", p.relative_to(ROOT))


def write(path: Path, content: str) -> None:
    path.write_text(content)
    print("logic", path.relative_to(ROOT))


def main() -> None:
    ensure_types(CAT / "internal/types/biz_types.go", CAT_QUERY_TYPES, "type PublicProductListReq struct")
    ensure_types(MER / "internal/types/biz_types.go", MER_QUERY_TYPES, "type CouponListReq struct")

    cat_routes = [
        ("get /api/v1/products/list (PageReq)", "get /api/v1/products/list (PublicProductListReq)"),
        ("get /api/v1/articles/list (PageReq)", "get /api/v1/articles/list (PublicArticleListReq)"),
        ("get /api/v1/articles/:id/comments (IdPathReq)", "get /api/v1/articles/:id/comments (IdPageReq)"),
        ("get /api/v1/merchant/products (PageReq)", "get /api/v1/merchant/products (MerchantProductListReq)"),
        ("get /api/v1/merchant/products/op-logs (PageReq)", "get /api/v1/merchant/products/op-logs (OpLogsReq)"),
        ("get /api/v1/merchant/articles (PageReq)", "get /api/v1/merchant/articles (MerchantArticleListReq)"),
        ("get /api/v1/merchant/article-comments (PageReq)", "get /api/v1/merchant/article-comments (ArticleCommentListReq)"),
        ("get /api/v1/merchant/notifications (PageReq)", "get /api/v1/merchant/notifications (NotificationListReq)"),
        ("get /api/v1/admin/products (PageReq)", "get /api/v1/admin/products (AdminProductListReq)"),
        ("get /api/v1/admin/articles (PageReq)", "get /api/v1/admin/articles (AdminArticleListReq)"),
        ("get /api/v1/admin/articles/recycle (PageReq)", "get /api/v1/admin/articles/recycle (AdminArticleListReq)"),
        ("get /api/v1/admin/article-comments (PageReq)", "get /api/v1/admin/article-comments (ArticleCommentListReq)"),
        ("get /api/v1/admin/users/:id/favorites (IdPathReq)", "get /api/v1/admin/users/:id/favorites (IdPageReq)"),
    ]
    patch_api(next((CAT / "api").glob("*.api")), CAT_API_STUBS, "type PublicProductListReq", cat_routes)

    mer_routes = [
        ("get /api/v1/user/coupons (PageReq)", "get /api/v1/user/coupons (StatusPageReq)"),
        ("get /api/v1/merchant/coupons (PageReq)", "get /api/v1/merchant/coupons (CouponListReq)"),
        ("get /api/v1/admin/coupons (PageReq)", "get /api/v1/admin/coupons (CouponListReq)"),
        ("get /api/v1/merchant/homepage-packages (PageReq)", "get /api/v1/merchant/homepage-packages (SlotTypePageReq)"),
        ("get /api/v1/admin/homepage-packages (PageReq)", "get /api/v1/admin/homepage-packages (SlotTypePageReq)"),
        ("get /api/v1/merchant/homepage-orders (PageReq)", "get /api/v1/merchant/homepage-orders (SlotOrderListReq)"),
        ("get /api/v1/admin/homepage-orders (PageReq)", "get /api/v1/admin/homepage-orders (SlotOrderListReq)"),
        ("get /api/v1/merchant/theme-packages (PageReq)", "get /api/v1/merchant/theme-packages (ThemePackageListReq)"),
        ("get /api/v1/admin/theme-packages (PageReq)", "get /api/v1/admin/theme-packages (ThemePackageListReq)"),
        ("get /api/v1/admin/theme-orders (PageReq)", "get /api/v1/admin/theme-orders (ThemeOrderListReq)"),
        ("get /api/v1/admin/applications (PageReq)", "get /api/v1/admin/applications (StatusPageReq)"),
        ("get /api/v1/admin/shops (PageReq)", "get /api/v1/admin/shops (ShopListReq)"),
        ("get /api/v1/admin/seckill/entries (PageReq)", "get /api/v1/admin/seckill/entries (SeckillEntryListReq)"),
        ("get /api/v1/merchant/coupons/:id/claims (IdPathReq)", "get /api/v1/merchant/coupons/:id/claims (IdPageReq)"),
        ("get /api/v1/merchant/coupons/:id/redeems (IdPathReq)", "get /api/v1/merchant/coupons/:id/redeems (IdPageReq)"),
        ("get /api/v1/admin/coupons/:id/claims (IdPathReq)", "get /api/v1/admin/coupons/:id/claims (IdPageReq)"),
        ("get /api/v1/admin/coupons/:id/redeems (IdPathReq)", "get /api/v1/admin/coupons/:id/redeems (IdPageReq)"),
        ("get /api/v1/admin/shops/:id/wallet/logs (IdPathReq)", "get /api/v1/admin/shops/:id/wallet/logs (IdPageReq)"),
    ]
    patch_api(next((MER / "api").glob("*.api")), MER_API_STUBS, "type CouponListReq", mer_routes)

    # order upload
    oapi = next((ORD / "api").glob("*.api"))
    ot = oapi.read_text().replace(
        "post /api/v1/user/review-uploads (JSONBody) returns (AnyResp)",
        "post /api/v1/user/review-uploads returns (AnyResp)",
    )
    oapi.write_text(ot)
    print("order api upload")

    rewrite_catalog_logic()
    rewrite_merchant_logic()
    update_handlers(CAT / "internal/handler", CAT_HANDLER_MAP)
    update_handlers(MER / "internal/handler", MER_HANDLER_MAP)


def rewrite_catalog_logic() -> None:
    write(CAT / "internal/logic/public/product/get_product_list_logic.go", '''package product

import (
	"context"
	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductListLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductListLogic {
	return &GetProductListLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *GetProductListLogic) GetProductList(ctx context.Context, req *types.PublicProductListReq) (resp *types.PageListResp, err error) {
	pageReq := &pagination.PageReq{Page: req.Page, PageSize: req.PageSize}
	data, err := plogic.NewCatalogLogic(l.svcCtx).GetProductListFiltered(ctx, pageReq, req.ShopId, "on_sale", req.CategoryId, req.OrderBy)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, "查询失败")
	}
	return &types.PageListResp{List: data}, nil
}
''')

    write(CAT / "internal/logic/public/article/public_list_articles_logic.go", '''package article

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicListArticlesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicListArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicListArticlesLogic {
	return &PublicListArticlesLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *PublicListArticlesLogic) PublicListArticles(ctx context.Context, req *types.PublicArticleListReq) (resp *types.PageListResp, err error) {
	data, err := clogic.NewArticleLogic(l.svcCtx).PublicList(ctx, req.Page, req.PageSize, req.Home == "1")
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
''')

    write(CAT / "internal/logic/public/article/list_comments_logic.go", '''package article

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCommentsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *ListCommentsLogic) ListComments(ctx context.Context, req *types.IdPageReq) (resp *types.PageListResp, err error) {
	data, err := clogic.NewArticleLogic(l.svcCtx).PublicListComments(ctx, req.Id, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
''')

    write(CAT / "internal/logic/merchant/product/merchant_list_products_logic.go", '''package product

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/product/repository"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListProductsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListProductsLogic {
	return &MerchantListProductsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *MerchantListProductsLogic) MerchantListProducts(ctx context.Context, req *types.MerchantProductListReq) (resp *types.PageListResp, err error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	f := repository.ProductListFilter{
		ShopID: shopID, Name: req.Name, ProductNo: req.ProductNo,
		CategoryID: req.CategoryId, Status: req.Status, ProductType: req.ProductType,
		StockWarnOnly: req.StockWarn == "1",
		Page: req.Page, PageSize: req.PageSize, OrderBy: req.OrderBy,
		Recycle: req.Recycle == "1",
	}
	data, err := plogic.NewProductAdminLogic(l.svcCtx).List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
''')

    # op_logs — read existing shell
    p = CAT / "internal/logic/merchant/product/op_logs_logic.go"
    t = p.read_text()
    t = t.replace("req *types.PageReq", "req *types.OpLogsReq")
    t = re.sub(
        r'pid, _ := strconv\.ParseUint\("" /\* was query:product_id \*/, 10, 64\)',
        "pid := req.ProductId",
        t,
    )
    t = t.replace("page, pageSize := req.Page, req.PageSize", "page, pageSize := req.Page, req.PageSize")
    if "strconv." not in t:
        t = t.replace('\n\t"strconv"\n', "\n")
    p.write_text(t)
    print("logic", p.relative_to(ROOT))

    write(CAT / "internal/logic/admin/product/admin_list_products_logic.go", '''package product

import (
	"context"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/product/repository"
	"net/http"
	"time"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListProductsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListProductsLogic {
	return &AdminListProductsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminListProductsLogic) AdminListProducts(ctx context.Context, req *types.AdminProductListReq) (resp *types.PageListResp, err error) {
	f := repository.ProductListFilter{
		ShopID: req.ShopId, Name: req.Name, ProductNo: req.ProductNo,
		CategoryID: req.CategoryId, Status: req.Status, ProductType: req.ProductType,
		Page: req.Page, PageSize: req.PageSize, OrderBy: req.OrderBy,
		PlatformScope: true,
	}
	if req.CreatedFrom != "" {
		if t, err := time.ParseInLocation("2006-01-02", req.CreatedFrom, time.Local); err == nil {
			f.CreatedFrom = &t
		}
	}
	if req.CreatedTo != "" {
		if t, err := time.ParseInLocation("2006-01-02", req.CreatedTo, time.Local); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			f.CreatedTo = &end
		}
	}
	if req.PublishFrom != "" {
		if t, err := time.ParseInLocation("2006-01-02", req.PublishFrom, time.Local); err == nil {
			f.PublishFrom = &t
		}
	}
	if req.PublishTo != "" {
		if t, err := time.ParseInLocation("2006-01-02", req.PublishTo, time.Local); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			f.PublishTo = &end
		}
	}
	data, err := plogic.NewPlatformProductLogic(l.svcCtx).List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
''')

    # merchant articles
    p = CAT / "internal/logic/merchant/article/merchant_list_articles_logic.go"
    t = p.read_text()
    t = t.replace("req *types.PageReq", "req *types.MerchantArticleListReq")
    t = t.replace('Title: "" /* was query:title */,', "Title: req.Title,")
    t = t.replace('AuditStatus: "" /* was query:audit_status */,', "AuditStatus: req.AuditStatus,")
    t = t.replace('Status:      "" /* was query:status */,', "Status:      req.Status,")
    p.write_text(t)
    print("logic", p.relative_to(ROOT))

    p = CAT / "internal/logic/merchant/article/merchant_list_article_comments_logic.go"
    t = p.read_text()
    t = t.replace("req *types.PageReq", "req *types.ArticleCommentListReq")
    t = re.sub(
        r'articleID, _ := strconv\.ParseUint\("" /\* was query:article_id \*/, 10, 64\)',
        "articleID := req.ArticleId",
        t,
    )
    t = t.replace('Status: "" /* was query:status */,', "Status: req.Status,")
    if "strconv." not in t:
        t = t.replace('\n\t"strconv"\n', "\n")
    p.write_text(t)
    print("logic", p.relative_to(ROOT))

    write(CAT / "internal/logic/merchant/notification/merchant_list_notifications_logic.go", '''package notification

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	nlogic "mymall/services/catalog-service/internal/notify/logic"
	"mymall/services/catalog-service/internal/notify/repository"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListNotificationsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListNotificationsLogic {
	return &MerchantListNotificationsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *MerchantListNotificationsLogic) MerchantListNotifications(ctx context.Context, req *types.NotificationListReq) (resp *types.PageListResp, err error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	f := repository.NotificationListFilter{ShopID: shopID, Page: req.Page, PageSize: req.PageSize}
	if req.IsRead == "0" || req.IsRead == "1" {
		v := int8(0)
		if req.IsRead == "1" {
			v = 1
		}
		f.IsRead = &v
	}
	data, err := nlogic.NewNotificationLogic(l.svcCtx).List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
''')

    for name, recycle in [
        ("admin_list_articles_logic.go", False),
        ("admin_list_article_recycle_logic.go", True),
    ]:
        method = "AdminListArticles" if not recycle else "AdminListArticleRecycle"
        write(CAT / "internal/logic/admin/article" / name, f'''package article

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"mymall/services/catalog-service/internal/content/repository"
	"net/http"
	"time"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type {method}Logic struct {{
	logx.Logger
	svcCtx *svc.ServiceContext
}}

func New{method}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *{method}Logic {{
	return &{method}Logic{{Logger: logx.WithContext(ctx), svcCtx: svcCtx}}
}}

func (l *{method}Logic) {method}(ctx context.Context, req *types.AdminArticleListReq) (resp *types.PageListResp, err error) {{
	f := repository.ArticleListFilter{{
		Title: req.Title, AuditStatus: req.AuditStatus, Status: req.Status,
		Page: req.Page, PageSize: req.PageSize, Recycle: {str(recycle).lower()},
	}}
	if req.ShopId > 0 {{
		f.ShopID = req.ShopId
		f.FilterShop = true
	}}
	if req.HasSchedule == "1" {{
		v := true
		f.HasSchedule = &v
	}} else if req.HasSchedule == "0" {{
		v := false
		f.HasSchedule = &v
	}}
	if req.CreatedFrom != "" {{
		if t, err := time.ParseInLocation("2006-01-02", req.CreatedFrom, time.Local); err == nil {{
			f.CreatedFrom = &t
		}}
	}}
	if req.CreatedTo != "" {{
		if t, err := time.ParseInLocation("2006-01-02", req.CreatedTo, time.Local); err == nil {{
			end := t.Add(24*time.Hour - time.Second)
			f.CreatedTo = &end
		}}
	}}
	data, err := clogic.NewArticleLogic(l.svcCtx).List(ctx, f)
	if err != nil {{
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}}
	return &types.PageListResp{{List: data}}, nil
}}
''')

    p = CAT / "internal/logic/admin/comment/admin_list_article_comments_logic.go"
    t = p.read_text()
    t = t.replace("req *types.PageReq", "req *types.ArticleCommentListReq")
    t = re.sub(
        r'articleID, _ := strconv\.ParseUint\("" /\* was query:article_id \*/, 10, 64\)\n\tshopID, _ := strconv\.ParseUint\("" /\* was query:shop_id \*/, 10, 64\)',
        "articleID := req.ArticleId\n\tshopID := req.ShopId",
        t,
    )
    t = t.replace('Status: "" /* was query:status */,', "Status: req.Status,")
    if "strconv." not in t:
        t = t.replace('\n\t"strconv"\n', "\n")
    p.write_text(t)
    print("logic", p.relative_to(ROOT))

    write(CAT / "internal/logic/admin/user_favorite/admin_list_user_favorites_logic.go", '''package user_favorite

import (
	"context"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListUserFavoritesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListUserFavoritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListUserFavoritesLogic {
	return &AdminListUserFavoritesLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminListUserFavoritesLogic) AdminListUserFavorites(ctx context.Context, req *types.IdPageReq) (resp *types.PageListResp, err error) {
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "用户ID无效")
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}
	list, total, err := plogic.NewFavoriteLogic(l.svcCtx).List(ctx, req.Id, req.Page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: map[string]interface{}{"list": list, "total": total}}, nil
}
''')


def rewrite_merchant_logic() -> None:
    def patch_file(rel: str, old_sig: str, new_sig: str, replacements: list[tuple[str, str]], drop_strconv=False) -> None:
        p = MER / rel
        t = p.read_text()
        t = t.replace(old_sig, new_sig)
        for a, b in replacements:
            t = t.replace(a, b)
        if drop_strconv and "strconv." not in t:
            t = t.replace('\n\t"strconv"\n', "\n")
        p.write_text(t)
        print("logic", p.relative_to(ROOT))

    patch_file(
        "internal/logic/user/coupon/list_my_coupons_logic.go",
        "req *types.PageReq",
        "req *types.StatusPageReq",
        [('"" /* was query:status */', "req.Status")],
    )
    patch_file(
        "internal/logic/merchant/coupon/merchant_list_coupons_logic.go",
        "req *types.PageReq",
        "req *types.CouponListReq",
        [('"" /* was query:status */', "req.Status"), ('"" /* was query:keyword */', "req.Keyword")],
    )
    patch_file(
        "internal/logic/admin/coupon/admin_list_coupons_logic.go",
        "req *types.PageReq",
        "req *types.CouponListReq",
        [('"" /* was query:status */', "req.Status"), ('"" /* was query:keyword */', "req.Keyword")],
    )
    patch_file(
        "internal/logic/merchant/homepage/merchant_list_slot_packages_logic.go",
        "req *types.PageReq",
        "req *types.SlotTypePageReq",
        [('"" /* was query:slot_type */', "req.SlotType")],
    )
    patch_file(
        "internal/logic/admin/homepage/admin_list_slot_packages_logic.go",
        "req *types.PageReq",
        "req *types.SlotTypePageReq",
        [('"" /* was query:slot_type */', "req.SlotType")],
    )
    patch_file(
        "internal/logic/merchant/homepage/merchant_list_slot_orders_logic.go",
        "req *types.PageReq",
        "req *types.SlotOrderListReq",
        [('"" /* was query:slot_type */', "req.SlotType"), ('"" /* was query:status */', "req.Status")],
    )
    # admin slot orders has shop_id parse
    p = MER / "internal/logic/admin/homepage/admin_list_slot_orders_logic.go"
    t = p.read_text()
    t = t.replace("req *types.PageReq", "req *types.SlotOrderListReq")
    t = t.replace('shopID, _ := strconv.ParseUint("" /* was query:shop_id */, 10, 64)', "shopID := req.ShopId")
    t = t.replace('"" /* was query:slot_type */', "req.SlotType")
    t = t.replace('"" /* was query:status */', "req.Status")
    if "strconv." not in t:
        t = t.replace('\n\t"strconv"\n', "\n")
    p.write_text(t)
    print("logic", p.relative_to(ROOT))

    for rel in [
        "internal/logic/merchant/theme/merchant_list_theme_packages_logic.go",
        "internal/logic/admin/theme/admin_list_theme_packages_logic.go",
    ]:
        p = MER / rel
        t = p.read_text()
        t = t.replace("req *types.PageReq", "req *types.ThemePackageListReq")
        t = t.replace(
            'slotID, _ := strconv.ParseUint("" /* was query:theme_slot_id */, 10, 64)',
            "slotID := req.ThemeSlotId",
        )
        if "strconv." not in t:
            t = t.replace('\n\t"strconv"\n', "\n")
        p.write_text(t)
        print("logic", p.relative_to(ROOT))

    p = MER / "internal/logic/admin/theme/admin_list_theme_orders_logic.go"
    t = p.read_text()
    t = t.replace("req *types.PageReq", "req *types.ThemeOrderListReq")
    t = t.replace('shopID, _ := strconv.ParseUint("" /* was query:shop_id */, 10, 64)', "shopID := req.ShopId")
    t = t.replace(
        'slotID, _ := strconv.ParseUint("" /* was query:theme_slot_id */, 10, 64)',
        "slotID := req.ThemeSlotId",
    )
    if "strconv." not in t:
        t = t.replace('\n\t"strconv"\n', "\n")
    p.write_text(t)
    print("logic", p.relative_to(ROOT))

    patch_file(
        "internal/logic/admin/application/admin_list_applications_logic.go",
        "req *types.PageReq",
        "req *types.StatusPageReq",
        [('"" /* was query:status */', "req.Status")],
    )
    patch_file(
        "internal/logic/admin/shop/admin_list_shops_logic.go",
        "req *types.PageReq",
        "req *types.ShopListReq",
        [('"" /* was query:status */', "req.Status"), ('"" /* was query:name */', "req.Name")],
    )
    p = MER / "internal/logic/admin/seckill/admin_list_seckill_entries_logic.go"
    t = p.read_text()
    t = t.replace("req *types.PageReq", "req *types.SeckillEntryListReq")
    t = t.replace('sid, _ := strconv.ParseUint("" /* was query:session_id */, 10, 64)', "sid := req.SessionId")
    if "strconv." not in t:
        t = t.replace('\n\t"strconv"\n', "\n")
    p.write_text(t)
    print("logic", p.relative_to(ROOT))

    # IdPageReq claims/redeems/wallet
    for rel, method in [
        ("internal/logic/merchant/coupon/merchant_coupon_claims_logic.go", "MerchantCouponClaims"),
        ("internal/logic/merchant/coupon/merchant_coupon_redeems_logic.go", "MerchantCouponRedeems"),
        ("internal/logic/admin/coupon/admin_coupon_claims_logic.go", "AdminCouponClaims"),
        ("internal/logic/admin/coupon/admin_coupon_redeems_logic.go", "AdminCouponRedeems"),
    ]:
        p = MER / rel
        t = p.read_text()
        t = t.replace("req *types.IdPathReq", "req *types.IdPageReq")
        t = t.replace("page, pageSize := 1, 10", "page, pageSize := req.Page, req.PageSize")
        p.write_text(t)
        print("logic", p.relative_to(ROOT))

    p = MER / "internal/logic/admin/wallet/admin_wallet_logs_logic.go"
    t = p.read_text()
    t = t.replace("req *types.IdPathReq", "req *types.IdPageReq")
    t = t.replace("p, ps := 1, 10", "p, ps := req.Page, req.PageSize")
    p.write_text(t)
    print("logic", p.relative_to(ROOT))


if __name__ == "__main__":
    main()
