package article

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantGetArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantGetArticleLogic {
	return &MerchantGetArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantGetArticleLogic) MerchantGetArticle(ctx context.Context, req *types.IdPathReq) (resp *types.ArticleResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id := req.Id
	data, err := clogic.NewArticleLogic(l.svcCtx).Detail(ctx, id, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return &types.ArticleResp{Data: data}, nil
}
