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

type AdminDeleteCategoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteCategoryLogic {
	return &AdminDeleteCategoryLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteCategoryLogic) AdminDeleteCategory(ctx context.Context, req *types.IdPathReq) (resp *types.EmptyResp, err error) {
	id := req.Id
	if err := plogic.NewCatalogLogic(l.svcCtx).DeleteCategory(ctx, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
