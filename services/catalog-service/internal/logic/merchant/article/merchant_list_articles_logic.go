package article

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"mymall/services/catalog-service/internal/content/repository"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListArticlesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListArticlesLogic {
	return &MerchantListArticlesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantListArticlesLogic) MerchantListArticles(ctx context.Context, req *types.MerchantArticleListReq) (resp *types.PageListResp, err error) {
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
	data, err := clogic.NewArticleLogic(l.svcCtx).List(ctx, repository.ArticleListFilter{
		ShopID: shopID, Title: req.Title,
		AuditStatus: req.AuditStatus,
		Status:      req.Status,
		Page:        page, PageSize: pageSize,
	})
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
