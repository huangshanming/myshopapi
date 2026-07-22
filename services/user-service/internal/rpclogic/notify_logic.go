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

type NotifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyLogic {
	return &NotifyLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *NotifyLogic) Notify(in *userv1.NotifyRequest) (*userv1.EmptyResponse, error) {
	if in.GetUserId() == 0 || in.GetTitle() == "" {
		return &userv1.EmptyResponse{}, nil
	}
	_, err := biz.NewUserLogic(l.svcCtx).CreateNotification(l.ctx, biz.NotifyCreateReq{
		UserID: in.GetUserId(), Title: in.GetTitle(), Content: in.GetContent(),
		MsgType: in.GetMsgType(), LinkType: in.GetLinkType(), LinkID: in.GetLinkId(), Extra: in.GetExtra(),
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}
	return &userv1.EmptyResponse{}, nil
}
