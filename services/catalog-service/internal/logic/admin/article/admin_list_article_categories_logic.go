package article

import (
	"context"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListArticleCategoriesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListArticleCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListArticleCategoriesLogic {
	return &AdminListArticleCategoriesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListArticleCategoriesLogic) AdminListArticleCategories(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {

	tree, err := clogic.NewArticleLogic(l.svcCtx).CategoryTree(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: tree}, nil
}
