package task

import (
	"context"
	"mymall/pkg/appinput"
	hinternal "mymall/services/user-service/internal/app/internalapi"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalEventLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalEventLogic {
	return &InternalEventLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalEventLogic) InternalEvent(ctx context.Context, req *types.TaskEventReq) error {
	_, err := hinternal.NewTaskHandler(l.svcCtx).InternalEvent(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return err
	}
	return nil
}
