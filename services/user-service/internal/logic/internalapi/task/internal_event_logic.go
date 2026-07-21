package task

import (
	"context"
	"mymall/pkg/httpinvoke"
	hinternal "mymall/services/user-service/internal/app/internalapi"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalEventLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalEventLogic(svcCtx *svc.ServiceContext) *InternalEventLogic {
	return &InternalEventLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *InternalEventLogic) InternalEvent(ctx context.Context, req *types.TaskEventReq) error {
	_, err := httpinvoke.Run(ctx, "POST", "/api/v1/internal/tasks/events", nil, nil, req, hinternal.NewTaskHandler(l.svcCtx).InternalEvent)
	if err != nil {
		return err
	}
	return nil
}
