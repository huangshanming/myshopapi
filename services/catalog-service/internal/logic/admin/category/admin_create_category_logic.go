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

type AdminCreateCategoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateCategoryLogic(svcCtx *svc.ServiceContext) *AdminCreateCategoryLogic {
	return &AdminCreateCategoryLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateCategoryLogic) AdminCreateCategory(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/categories", nil, nil, req, hadmin.NewCatalogHandler(l.svcCtx).AdminCreateCategory)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
