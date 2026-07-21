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

type MerchantSetSeckillAutoRenewLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantSetSeckillAutoRenewLogic(svcCtx *svc.ServiceContext) *MerchantSetSeckillAutoRenewLogic {
	return &MerchantSetSeckillAutoRenewLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantSetSeckillAutoRenewLogic) MerchantSetSeckillAutoRenew(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/merchant/seckill/entries/:id/auto-renew", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hmerchant.NewSeckillHandler(l.svcCtx).MerchantSetSeckillAutoRenew)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
