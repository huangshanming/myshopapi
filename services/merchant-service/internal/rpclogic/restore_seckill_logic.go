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

type RestoreSeckillLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRestoreSeckillLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreSeckillLogic {
	return &RestoreSeckillLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *RestoreSeckillLogic) RestoreSeckill(in *merchantv1.RestoreSeckillRequest) (*merchantv1.EmptyResponse, error) {
	if err := biz.NewMerchantLogic(l.svcCtx).RestoreSeckill(in.GetEntryId(), int(in.GetQuantity())); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &merchantv1.EmptyResponse{}, nil
}
