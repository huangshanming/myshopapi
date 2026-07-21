package seckill

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
	hmerchant "mymall/services/merchant-service/internal/app/merchant"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantApplySeckillLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantApplySeckillLogic(svcCtx *svc.ServiceContext) *MerchantApplySeckillLogic {
	return &MerchantApplySeckillLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantApplySeckillLogic) MerchantApplySeckill(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/seckill/entries", nil, nil, req, hmerchant.NewSeckillHandler(l.svcCtx).MerchantApplySeckill)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
