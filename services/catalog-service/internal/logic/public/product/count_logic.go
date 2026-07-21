package product

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	huser "mymall/services/catalog-service/internal/product/app/user"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCountLogic(svcCtx *svc.ServiceContext) *CountLogic {
	return &CountLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *CountLogic) Count(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/products/:id/favorite-count", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, huser.NewFavoriteHandler(l.svcCtx).Count)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
