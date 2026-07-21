package article

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"mymall/services/catalog-service/internal/content/repository"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListArticleCommentsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListArticleCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListArticleCommentsLogic {
	return &MerchantListArticleCommentsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantListArticleCommentsLogic) MerchantListArticleComments(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	page, pageSize := req.Page, req.PageSize
	articleID, _ := strconv.ParseUint("" /* was query:article_id */, 10, 64)
	data, err := clogic.NewArticleLogic(l.svcCtx).ListComments(ctx, repository.CommentListFilter{
		ShopID: shopID, ArticleID: articleID, Status: "" /* was query:status */,
		Page: page, PageSize: pageSize,
	})
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
