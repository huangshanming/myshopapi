package article

import (
	"context"
	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantCreateArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCreateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCreateArticleLogic {
	return &MerchantCreateArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCreateArticleLogic) MerchantCreateArticle(ctx context.Context, req *types.ArticleSaveReq) (resp *types.ArticleResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

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

	if err := requirePerm(ctx, "article:add"); err != nil {
		return nil, err
	}
	shopID, uid, _ := shopUser(ctx)
	a, err := clogic.NewArticleLogic(l.svcCtx).MerchantCreate(ctx, shopID, uid, req.ToContent())
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.ArticleResp{Data: a}, nil
}
