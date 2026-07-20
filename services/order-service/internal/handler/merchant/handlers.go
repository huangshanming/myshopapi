package merchant

import (
	"context"

	"mymall/services/order-service/internal/logic"
	"mymall/services/order-service/internal/svc"
)

type OrderHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.OrderLogic
}

func NewOrderHandler(svcCtx *svc.ServiceContext) *OrderHandler {
	return &OrderHandler{
		svcCtx: svcCtx,
		logic:  logic.NewOrderLogic(context.Background(), svcCtx),
	}
}

type ReviewHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.ReviewLogic
}

func NewReviewHandler(svcCtx *svc.ServiceContext) *ReviewHandler {
	return &ReviewHandler{
		svcCtx: svcCtx,
		logic:  logic.NewReviewLogic(context.Background(), svcCtx),
	}
}
