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

type OrderGiftCouponsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderGiftCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderGiftCouponsLogic {
	return &OrderGiftCouponsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *OrderGiftCouponsLogic) OrderGiftCoupons(in *merchantv1.OrderGiftCouponsRequest) (*merchantv1.OrderGiftCouponsResponse, error) {
	n, err := biz.NewMerchantLogic(l.svcCtx).OrderGiftCoupons(in.GetUserId(), in.GetShopId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	return &merchantv1.OrderGiftCouponsResponse{Granted: int32(n)}, nil
}
