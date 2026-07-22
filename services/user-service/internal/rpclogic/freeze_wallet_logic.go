package rpclogic

import (
	"context"

	userv1 "mymall/api/gen/user/v1"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FreezeWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFreezeWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FreezeWalletLogic {
	return &FreezeWalletLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *FreezeWalletLogic) FreezeWallet(in *userv1.WalletOpRequest) (*userv1.EmptyResponse, error) {
	if err := biz.NewWalletLogic(l.svcCtx).FreezeForOrder(l.ctx, in.GetUserId(), in.GetAmount(), in.GetOrderId(), in.GetOrderNo()); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &userv1.EmptyResponse{}, nil
}
