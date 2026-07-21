package admin

import (
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
)

type LogisticsHandler struct {
	logic *biz.LogisticsLogic
}

func NewLogisticsHandler(svcCtx *svc.ServiceContext) *LogisticsHandler {
	return &LogisticsHandler{logic: biz.NewLogisticsLogic(svcCtx)}
}

type OrderHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.OrderLogic
}

func NewOrderHandler(svcCtx *svc.ServiceContext) *OrderHandler {
	return &OrderHandler{
		svcCtx: svcCtx,
		logic:  biz.NewOrderLogic(svcCtx),
	}
}

type ReviewHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.ReviewLogic
}

func NewReviewHandler(svcCtx *svc.ServiceContext) *ReviewHandler {
	return &ReviewHandler{
		svcCtx: svcCtx,
		logic:  biz.NewReviewLogic(svcCtx),
	}
}
