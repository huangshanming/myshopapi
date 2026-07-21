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

type AdminUpdateCategoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateCategoryLogic {
	return &AdminUpdateCategoryLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateCategoryLogic) AdminUpdateCategory(ctx context.Context, req *types.CategoryUpdateBodyReq) (resp *types.AnyResp, err error) {
	id := req.Id
	if req.Name == "" {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := plogic.NewCatalogLogic(l.svcCtx).UpdateCategory(ctx, id, req.ToProduct()); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
