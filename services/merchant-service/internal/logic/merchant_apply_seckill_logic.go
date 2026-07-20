package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantApplySeckillLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantApplySeckillLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantApplySeckillLogic {
	return &MerchantApplySeckillLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantApplySeckillLogic) MerchantApplySeckill(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewSeckillHandler(l.svcCtx).MerchantApplySeckill(w, r)
}
