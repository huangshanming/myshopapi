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

type RedeemCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRedeemCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedeemCouponLogic {
	return &RedeemCouponLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *RedeemCouponLogic) RedeemCoupon(in *merchantv1.RedeemCouponRequest) (*merchantv1.EmptyResponse, error) {
	if err := biz.NewMerchantLogic(l.svcCtx).RedeemCoupon(in.GetUserCouponId(), in.GetOrderId(), in.GetDiscountAmount()); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &merchantv1.EmptyResponse{}, nil
}
