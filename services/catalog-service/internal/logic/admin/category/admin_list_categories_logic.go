package category

import (
	"context"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListCategoriesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListCategoriesLogic {
	return &AdminListCategoriesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListCategoriesLogic) AdminListCategories(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {

	list, err := plogic.NewCatalogLogic(l.svcCtx).ListAllCategories(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list}, nil
}
