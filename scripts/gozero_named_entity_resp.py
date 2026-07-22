#!/usr/bin/env python3
"""Promote DataResp routes to named *Resp across all HTTP services.

Entity payloads keep wire compatibility via `type XxxResp = DataResp` (unwrap).
Known object shapes get real structs (ListResp, AuthMeResp, …).
"""
from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

DATA_SHAPE = '{\n\tData interface{} `json:"data,optional"`\n}'


def ensure_api_types(api_text: str, type_names: list[str], concrete: dict[str, str] | None = None) -> str:
    concrete = concrete or {}
    for name, body in concrete.items():
        if f"type {name} {{" not in api_text:
            # insert after DataResp
            needle = "type DataResp {\n\tData interface{} `json:\"data,optional\"`\n}"
            if needle in api_text:
                api_text = api_text.replace(needle, needle + f"\n\ntype {name} {body}", 1)
            else:
                api_text = api_text.replace(
                    "type EmptyResp {}",
                    f"type EmptyResp {{}}\n\ntype {name} {body}",
                    1,
                )
    for name in type_names:
        if name in concrete:
            continue
        if f"type {name} {{" not in api_text:
            needle = "type DataResp {\n\tData interface{} `json:\"data,optional\"`\n}"
            if needle in api_text:
                api_text = api_text.replace(needle, needle + f"\n\ntype {name} {DATA_SHAPE}", 1)
            else:
                api_text += f"\n\ntype {name} {DATA_SHAPE}\n"
    return api_text


def write_aliases(go_path: Path, aliases: list[str], extra_structs: str = "") -> None:
    lines = [
        "package types",
        "",
        "// Named entity responses — aliases of DataResp keep unwrap MarshalJSON / legacy wire.",
        "",
    ]
    for a in sorted(set(aliases)):
        lines.append(f"type {a} = DataResp")
    lines.append("")
    if extra_structs:
        lines.append(extra_structs.rstrip())
        lines.append("")
    go_path.write_text("\n".join(lines))
    print("wrote", go_path.relative_to(ROOT))


def patch_logic(path: Path, resp_type: str) -> None:
    text = path.read_text()
    orig = text
    text = text.replace("*types.DataResp", f"*types.{resp_type}")
    text = text.replace("&types.DataResp{", f"&types.{resp_type}{{")
    if text != orig:
        path.write_text(text)
        print("logic", path.relative_to(ROOT), "->", resp_type)


def patch_api_routes(api: Path, replacements: list[tuple[str, str]]) -> None:
    at = api.read_text()
    for old, new in replacements:
        if old not in at:
            # already done or slightly different
            if new not in at:
                print("WARN missing route:", old[:80])
            continue
        at = at.replace(old, new)
    api.write_text(at)


def order() -> None:
    svc = ROOT / "services/order-service"
    aliases = [
        "OrderResp",
        "AfterSaleResp",
        "ReviewResp",
        "LogisticsCompanyResp",
    ]
    write_aliases(svc / "internal/types/entity_resp.go", aliases)

    mapping = [
        ("internal/logic/user/order/user_create_order_logic.go", "OrderResp"),
        ("internal/logic/user/order/user_get_order_logic.go", "OrderResp"),
        ("internal/logic/user/order/create_after_sale_logic.go", "AfterSaleResp"),
        ("internal/logic/user/review/user_create_review_logic.go", "ReviewResp"),
        ("internal/logic/user/review/get_by_order_logic.go", "ReviewResp"),
        ("internal/logic/admin/logistics/admin_create_logistics_logic.go", "LogisticsCompanyResp"),
    ]
    for rel, t in mapping:
        patch_logic(svc / rel, t)

    api = next((svc / "api").glob("*.api"))
    at = api.read_text()
    at = ensure_api_types(at, aliases)
    api.write_text(at)
    patch_api_routes(
        api,
        [
            ("post /api/v1/orders (CreateOrderReq) returns (DataResp)", "post /api/v1/orders (CreateOrderReq) returns (OrderResp)"),
            ("get /api/v1/orders/:id (IdPathReq) returns (DataResp)", "get /api/v1/orders/:id (IdPathReq) returns (OrderResp)"),
            ("post /api/v1/orders/:id/after-sales (CreateAfterSaleBodyReq) returns (DataResp)", "post /api/v1/orders/:id/after-sales (CreateAfterSaleBodyReq) returns (AfterSaleResp)"),
            ("post /api/v1/orders/:id/reviews (CreateReviewBodyReq) returns (DataResp)", "post /api/v1/orders/:id/reviews (CreateReviewBodyReq) returns (ReviewResp)"),
            ("get /api/v1/orders/:id/review (IdPathReq) returns (DataResp)", "get /api/v1/orders/:id/review (IdPathReq) returns (ReviewResp)"),
            ("post /api/v1/admin/logistics (LogisticsSaveBodyReq) returns (DataResp)", "post /api/v1/admin/logistics (LogisticsSaveBodyReq) returns (LogisticsCompanyResp)"),
        ],
    )
    print("order done")


def user() -> None:
    svc = ROOT / "services/user-service"
    aliases = [
        "UserResp",
        "WalletResp",
        "AddressResp",
        "MenuResp",
        "RoleResp",
        "TaskResp",
        "PointsProductResp",
        "PointsOrderResp",
        "NotificationBatchResp",
        "NotificationResp",
        "CheckinResp",
    ]
    write_aliases(svc / "internal/types/entity_resp.go", aliases)

    mapping = [
        ("internal/logic/public/auth/register_logic.go", "UserResp"),
        ("internal/logic/user/profile/user_profile_logic.go", "UserResp"),
        ("internal/logic/admin/user/get_user_logic.go", "UserResp"),
        ("internal/logic/admin/staff/create_admin_logic.go", "UserResp"),
        ("internal/logic/user/wallet/user_get_wallet_logic.go", "WalletResp"),
        ("internal/logic/admin/user/admin_get_wallet_logic.go", "WalletResp"),
        ("internal/logic/admin/user/admin_adjust_wallet_logic.go", "WalletResp"),
        ("internal/logic/user/address/user_create_address_logic.go", "AddressResp"),
        ("internal/logic/internalapi/address/internal_get_logic.go", "AddressResp"),
        ("internal/logic/admin/menu/create_menu_logic.go", "MenuResp"),
        ("internal/logic/admin/role/create_role_logic.go", "RoleResp"),
        ("internal/logic/admin/task/admin_update_task_logic.go", "TaskResp"),
        ("internal/logic/user/task/user_checkin_logic.go", "CheckinResp"),
        ("internal/logic/admin/points_mall/create_points_product_logic.go", "PointsProductResp"),
        ("internal/logic/admin/points_mall/detail_points_product_logic.go", "PointsProductResp"),
        ("internal/logic/admin/points_mall/update_points_product_logic.go", "PointsProductResp"),
        ("internal/logic/admin/points_mall/set_points_product_status_logic.go", "PointsProductResp"),
        ("internal/logic/user/points_mall/exchange_logic.go", "PointsOrderResp"),
        ("internal/logic/user/points_mall/detail_user_points_order_logic.go", "PointsOrderResp"),
        ("internal/logic/admin/points_mall/detail_points_order_logic.go", "PointsOrderResp"),
        ("internal/logic/admin/points_mall/ship_points_order_logic.go", "PointsOrderResp"),
        ("internal/logic/admin/points_mall/complete_points_order_logic.go", "PointsOrderResp"),
        ("internal/logic/admin/points_mall/cancel_points_order_logic.go", "PointsOrderResp"),
        ("internal/logic/admin/points_mall/remark_points_order_logic.go", "PointsOrderResp"),
        ("internal/logic/admin/notification/admin_send_notification_logic.go", "NotificationBatchResp"),
        ("internal/logic/internalapi/notification/internal_create_notification_logic.go", "NotificationResp"),
    ]
    for rel, t in mapping:
        patch_logic(svc / rel, t)

    api = next((svc / "api").glob("*.api"))
    at = api.read_text()
    at = ensure_api_types(at, aliases)
    api.write_text(at)
    patch_api_routes(
        api,
        [
            ("post /api/v1/user/register (RegisterReq) returns (DataResp)", "post /api/v1/user/register (RegisterReq) returns (UserResp)"),
            ("get /api/v1/user/addresses/internal (InternalAddressReq) returns (DataResp)", "get /api/v1/user/addresses/internal (InternalAddressReq) returns (AddressResp)"),
            ("post /api/v1/internal/notifications (NotifyCreateReq) returns (DataResp)", "post /api/v1/internal/notifications (NotifyCreateReq) returns (NotificationResp)"),
            ("get /api/v1/user/profile returns (DataResp)", "get /api/v1/user/profile returns (UserResp)"),
            ("get /api/v1/user/wallet returns (DataResp)", "get /api/v1/user/wallet returns (WalletResp)"),
            ("post /api/v1/user/addresses (AddressReq) returns (DataResp)", "post /api/v1/user/addresses (AddressReq) returns (AddressResp)"),
            ("post /api/v1/user/tasks/checkin returns (DataResp)", "post /api/v1/user/tasks/checkin returns (CheckinResp)"),
            ("post /api/v1/admin/menus (MenuReq) returns (DataResp)", "post /api/v1/admin/menus (MenuReq) returns (MenuResp)"),
            ("post /api/v1/admin/roles (RoleReq) returns (DataResp)", "post /api/v1/admin/roles (RoleReq) returns (RoleResp)"),
            ("get /api/v1/admin/users/:id (IdPathReq) returns (DataResp)", "get /api/v1/admin/users/:id (IdPathReq) returns (UserResp)"),
            ("get /api/v1/admin/users/:id/wallet (IdPathReq) returns (DataResp)", "get /api/v1/admin/users/:id/wallet (IdPathReq) returns (WalletResp)"),
            ("post /api/v1/admin/users/:id/wallet/adjust (WalletAdjustReq) returns (DataResp)", "post /api/v1/admin/users/:id/wallet/adjust (WalletAdjustReq) returns (WalletResp)"),
            ("post /api/v1/admin/admins (AdminCreateReq) returns (DataResp)", "post /api/v1/admin/admins (AdminCreateReq) returns (UserResp)"),
            ("post /api/v1/admin/notifications/send (AdminSendReq) returns (DataResp)", "post /api/v1/admin/notifications/send (AdminSendReq) returns (NotificationBatchResp)"),
            ("put /api/v1/admin/tasks/:id (UpdateTaskReq) returns (DataResp)", "put /api/v1/admin/tasks/:id (UpdateTaskReq) returns (TaskResp)"),
            ("post /api/v1/user/points-mall/exchange (ExchangeReq) returns (DataResp)", "post /api/v1/user/points-mall/exchange (ExchangeReq) returns (PointsOrderResp)"),
            ("get /api/v1/user/points-mall/orders/:id (IdPathReq) returns (DataResp)", "get /api/v1/user/points-mall/orders/:id (IdPathReq) returns (PointsOrderResp)"),
            ("post /api/v1/admin/points-products (PointsProductSaveReq) returns (DataResp)", "post /api/v1/admin/points-products (PointsProductSaveReq) returns (PointsProductResp)"),
            ("get /api/v1/admin/points-products/:id (IdPathReq) returns (DataResp)", "get /api/v1/admin/points-products/:id (IdPathReq) returns (PointsProductResp)"),
            ("put /api/v1/admin/points-products/:id (PointsProductUpdateReq) returns (DataResp)", "put /api/v1/admin/points-products/:id (PointsProductUpdateReq) returns (PointsProductResp)"),
            ("put /api/v1/admin/points-products/:id/status (PointsProductStatusReq) returns (DataResp)", "put /api/v1/admin/points-products/:id/status (PointsProductStatusReq) returns (PointsProductResp)"),
            ("get /api/v1/admin/points-orders/:id (IdPathReq) returns (DataResp)", "get /api/v1/admin/points-orders/:id (IdPathReq) returns (PointsOrderResp)"),
            ("post /api/v1/admin/points-orders/:id/ship (ShipReq) returns (DataResp)", "post /api/v1/admin/points-orders/:id/ship (ShipReq) returns (PointsOrderResp)"),
            ("post /api/v1/admin/points-orders/:id/complete (IdPathReq) returns (DataResp)", "post /api/v1/admin/points-orders/:id/complete (IdPathReq) returns (PointsOrderResp)"),
            ("post /api/v1/admin/points-orders/:id/cancel (RemarkReq) returns (DataResp)", "post /api/v1/admin/points-orders/:id/cancel (RemarkReq) returns (PointsOrderResp)"),
            ("put /api/v1/admin/points-orders/:id/remark (RemarkReq) returns (DataResp)", "put /api/v1/admin/points-orders/:id/remark (RemarkReq) returns (PointsOrderResp)"),
        ],
    )
    print("user done")


def catalog() -> None:
    svc = ROOT / "services/catalog-service"
    aliases = [
        "ProductResp",
        "CategoryResp",
        "ArticleResp",
        "CommentResp",
        "ProductJobResp",
        "ImportResultResp",
        "StockWarningsResp",
        "TagResp",
        "AttrTemplateResp",
        "ShopRoleResp",
        "BannerResp",
        "EmojiResp",
        "ArticleStatsResp",
        "SalesRankResp",
        "RoleMenuIdsResp",
    ]
    concrete = {
        "ShopAuthMeResp": "{\n\tPerms interface{} `json:\"perms\"`\n\tMenus interface{} `json:\"menus\"`\n"
        "\tMenuTree interface{} `json:\"menu_tree\"`\n\tIsOwner bool `json:\"is_owner\"`\n}",
        "BindStaffResp": "{\n\tUserId uint64 `json:\"user_id\"`\n\tMsg string `json:\"msg\"`\n}",
    }
    extra = """
type ShopAuthMeResp struct {
	Perms    interface{} `json:"perms"`
	Menus    interface{} `json:"menus"`
	MenuTree interface{} `json:"menu_tree"`
	IsOwner  bool        `json:"is_owner"`
}

type BindStaffResp struct {
	UserId uint64 `json:"user_id"`
	Msg    string `json:"msg"`
}
"""
    write_aliases(svc / "internal/types/entity_resp.go", aliases, extra)

    mapping = [
        ("internal/logic/public/product/get_product_detail_logic.go", "ProductResp"),
        ("internal/logic/public/product/get_sales_rank_logic.go", "SalesRankResp"),
        ("internal/logic/public/category/get_category_detail_logic.go", "CategoryResp"),
        ("internal/logic/public/article/public_get_article_logic.go", "ArticleResp"),
        ("internal/logic/user/article/create_comment_logic.go", "CommentResp"),
        ("internal/logic/user/article/create_mine_logic.go", "ArticleResp"),
        ("internal/logic/user/article/detail_mine_logic.go", "ArticleResp"),
        ("internal/logic/merchant/product/merchant_create_product_logic.go", "ProductResp"),
        ("internal/logic/merchant/product/merchant_batch_products_logic.go", "ProductJobResp"),
        ("internal/logic/merchant/product/job_status_logic.go", "ProductJobResp"),
        ("internal/logic/merchant/product/merchant_import_products_logic.go", "ImportResultResp"),
        ("internal/logic/merchant/product/merchant_get_product_logic.go", "ProductResp"),
        ("internal/logic/merchant/product/merchant_update_product_logic.go", "ProductResp"),
        ("internal/logic/merchant/product/merchant_copy_product_logic.go", "ProductResp"),
        ("internal/logic/merchant/product/stock_warnings_logic.go", "StockWarningsResp"),
        ("internal/logic/merchant/product/merchant_create_tag_logic.go", "TagResp"),
        ("internal/logic/merchant/product/merchant_update_tag_logic.go", "TagResp"),
        ("internal/logic/merchant/product/merchant_create_attr_template_logic.go", "AttrTemplateResp"),
        ("internal/logic/merchant/product/merchant_update_attr_template_logic.go", "AttrTemplateResp"),
        ("internal/logic/merchant/shopops/merchant_create_role_logic.go", "ShopRoleResp"),
        ("internal/logic/merchant/shopops/merchant_update_role_logic.go", "ShopRoleResp"),
        ("internal/logic/merchant/shopops/role_menus_logic.go", "RoleMenuIdsResp"),
        ("internal/logic/merchant/article/merchant_create_article_logic.go", "ArticleResp"),
        ("internal/logic/merchant/article/merchant_get_article_logic.go", "ArticleResp"),
        ("internal/logic/admin/category/admin_create_category_logic.go", "CategoryResp"),
        ("internal/logic/admin/article/admin_article_stats_logic.go", "ArticleStatsResp"),
        ("internal/logic/admin/article/admin_create_article_logic.go", "ArticleResp"),
        ("internal/logic/admin/article/admin_get_article_logic.go", "ArticleResp"),
        ("internal/logic/admin/comment/emoji_create_logic.go", "EmojiResp"),
        ("internal/logic/admin/banner/create_banner_logic.go", "BannerResp"),
        ("internal/logic/admin/banner/get_banner_logic.go", "BannerResp"),
    ]
    for rel, t in mapping:
        patch_logic(svc / rel, t)

    # special concrete conversions
    auth = svc / "internal/logic/merchant/shopops/merchant_auth_me_logic.go"
    text = auth.read_text()
    text = text.replace("*types.DataResp", "*types.ShopAuthMeResp")
    text = re.sub(
        r"return &types\.DataResp\{Data: map\[string\]interface\{\}\{\n"
        r'\t\t"perms": perms, "menus": menus, "menu_tree": repository\.BuildShopMenuTree\(menus\),\n'
        r'\t\t"is_owner": l\.svcCtx\.ShopRBAC\.IsOwner\(ctx, shopID, uid\),\n'
        r"\t\}\}, nil",
        "return &types.ShopAuthMeResp{\n"
        "\t\tPerms: perms, Menus: menus, MenuTree: repository.BuildShopMenuTree(menus),\n"
        "\t\tIsOwner: l.svcCtx.ShopRBAC.IsOwner(ctx, shopID, uid),\n"
        "\t}, nil",
        text,
    )
    # fallback if already DataResp replaced
    if "ShopAuthMeResp{" not in text and "map[string]interface{}" in text:
        text = re.sub(
            r"return &types\.ShopAuthMeResp\{Data: map\[string\]interface\{\}\{\n"
            r'\t\t"perms": perms, "menus": menus, "menu_tree": repository\.BuildShopMenuTree\(menus\),\n'
            r'\t\t"is_owner": l\.svcCtx\.ShopRBAC\.IsOwner\(ctx, shopID, uid\),\n'
            r"\t\}\}, nil",
            "return &types.ShopAuthMeResp{\n"
            "\t\tPerms: perms, Menus: menus, MenuTree: repository.BuildShopMenuTree(menus),\n"
            "\t\tIsOwner: l.svcCtx.ShopRBAC.IsOwner(ctx, shopID, uid),\n"
            "\t}, nil",
            text,
        )
    auth.write_text(text)
    print("logic auth me")

    bind = svc / "internal/logic/merchant/shopops/bind_staff_logic.go"
    bt = bind.read_text()
    bt = bt.replace("*types.DataResp", "*types.BindStaffResp")
    bt = re.sub(
        r'return &types\.(DataResp|BindStaffResp)\{Data: map\[string\]interface\{\}\{"user_id": userID, "msg": msg\}\}, nil',
        "return &types.BindStaffResp{UserId: userID, Msg: msg}, nil",
        bt,
    )
    bind.write_text(bt)
    print("logic bind staff")

    api = next((svc / "api").glob("*.api"))
    at = api.read_text()
    at = ensure_api_types(at, aliases, concrete)
    api.write_text(at)
    patch_api_routes(
        api,
        [
            ("get /api/v1/products/detail (IdQueryReq) returns (DataResp)", "get /api/v1/products/detail (IdQueryReq) returns (ProductResp)"),
            ("get /api/v1/products/sales-rank (PageReq) returns (DataResp)", "get /api/v1/products/sales-rank (PageReq) returns (SalesRankResp)"),
            ("get /api/v1/product_category/detail (IdQueryReq) returns (DataResp)", "get /api/v1/product_category/detail (IdQueryReq) returns (CategoryResp)"),
            ("get /api/v1/articles/:id (IdPathReq) returns (DataResp)", "get /api/v1/articles/:id (IdPathReq) returns (ArticleResp)"),
            ("post /api/v1/articles/:id/comments (CreateCommentBodyReq) returns (DataResp)", "post /api/v1/articles/:id/comments (CreateCommentBodyReq) returns (CommentResp)"),
            ("post /api/v1/user/articles (UserArticleCreateReq) returns (DataResp)", "post /api/v1/user/articles (UserArticleCreateReq) returns (ArticleResp)"),
            ("get /api/v1/user/articles/:id (IdPathReq) returns (DataResp)", "get /api/v1/user/articles/:id (IdPathReq) returns (ArticleResp)"),
            ("post /api/v1/merchant/products (MerchantProductSaveReq) returns (DataResp)", "post /api/v1/merchant/products (MerchantProductSaveReq) returns (ProductResp)"),
            ("post /api/v1/merchant/products/batch (BatchProductReq) returns (DataResp)", "post /api/v1/merchant/products/batch (BatchProductReq) returns (ProductJobResp)"),
            ("get /api/v1/merchant/products/jobs/:id (IdPathReq) returns (DataResp)", "get /api/v1/merchant/products/jobs/:id (IdPathReq) returns (ProductJobResp)"),
            ("post /api/v1/merchant/products/import returns (DataResp)", "post /api/v1/merchant/products/import returns (ImportResultResp)"),
            ("get /api/v1/merchant/products/:id (IdPathReq) returns (DataResp)", "get /api/v1/merchant/products/:id (IdPathReq) returns (ProductResp)"),
            ("put /api/v1/merchant/products/:id (ProductUpdateBodyReq) returns (DataResp)", "put /api/v1/merchant/products/:id (ProductUpdateBodyReq) returns (ProductResp)"),
            ("post /api/v1/merchant/products/:id/copy (IdPathReq) returns (DataResp)", "post /api/v1/merchant/products/:id/copy (IdPathReq) returns (ProductResp)"),
            ("get /api/v1/merchant/stocks/warnings (PageReq) returns (DataResp)", "get /api/v1/merchant/stocks/warnings (PageReq) returns (StockWarningsResp)"),
            ("post /api/v1/merchant/tags (TagReq) returns (DataResp)", "post /api/v1/merchant/tags (TagReq) returns (TagResp)"),
            ("put /api/v1/merchant/tags/:id (TagUpdateBodyReq) returns (DataResp)", "put /api/v1/merchant/tags/:id (TagUpdateBodyReq) returns (TagResp)"),
            ("post /api/v1/merchant/attr-templates (AttrTemplateReq) returns (DataResp)", "post /api/v1/merchant/attr-templates (AttrTemplateReq) returns (AttrTemplateResp)"),
            ("put /api/v1/merchant/attr-templates/:id (AttrTemplateUpdateBodyReq) returns (DataResp)", "put /api/v1/merchant/attr-templates/:id (AttrTemplateUpdateBodyReq) returns (AttrTemplateResp)"),
            ("get /api/v1/merchant/auth/me returns (DataResp)", "get /api/v1/merchant/auth/me returns (ShopAuthMeResp)"),
            ("get /api/v1/merchant/shop/roles/:id/menus (IdPathReq) returns (DataResp)", "get /api/v1/merchant/shop/roles/:id/menus (IdPathReq) returns (RoleMenuIdsResp)"),
            ("post /api/v1/merchant/shop/roles (ShopRoleReq) returns (DataResp)", "post /api/v1/merchant/shop/roles (ShopRoleReq) returns (ShopRoleResp)"),
            ("put /api/v1/merchant/shop/roles/:id (ShopRoleUpdateBodyReq) returns (DataResp)", "put /api/v1/merchant/shop/roles/:id (ShopRoleUpdateBodyReq) returns (ShopRoleResp)"),
            ("post /api/v1/merchant/shop/staff (ShopStaffReq) returns (DataResp)", "post /api/v1/merchant/shop/staff (ShopStaffReq) returns (BindStaffResp)"),
            ("post /api/v1/merchant/articles (ArticleSaveReq) returns (DataResp)", "post /api/v1/merchant/articles (ArticleSaveReq) returns (ArticleResp)"),
            ("get /api/v1/merchant/articles/:id (IdPathReq) returns (DataResp)", "get /api/v1/merchant/articles/:id (IdPathReq) returns (ArticleResp)"),
            ("post /api/v1/admin/categories (CategoryReq) returns (DataResp)", "post /api/v1/admin/categories (CategoryReq) returns (CategoryResp)"),
            ("get /api/v1/admin/articles/stats returns (DataResp)", "get /api/v1/admin/articles/stats returns (ArticleStatsResp)"),
            ("post /api/v1/admin/articles (ArticleSaveReq) returns (DataResp)", "post /api/v1/admin/articles (ArticleSaveReq) returns (ArticleResp)"),
            ("get /api/v1/admin/articles/:id (IdPathReq) returns (DataResp)", "get /api/v1/admin/articles/:id (IdPathReq) returns (ArticleResp)"),
            ("post /api/v1/admin/comment-emojis (EmojiSaveReq) returns (DataResp)", "post /api/v1/admin/comment-emojis (EmojiSaveReq) returns (EmojiResp)"),
            ("post /api/v1/admin/banners (BannerSaveReq) returns (DataResp)", "post /api/v1/admin/banners (BannerSaveReq) returns (BannerResp)"),
            ("get /api/v1/admin/banners/:id (IdPathReq) returns (DataResp)", "get /api/v1/admin/banners/:id (IdPathReq) returns (BannerResp)"),
        ],
    )
    print("catalog done")


def merchant() -> None:
    svc = ROOT / "services/merchant-service"
    aliases = [
        "ShopResp",
        "ShopListResp",
        "HomeSlotsResp",
        "SeckillCurrentResp",
        "SeckillEntryResp",
        "SeckillSessionsResp",
        "MatchCouponsResp",
        "SeckillConsumeResp",
        "UserCouponResp",
        "ShopApplicationResp",
        "WalletResp",
        "HomepageOrderResp",
        "ThemeOrderResp",
        "CouponResp",
        "CouponGrantResp",
        "CouponStatsResp",
        "SeckillRuleResp",
        "SlotPackageResp",
        "ThemePackageResp",
    ]
    concrete = {
        "GrantedCountResp": "{\n\tGranted int64 `json:\"granted\"`\n}",
    }
    # ensure ListResp in Go
    biz = svc / "internal/types/biz_types.go"
    bt = biz.read_text()
    if "type ListResp struct" not in bt and "type ListResp =" not in bt:
        if "type ListData struct" in bt:
            bt = bt.replace(
                "type ListData struct {\n\tList interface{} `json:\"list\"`\n}",
                "type ListResp struct {\n\tList interface{} `json:\"list\"`\n}\n\ntype ListData = ListResp",
            )
        else:
            bt += "\n\ntype ListResp struct {\n\tList interface{} `json:\"list\"`\n}\n"
        biz.write_text(bt)

    types_go = svc / "internal/types/types.go"
    tg = types_go.read_text()
    if "type ListResp struct" not in tg:
        tg = tg.replace(
            "type EmptyResp struct {\n}\n",
            "type EmptyResp struct {\n}\n\ntype ListResp struct {\n\tList interface{} `json:\"list\"`\n}\n",
        )
        types_go.write_text(tg)

    extra = """
type GrantedCountResp struct {
	Granted int64 `json:"granted"`
}
"""
    write_aliases(svc / "internal/types/entity_resp.go", aliases, extra)

    mapping = [
        ("internal/logic/public/shop/public_home_slots_logic.go", "HomeSlotsResp"),
        ("internal/logic/public/shop/public_get_shop_logic.go", "ShopResp"),
        ("internal/logic/public/seckill/public_seckill_current_logic.go", "SeckillCurrentResp"),
        ("internal/logic/public/seckill/public_seckill_entry_logic.go", "SeckillEntryResp"),
        ("internal/logic/internalapi/seckill/seckill_consume_logic.go", "SeckillConsumeResp"),
        ("internal/logic/internalapi/coupon/internal_match_coupons_logic.go", "MatchCouponsResp"),
        ("internal/logic/user/coupon/claim_coupon_logic.go", "UserCouponResp"),
        ("internal/logic/merchant/shop/apply_logic.go", "ShopApplicationResp"),
        ("internal/logic/merchant/shop/my_shops_logic.go", "ShopListResp"),
        ("internal/logic/merchant/wallet/merchant_get_wallet_logic.go", "WalletResp"),
        ("internal/logic/merchant/seckill/merchant_seckill_sessions_logic.go", "SeckillSessionsResp"),
        ("internal/logic/merchant/seckill/merchant_apply_seckill_logic.go", "SeckillEntryResp"),
        ("internal/logic/merchant/seckill/merchant_set_seckill_auto_renew_logic.go", "SeckillEntryResp"),
        ("internal/logic/merchant/homepage/merchant_buy_slot_logic.go", "HomepageOrderResp"),
        ("internal/logic/merchant/theme/merchant_buy_theme_logic.go", "ThemeOrderResp"),
        ("internal/logic/merchant/coupon/merchant_create_coupon_logic.go", "CouponResp"),
        ("internal/logic/merchant/coupon/merchant_copy_coupon_logic.go", "CouponResp"),
        ("internal/logic/merchant/coupon/merchant_grant_coupon_logic.go", "CouponGrantResp"),
        ("internal/logic/merchant/coupon/merchant_coupon_stats_logic.go", "CouponStatsResp"),
        ("internal/logic/admin/application/admin_approve_logic.go", "ShopResp"),
        ("internal/logic/admin/shop/admin_create_shop_logic.go", "ShopResp"),
        ("internal/logic/admin/shop/admin_get_shop_logic.go", "ShopResp"),
        ("internal/logic/admin/wallet/admin_get_wallet_logic.go", "WalletResp"),
        ("internal/logic/admin/wallet/admin_adjust_wallet_logic.go", "WalletResp"),
        ("internal/logic/admin/seckill/admin_get_seckill_rule_logic.go", "SeckillRuleResp"),
        ("internal/logic/admin/seckill/admin_update_seckill_rule_logic.go", "SeckillRuleResp"),
        ("internal/logic/admin/homepage/admin_create_slot_package_logic.go", "SlotPackageResp"),
        ("internal/logic/admin/homepage/admin_grant_slot_logic.go", "HomepageOrderResp"),
        ("internal/logic/admin/theme/admin_create_theme_package_logic.go", "ThemePackageResp"),
        ("internal/logic/admin/theme/admin_grant_theme_logic.go", "ThemeOrderResp"),
        ("internal/logic/admin/coupon/admin_create_coupon_logic.go", "CouponResp"),
        ("internal/logic/admin/coupon/admin_copy_coupon_logic.go", "CouponResp"),
        ("internal/logic/admin/coupon/admin_grant_coupon_logic.go", "CouponGrantResp"),
        ("internal/logic/admin/coupon/admin_coupon_stats_logic.go", "CouponStatsResp"),
    ]
    for rel, t in mapping:
        patch_logic(svc / rel, t)

    # ListResp conversions
    for rel, old_ret in [
        (
            "internal/logic/public/coupon/public_coupon_center_logic.go",
            None,
        ),
        (
            "internal/logic/public/coupon/public_coupon_popup_logic.go",
            None,
        ),
        (
            "internal/logic/public/shop/public_theme_tiles_logic.go",
            None,
        ),
    ]:
        p = svc / rel
        text = p.read_text()
        text = text.replace("*types.DataResp", "*types.ListResp")
        text = text.replace(
            "return &types.DataResp{Data: types.ListData{List: list}}, nil",
            "return &types.ListResp{List: list}, nil",
        )
        text = text.replace(
            "return &types.ListResp{Data: types.ListData{List: list}}, nil",
            "return &types.ListResp{List: list}, nil",
        )
        text = text.replace(
            'return &types.DataResp{Data: map[string]interface{}{"list": list}}, nil',
            "return &types.ListResp{List: list}, nil",
        )
        text = text.replace(
            'return &types.ListResp{Data: map[string]interface{}{"list": list}}, nil',
            "return &types.ListResp{List: list}, nil",
        )
        p.write_text(text)
        print("list", rel)

    gift = svc / "internal/logic/internalapi/coupon/internal_order_gift_logic.go"
    gt = gift.read_text()
    gt = gt.replace("*types.DataResp", "*types.GrantedCountResp")
    gt = re.sub(
        r'return &types\.(DataResp|GrantedCountResp)\{Data: map\[string\]interface\{\}\{"granted": n\}\}, nil',
        "return &types.GrantedCountResp{Granted: int64(n)}, nil",
        gt,
    )
    # n might already be int64
    if "GrantedCountResp{Granted:" not in gt:
        gt = re.sub(
            r'return &types\.GrantedCountResp\{Data: map\[string\]interface\{\}\{"granted": n\}\}, nil',
            "return &types.GrantedCountResp{Granted: int64(n)}, nil",
            gt,
        )
    gift.write_text(gt)
    print("gift granted")

    api = next((svc / "api").glob("*.api"))
    at = api.read_text()
    if "type ListResp {" not in at:
        at = ensure_api_types(at, [], {"ListResp": "{\n\tList interface{} `json:\"list\"`\n}"})
    at = ensure_api_types(at, aliases, concrete)
    api.write_text(at)
    patch_api_routes(
        api,
        [
            ("get /api/v1/shops/home-slots (SlotTypeQueryReq) returns (DataResp)", "get /api/v1/shops/home-slots (SlotTypeQueryReq) returns (HomeSlotsResp)"),
            ("get /api/v1/home/theme-tiles returns (DataResp)", "get /api/v1/home/theme-tiles returns (ListResp)"),
            ("get /api/v1/shops/:id (IdPathReq) returns (DataResp)", "get /api/v1/shops/:id (IdPathReq) returns (ShopResp)"),
            ("get /api/v1/seckill/current returns (DataResp)", "get /api/v1/seckill/current returns (SeckillCurrentResp)"),
            ("get /api/v1/seckill/entries/:id (IdPathReq) returns (DataResp)", "get /api/v1/seckill/entries/:id (IdPathReq) returns (SeckillEntryResp)"),
            ("get /api/v1/coupons/center (ShopIdQueryReq) returns (DataResp)", "get /api/v1/coupons/center (ShopIdQueryReq) returns (ListResp)"),
            ("get /api/v1/coupons/popup returns (DataResp)", "get /api/v1/coupons/popup returns (ListResp)"),
            ("post /api/v1/seckill/consume (SeckillConsumeReq) returns (DataResp)", "post /api/v1/seckill/consume (SeckillConsumeReq) returns (SeckillConsumeResp)"),
            ("post /api/v1/internal/coupons/match (MatchCouponsReq) returns (DataResp)", "post /api/v1/internal/coupons/match (MatchCouponsReq) returns (MatchCouponsResp)"),
            ("post /api/v1/internal/coupons/order-gift (OrderGiftCouponReq) returns (DataResp)", "post /api/v1/internal/coupons/order-gift (OrderGiftCouponReq) returns (GrantedCountResp)"),
            ("post /api/v1/coupons/:id/claim (ClaimCouponBodyReq) returns (DataResp)", "post /api/v1/coupons/:id/claim (ClaimCouponBodyReq) returns (UserCouponResp)"),
            ("post /api/v1/merchant/apply (ApplyReq) returns (DataResp)", "post /api/v1/merchant/apply (ApplyReq) returns (ShopApplicationResp)"),
            ("get /api/v1/merchant/shops returns (DataResp)", "get /api/v1/merchant/shops returns (ShopListResp)"),
            ("get /api/v1/merchant/wallet returns (DataResp)", "get /api/v1/merchant/wallet returns (WalletResp)"),
            ("get /api/v1/merchant/seckill/sessions returns (DataResp)", "get /api/v1/merchant/seckill/sessions returns (SeckillSessionsResp)"),
            ("post /api/v1/merchant/seckill/entries (SeckillApplyReq) returns (DataResp)", "post /api/v1/merchant/seckill/entries (SeckillApplyReq) returns (SeckillEntryResp)"),
            ("put /api/v1/merchant/seckill/entries/:id/auto-renew (SeckillAutoRenewBodyReq) returns (DataResp)", "put /api/v1/merchant/seckill/entries/:id/auto-renew (SeckillAutoRenewBodyReq) returns (SeckillEntryResp)"),
            ("post /api/v1/merchant/homepage-orders (BuySlotReq) returns (DataResp)", "post /api/v1/merchant/homepage-orders (BuySlotReq) returns (HomepageOrderResp)"),
            ("post /api/v1/merchant/theme-orders (ThemeBuyReq) returns (DataResp)", "post /api/v1/merchant/theme-orders (ThemeBuyReq) returns (ThemeOrderResp)"),
            ("post /api/v1/merchant/coupons (CouponSaveReq) returns (DataResp)", "post /api/v1/merchant/coupons (CouponSaveReq) returns (CouponResp)"),
            ("post /api/v1/merchant/coupons/:id/copy (IdPathReq) returns (DataResp)", "post /api/v1/merchant/coupons/:id/copy (IdPathReq) returns (CouponResp)"),
            ("post /api/v1/merchant/coupons/grant (GrantCouponReq) returns (DataResp)", "post /api/v1/merchant/coupons/grant (GrantCouponReq) returns (CouponGrantResp)"),
            ("get /api/v1/merchant/coupons/:id/stats (IdPathReq) returns (DataResp)", "get /api/v1/merchant/coupons/:id/stats (IdPathReq) returns (CouponStatsResp)"),
            ("post /api/v1/admin/applications/:id/approve (IdPathReq) returns (DataResp)", "post /api/v1/admin/applications/:id/approve (IdPathReq) returns (ShopResp)"),
            ("post /api/v1/admin/shops (AdminCreateShopReq) returns (DataResp)", "post /api/v1/admin/shops (AdminCreateShopReq) returns (ShopResp)"),
            ("get /api/v1/admin/shops/:id (IdPathReq) returns (DataResp)", "get /api/v1/admin/shops/:id (IdPathReq) returns (ShopResp)"),
            ("get /api/v1/admin/shops/:id/wallet (IdPathReq) returns (DataResp)", "get /api/v1/admin/shops/:id/wallet (IdPathReq) returns (WalletResp)"),
            ("post /api/v1/admin/shops/:id/wallet/adjust (WalletAdjustBodyReq) returns (DataResp)", "post /api/v1/admin/shops/:id/wallet/adjust (WalletAdjustBodyReq) returns (WalletResp)"),
            ("get /api/v1/admin/seckill/rule returns (DataResp)", "get /api/v1/admin/seckill/rule returns (SeckillRuleResp)"),
            ("put /api/v1/admin/seckill/rule (SeckillRuleReq) returns (DataResp)", "put /api/v1/admin/seckill/rule (SeckillRuleReq) returns (SeckillRuleResp)"),
            ("post /api/v1/admin/homepage-packages (SlotPackageSaveReq) returns (DataResp)", "post /api/v1/admin/homepage-packages (SlotPackageSaveReq) returns (SlotPackageResp)"),
            ("post /api/v1/admin/homepage-orders/grant (GrantSlotReq) returns (DataResp)", "post /api/v1/admin/homepage-orders/grant (GrantSlotReq) returns (HomepageOrderResp)"),
            ("post /api/v1/admin/theme-packages (ThemePackageSaveReq) returns (DataResp)", "post /api/v1/admin/theme-packages (ThemePackageSaveReq) returns (ThemePackageResp)"),
            ("post /api/v1/admin/theme-orders/grant (ThemeGrantReq) returns (DataResp)", "post /api/v1/admin/theme-orders/grant (ThemeGrantReq) returns (ThemeOrderResp)"),
            ("post /api/v1/admin/coupons (CouponSaveReq) returns (DataResp)", "post /api/v1/admin/coupons (CouponSaveReq) returns (CouponResp)"),
            ("post /api/v1/admin/coupons/:id/copy (IdPathReq) returns (DataResp)", "post /api/v1/admin/coupons/:id/copy (IdPathReq) returns (CouponResp)"),
            ("post /api/v1/admin/coupons/grant (GrantCouponReq) returns (DataResp)", "post /api/v1/admin/coupons/grant (GrantCouponReq) returns (CouponGrantResp)"),
            ("get /api/v1/admin/coupons/:id/stats (IdPathReq) returns (DataResp)", "get /api/v1/admin/coupons/:id/stats (IdPathReq) returns (CouponStatsResp)"),
        ],
    )
    print("merchant done")


def main() -> None:
    order()
    user()
    catalog()
    merchant()


if __name__ == "__main__":
    main()
