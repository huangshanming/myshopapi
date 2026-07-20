package seckill

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"
)

type MerchantSetSeckillAutoRenewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantSetSeckillAutoRenewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantSetSeckillAutoRenewLogic {
	return &MerchantSetSeckillAutoRenewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantSetSeckillAutoRenewLogic) MerchantSetSeckillAutoRenew(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewSeckillHandler(l.svcCtx).MerchantSetSeckillAutoRenew(w, r)
}
