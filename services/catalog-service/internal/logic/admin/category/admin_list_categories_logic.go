package category

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/catalog-service/internal/product/app/admin"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListCategoriesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListCategoriesLogic(svcCtx *svc.ServiceContext) *AdminListCategoriesLogic {
	return &AdminListCategoriesLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminListCategoriesLogic) AdminListCategories(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/categories", nil, url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}, nil, hadmin.NewCatalogHandler(l.svcCtx).AdminListCategories)
	if err != nil {
		return nil, err
	}
	var out types.PageListResp
	if err := httpinvoke.Decode(raw, &out); err != nil {
		var list interface{}
		if err2 := httpinvoke.Decode(raw, &list); err2 == nil {
			return &types.PageListResp{List: list}, nil
		}
		return nil, err
	}
	return &out, nil
}
