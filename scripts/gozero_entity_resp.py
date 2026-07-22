#!/usr/bin/env python3
"""Promote typed AnyResp.Data payloads to first-class *Resp (wire-compatible via unwrap)."""
from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
CAT = ROOT / "services/catalog-service"
ORD = ROOT / "services/order-service"
USR = ROOT / "services/user-service"


def catalog() -> None:
    biz = CAT / "internal/types/biz_types.go"
    t = biz.read_text()
    # Add first-class resp types if missing
    if "type URLResp struct" not in t:
        t = t.replace(
            "// ---- typed response payloads (used inside AnyResp.Data) ----\n\n"
            "type URLData struct {\n\tUrl string `json:\"url\"`\n}\n\n"
            "type CountData struct {\n\tCount int64 `json:\"count\"`\n}\n\n"
            "type FavoriteStatusData struct {\n\tFavorited bool `json:\"favorited\"`\n}\n\n"
            "type EngagementData struct {\n\tLiked     bool `json:\"liked\"`\n\tFavorited bool `json:\"favorited\"`\n}\n",
            """// ---- first-class response DTOs ----

type URLResp struct {
	Url string `json:"url"`
}

type CountResp struct {
	Count int64 `json:"count"`
}

type FavoriteStatusResp struct {
	Favorited bool `json:"favorited"`
}

type EngagementResp struct {
	Liked     bool `json:"liked"`
	Favorited bool `json:"favorited"`
}

// Aliases kept for any residual references
type URLData = URLResp
type CountData = CountResp
type FavoriteStatusData = FavoriteStatusResp
type EngagementData = EngagementResp
""",
        )
        biz.write_text(t)
        print("catalog biz resp types")

    types_go = CAT / "internal/types/types.go"
    tg = types_go.read_text()
    if "type URLResp struct" not in tg and "type CountResp struct" not in tg:
        tg = tg.replace(
            "type EmptyResp struct {\n}\n",
            "type EmptyResp struct {\n}\n\n"
            "type URLResp struct {\n\tUrl string `json:\"url\"`\n}\n\n"
            "type CountResp struct {\n\tCount int64 `json:\"count\"`\n}\n\n"
            "type FavoriteStatusResp struct {\n\tFavorited bool `json:\"favorited\"`\n}\n\n"
            "type EngagementResp struct {\n\tLiked bool `json:\"liked\"`\n\tFavorited bool `json:\"favorited\"`\n}\n",
        )
        types_go.write_text(tg)

    api = next((CAT / "api").glob("*.api"))
    at = api.read_text()
    if "type URLResp {" not in at:
        at = at.replace(
            "type EmptyResp {}",
            "type EmptyResp {}\n\n"
            "type URLResp {\n\tUrl string `json:\"url\"`\n}\n\n"
            "type CountResp {\n\tCount int64 `json:\"count\"`\n}\n\n"
            "type FavoriteStatusResp {\n\tFavorited bool `json:\"favorited\"`\n}\n\n"
            "type EngagementResp {\n\tLiked bool `json:\"liked\"`\n\tFavorited bool `json:\"favorited\"`\n}\n",
        )
    route_repls = [
        ("get /api/v1/products/:id/favorite-count (IdPathReq) returns (AnyResp)", "get /api/v1/products/:id/favorite-count (IdPathReq) returns (CountResp)"),
        ("get /api/v1/articles/:id/engagement (IdPathReq) returns (AnyResp)", "get /api/v1/articles/:id/engagement (IdPathReq) returns (EngagementResp)"),
        ("get /api/v1/products/:id/favorite (IdPathReq) returns (AnyResp)", "get /api/v1/products/:id/favorite (IdPathReq) returns (FavoriteStatusResp)"),
        ("post /api/v1/user/article-uploads returns (AnyResp)", "post /api/v1/user/article-uploads returns (URLResp)"),
        ("post /api/v1/merchant/uploads/images returns (AnyResp)", "post /api/v1/merchant/uploads/images returns (URLResp)"),
        ("post /api/v1/merchant/article-uploads returns (AnyResp)", "post /api/v1/merchant/article-uploads returns (URLResp)"),
        ("post /api/v1/admin/article-uploads returns (AnyResp)", "post /api/v1/admin/article-uploads returns (URLResp)"),
        ("post /api/v1/admin/banners/upload returns (AnyResp)", "post /api/v1/admin/banners/upload returns (URLResp)"),
        ("post /api/v1/admin/shop-uploads returns (AnyResp)", "post /api/v1/admin/shop-uploads returns (URLResp)"),
        ("get /api/v1/merchant/notifications/unread-count returns (AnyResp)", "get /api/v1/merchant/notifications/unread-count returns (CountResp)"),
    ]
    for a, b in route_repls:
        at = at.replace(a, b)
    api.write_text(at)

    # logic conversions
    convs = [
        (
            "internal/logic/public/product/count_logic.go",
            "*types.AnyResp",
            "*types.CountResp",
            "return &types.AnyResp{Data: types.CountData{Count: n}}, nil",
            "return &types.CountResp{Count: n}, nil",
        ),
        (
            "internal/logic/user/favorite/user_favorite_status_logic.go",
            "*types.AnyResp",
            "*types.FavoriteStatusResp",
            "return &types.AnyResp{Data: types.FavoriteStatusData{Favorited: okFav}}, nil",
            "return &types.FavoriteStatusResp{Favorited: okFav}, nil",
        ),
        (
            "internal/logic/user/article/user_article_engagement_logic.go",
            "*types.AnyResp",
            "*types.EngagementResp",
            "return &types.AnyResp{Data: types.EngagementData{Liked: liked, Favorited: favorited}}, nil",
            "return &types.EngagementResp{Liked: liked, Favorited: favorited}, nil",
        ),
    ]
    upload_files = [
        "internal/logic/user/article/upload_mine_logic.go",
        "internal/logic/merchant/product/merchant_upload_image_logic.go",
        "internal/logic/merchant/article/merchant_upload_article_logic.go",
        "internal/logic/admin/article/admin_upload_article_logic.go",
        "internal/logic/admin/banner/upload_banner_logic.go",
        "internal/logic/admin/shop/admin_upload_shop_logic.go",
    ]
    for rel, old_t, new_t, old_r, new_r in convs:
        p = CAT / rel
        text = p.read_text().replace(old_t, new_t).replace(old_r, new_r)
        # also replace CountData alias forms
        text = text.replace("types.CountData", "types.CountResp")
        p.write_text(text)
        print("cat", rel)

    for rel in upload_files:
        p = CAT / rel
        text = p.read_text()
        text = re.sub(r"\*types\.AnyResp", "*types.URLResp", text)
        text = text.replace(
            "return &types.AnyResp{Data: types.URLData{Url: url}}, nil",
            "return &types.URLResp{Url: url}, nil",
        )
        text = text.replace(
            "return &types.AnyResp{Data: types.URLResp{Url: url}}, nil",
            "return &types.URLResp{Url: url}, nil",
        )
        p.write_text(text)
        print("cat upload", rel)

    # unread count
    p = CAT / "internal/logic/merchant/notification/merchant_unread_notification_count_logic.go"
    if p.exists():
        text = p.read_text()
        # see what it returns
        if "AnyResp" in text:
            text = text.replace("*types.AnyResp", "*types.CountResp")
            # common patterns
            text = re.sub(
                r"return &types\.AnyResp\{Data: ([^}]+)\}, nil",
                r"return &types.CountResp{Count: \1}, nil",
                text,
            )
            # if data is already a number variable
            if "UnreadCount" in text and "CountResp" in text:
                # fix if Data: data where data might be int64 or map
                pass
            p.write_text(text)
            print("cat unread", p.read_text()[-200:])


def order() -> None:
    biz = ORD / "internal/types/biz_types.go"
    t = biz.read_text()
    t = t.replace("type OrderDetailData struct {", "type OrderDetailResp struct {")
    t = t.replace("type ListData struct {", "type ListResp struct {")
    if "type OrderDetailData =" not in t:
        t += "\ntype OrderDetailData = OrderDetailResp\ntype ListData = ListResp\n"
    biz.write_text(t)

    api = next((ORD / "api").glob("*.api"))
    at = api.read_text()
    if "type OrderDetailResp {" not in at:
        at = at.replace(
            "type EmptyResp {}",
            "type EmptyResp {}\n\n"
            "type OrderDetailResp {\n\tOrder      interface{} `json:\"order\"`\n\tAfterSales interface{} `json:\"after_sales\"`\n}\n\n"
            "type ListResp {\n\tList interface{} `json:\"list\"`\n}\n",
        )
    at = at.replace(
        "get /api/v1/merchant/orders/:id (IdPathReq) returns (AnyResp)",
        "get /api/v1/merchant/orders/:id (IdPathReq) returns (OrderDetailResp)",
    )
    at = at.replace(
        "get /api/v1/admin/orders/:id (IdPathReq) returns (AnyResp)",
        "get /api/v1/admin/orders/:id (IdPathReq) returns (OrderDetailResp)",
    )
    at = at.replace(
        "get /api/v1/logistics/options returns (AnyResp)",
        "get /api/v1/logistics/options returns (ListResp)",
    )
    api.write_text(at)

    for rel in [
        "internal/logic/admin/order/admin_detail_logic.go",
        "internal/logic/merchant/order/merchant_detail_logic.go",
    ]:
        p = ORD / rel
        text = p.read_text()
        text = text.replace("*types.AnyResp", "*types.OrderDetailResp")
        text = text.replace(
            "return &types.AnyResp{Data: types.OrderDetailData{Order: order, AfterSales: as}}, nil",
            "return &types.OrderDetailResp{Order: order, AfterSales: as}, nil",
        )
        text = text.replace(
            "return &types.AnyResp{Data: types.OrderDetailResp{Order: order, AfterSales: as}}, nil",
            "return &types.OrderDetailResp{Order: order, AfterSales: as}, nil",
        )
        p.write_text(text)
        print("ord detail", rel)

    p = ORD / "internal/logic/shared/logistics/logistics_options_logic.go"
    text = p.read_text()
    text = text.replace("*types.AnyResp", "*types.ListResp")
    text = text.replace(
        "return &types.AnyResp{Data: types.ListData{List: list}}, nil",
        "return &types.ListResp{List: list}, nil",
    )
    text = text.replace(
        "return &types.AnyResp{Data: types.ListResp{List: list}}, nil",
        "return &types.ListResp{List: list}, nil",
    )
    p.write_text(text)
    print("ord list options")


def user() -> None:
    # UserTokenResp add ExpireHours
    types_go = USR / "internal/types/types.go"
    t = types_go.read_text()
    if "ExpireHours" not in t:
        t = t.replace(
            """type UserTokenResp struct {
	Token    string `json:"token"`
	UserId   uint64 `json:"user_id"`
	Mobile   string `json:"mobile"`
	Nickname string `json:"nickname"`
	Role     string `json:"role,optional"`
	ShopId   uint64 `json:"shop_id,optional"`
}""",
            """type UserTokenResp struct {
	Token       string `json:"token"`
	UserId      uint64 `json:"user_id"`
	Mobile      string `json:"mobile"`
	Nickname    string `json:"nickname"`
	Role        string `json:"role,optional"`
	ShopId      uint64 `json:"shop_id,optional"`
	ExpireHours int    `json:"expire_hours,optional"`
}""",
        )
        types_go.write_text(t)

    api = next((USR / "api").glob("*.api"))
    at = api.read_text()
    if "type UserTokenResp {" not in at:
        at = at.replace(
            "type TokenResp {\n\tToken string `json:\"token\"`\n}",
            "type TokenResp {\n\tToken string `json:\"token\"`\n}\n\n"
            "type UserTokenResp {\n\tToken string `json:\"token\"`\n\tUserId uint64 `json:\"user_id\"`\n"
            "\tMobile string `json:\"mobile\"`\n\tNickname string `json:\"nickname\"`\n"
            "\tRole string `json:\"role,optional\"`\n\tShopId uint64 `json:\"shop_id,optional\"`\n"
            "\tExpireHours int `json:\"expire_hours,optional\"`\n}\n",
        )
    at = at.replace(
        "post /api/v1/admin/users/:id/token (IdPathReq) returns (AnyResp)",
        "post /api/v1/admin/users/:id/token (IdPathReq) returns (UserTokenResp)",
    )
    api.write_text(at)

    # rewrite generate_user_token logic
    p = USR / "internal/logic/admin/user/generate_user_token_logic.go"
    p.write_text('''package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type GenerateUserTokenLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGenerateUserTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateUserTokenLogic {
	return &GenerateUserTokenLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GenerateUserTokenLogic) GenerateUserToken(ctx context.Context, req *types.IdPathReq) (resp *types.UserTokenResp, err error) {
	data, err := biz.NewRBACLogic(l.svcCtx).GenerateUserToken(ctx, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	out := &types.UserTokenResp{}
	if v, ok := data["token"].(string); ok {
		out.Token = v
	}
	if v, ok := data["user_id"].(uint64); ok {
		out.UserId = v
	}
	if v, ok := data["mobile"].(string); ok {
		out.Mobile = v
	}
	if v, ok := data["nickname"].(string); ok {
		out.Nickname = v
	}
	if v, ok := data["role"].(string); ok {
		out.Role = v
	}
	if v, ok := data["shop_id"].(uint64); ok {
		out.ShopId = v
	}
	switch v := data["expire_hours"].(type) {
	case int:
		out.ExpireHours = v
	case int64:
		out.ExpireHours = int(v)
	}
	return out, nil
}
''')
    print("user token resp")


def fix_catalog_unread() -> None:
    p = CAT / "internal/logic/merchant/notification/merchant_unread_notification_count_logic.go"
    text = p.read_text()
    print("unread before signature check")
    # rewrite whole if needed
    if "CountResp" in text and "AnyResp" not in text:
        return
    # read UnreadCount return type from domain
    p.write_text('''package notification

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	nlogic "mymall/services/catalog-service/internal/notify/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantUnreadNotificationCountLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantUnreadNotificationCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantUnreadNotificationCountLogic {
	return &MerchantUnreadNotificationCountLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *MerchantUnreadNotificationCountLogic) MerchantUnreadNotificationCount(ctx context.Context) (resp *types.CountResp, err error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	n, err := nlogic.NewNotificationLogic(l.svcCtx).UnreadCount(ctx, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.CountResp{Count: n}, nil
}
''')
    print("cat unread rewritten")


def main() -> None:
    catalog()
    fix_catalog_unread()
    order()
    user()


if __name__ == "__main__":
    main()
