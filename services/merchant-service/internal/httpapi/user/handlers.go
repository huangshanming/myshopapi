package user

import (
	"context"

	"mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/svc"
)

type CouponHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.MerchantLogic
}

func NewCouponHandler(svcCtx *svc.ServiceContext) *CouponHandler {
	return &CouponHandler{
		svcCtx: svcCtx,
		logic:  biz.NewMerchantLogic(context.Background(), svcCtx),
	}
}

type PointsOrderHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.PointsOrderLogic
}

func NewPointsOrderHandler(svcCtx *svc.ServiceContext) *PointsOrderHandler {
	return &PointsOrderHandler{
		svcCtx: svcCtx,
		logic:  biz.NewPointsOrderLogic(context.Background(), svcCtx),
	}
}
