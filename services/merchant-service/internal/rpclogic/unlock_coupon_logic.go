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

type UnlockCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnlockCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlockCouponLogic {
	return &UnlockCouponLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UnlockCouponLogic) UnlockCoupon(in *merchantv1.UnlockCouponRequest) (*merchantv1.EmptyResponse, error) {
	if err := biz.NewMerchantLogic(l.svcCtx).UnlockCoupon(in.GetUserCouponId(), in.GetOrderId()); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &merchantv1.EmptyResponse{}, nil
}
