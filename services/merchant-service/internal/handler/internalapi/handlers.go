package internalapi

import (
	"context"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

type CouponHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.MerchantLogic
}

func NewCouponHandler(svcCtx *svc.ServiceContext) *CouponHandler {
	return &CouponHandler{
		svcCtx: svcCtx,
		logic:  logic.NewMerchantLogic(context.Background(), svcCtx),
	}
}

type SeckillHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.MerchantLogic
}

func NewSeckillHandler(svcCtx *svc.ServiceContext) *SeckillHandler {
	return &SeckillHandler{
		svcCtx: svcCtx,
		logic:  logic.NewMerchantLogic(context.Background(), svcCtx),
	}
}
