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

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetUserLogic) GetUser(in *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := biz.NewUserLogic(l.svcCtx).GetProfile(l.ctx, in.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &userv1.GetUserResponse{
		Id: user.ID, Mobile: user.Mobile, Nickname: user.Nickname, Status: int32(user.Status),
	}, nil
}
