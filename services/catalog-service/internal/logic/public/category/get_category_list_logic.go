package category

import (
	"context"
	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryListLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryListLogic {
	return &GetCategoryListLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryListLogic) GetCategoryList(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	page, pageSize := req.Page, req.PageSize
	pageReq := &pagination.PageReq{Page: page, PageSize: pageSize}
	data, err := plogic.NewCatalogLogic(l.svcCtx).GetCategoryList(ctx, pageReq)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, "查询失败")
	}
	return &types.PageListResp{List: data}, nil
}
