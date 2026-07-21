package task

import (
	"context"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserReportEventLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserReportEventLogic(svcCtx *svc.ServiceContext) *UserReportEventLogic {
	return &UserReportEventLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserReportEventLogic) UserReportEvent(ctx context.Context, req *types.TaskEventReq) error {
	_, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/tasks/events", nil, nil, req, huser.NewTaskHandler(l.svcCtx).UserReportEvent)
	if err != nil {
		return err
	}
	return nil
}
