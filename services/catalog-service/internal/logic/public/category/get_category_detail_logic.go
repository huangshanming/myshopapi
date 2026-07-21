package category

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hpublic "mymall/services/catalog-service/internal/product/app/public"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryDetailLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetCategoryDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryDetailLogic {
	return &GetCategoryDetailLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryDetailLogic) GetCategoryDetail(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hpublic.NewCatalogHandler(l.svcCtx).GetCategoryDetail(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
