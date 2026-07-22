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

type AdminCreateCategoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateCategoryLogic {
	return &AdminCreateCategoryLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateCategoryLogic) AdminCreateCategory(ctx context.Context, req *types.CategoryReq) (resp *types.CategoryResp, err error) {
	if req.Name == "" {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	cat, err := plogic.NewCatalogLogic(l.svcCtx).CreateCategory(ctx, req.ToProduct())
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.CategoryResp{Data: cat}, nil
}
