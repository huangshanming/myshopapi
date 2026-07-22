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

type TaskEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTaskEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskEventLogic {
	return &TaskEventLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *TaskEventLogic) TaskEvent(in *userv1.TaskEventRequest) (*userv1.EmptyResponse, error) {
	if in.GetUserId() == 0 || in.GetTaskCode() == "" {
		return &userv1.EmptyResponse{}, nil
	}
	delta := int(in.GetDelta())
	if delta < 1 {
		delta = 1
	}
	if err := biz.NewTaskLogic(l.svcCtx).HandleEvent(l.ctx, biz.TaskEventReq{
		UserID: in.GetUserId(), TaskCode: in.GetTaskCode(), Delta: delta,
		RefType: in.GetRefType(), RefID: in.GetRefId(),
	}); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}
	return &userv1.EmptyResponse{}, nil
}
