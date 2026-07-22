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

type LockCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLockCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LockCouponLogic {
	return &LockCouponLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *LockCouponLogic) LockCoupon(in *merchantv1.LockCouponRequest) (*merchantv1.EmptyResponse, error) {
	if err := biz.NewMerchantLogic(l.svcCtx).LockCoupon(in.GetUserCouponId(), in.GetUserId(), in.GetOrderId(), in.GetDiscountAmount()); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &merchantv1.EmptyResponse{}, nil
}
