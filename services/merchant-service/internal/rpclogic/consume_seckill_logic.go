package rpclogic

import (
	"context"

	merchantv1 "mymall/api/gen/merchant/v1"
	"mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConsumeSeckillLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConsumeSeckillLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConsumeSeckillLogic {
	return &ConsumeSeckillLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ConsumeSeckillLogic) ConsumeSeckill(in *merchantv1.ConsumeSeckillRequest) (*merchantv1.ConsumeSeckillResponse, error) {
	out, err := biz.NewMerchantLogic(l.svcCtx).ConsumeSeckill(in.GetEntryId(), in.GetProductId(), int(in.GetQuantity()))
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	price, _ := out["seckill_price"].(float64)
	return &merchantv1.ConsumeSeckillResponse{
		EntryId: in.GetEntryId(), ProductId: in.GetProductId(), SeckillPrice: price, Quantity: in.GetQuantity(),
	}, nil
}
