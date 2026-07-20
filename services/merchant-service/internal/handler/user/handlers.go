package user

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

type PointsOrderHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.PointsOrderLogic
}

func NewPointsOrderHandler(svcCtx *svc.ServiceContext) *PointsOrderHandler {
	return &PointsOrderHandler{
		svcCtx: svcCtx,
		logic:  logic.NewPointsOrderLogic(context.Background(), svcCtx),
	}
}
