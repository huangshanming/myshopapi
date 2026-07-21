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

type MerchantSeckillSessionsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantSeckillSessionsLogic(svcCtx *svc.ServiceContext) *MerchantSeckillSessionsLogic {
	return &MerchantSeckillSessionsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantSeckillSessionsLogic) MerchantSeckillSessions(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/merchant/seckill/sessions", nil, nil, nil, hmerchant.NewSeckillHandler(l.svcCtx).MerchantSeckillSessions)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
