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

type GetAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAddressLogic {
	return &GetAddressLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetAddressLogic) GetAddress(in *userv1.GetAddressRequest) (*userv1.GetAddressResponse, error) {
	a, err := biz.NewAddressLogic(l.svcCtx).Get(l.ctx, in.GetUserId(), in.GetAddressId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%s", err.Error())
	}
	return &userv1.GetAddressResponse{
		Id: a.ID, UserId: a.UserID, ReceiverName: a.ReceiverName, ReceiverPhone: a.ReceiverPhone,
		Province: a.Province, City: a.City, District: a.District, Detail: a.Detail, IsDefault: int32(a.IsDefault),
	}, nil
}
