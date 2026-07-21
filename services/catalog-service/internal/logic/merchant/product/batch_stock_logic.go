package product

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hmerchant "mymall/services/catalog-service/internal/product/app/merchant"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchStockLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewBatchStockLogic(svcCtx *svc.ServiceContext) *BatchStockLogic {
	return &BatchStockLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *BatchStockLogic) BatchStock(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/skus/batch-stock", nil, nil, req, hmerchant.NewProductHandler(l.svcCtx).BatchStock)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
