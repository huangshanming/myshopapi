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

type MerchantListArticleCategoriesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListArticleCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListArticleCategoriesLogic {
	return &MerchantListArticleCategoriesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantListArticleCategoriesLogic) MerchantListArticleCategories(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {

	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	if _, _, ok := shopUser(ctx); !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	tree, err := clogic.NewArticleLogic(l.svcCtx).CategoryTree(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: tree}, nil
}
