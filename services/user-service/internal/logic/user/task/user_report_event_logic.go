package task

import (
	"context"
	"mymall/pkg/appinput"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserReportEventLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserReportEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserReportEventLogic {
	return &UserReportEventLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UserReportEventLogic) UserReportEvent(ctx context.Context, req *types.TaskEventReq) error {
	_, err := huser.NewTaskHandler(l.svcCtx).UserReportEvent(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return err
	}
	return nil
}
