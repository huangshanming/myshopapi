package notification

import (
	"context"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkAllNotificationsReadLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMarkAllNotificationsReadLogic(svcCtx *svc.ServiceContext) *MarkAllNotificationsReadLogic {
	return &MarkAllNotificationsReadLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MarkAllNotificationsReadLogic) MarkAllNotificationsRead(ctx context.Context) error {
	_, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/notifications/read-all", nil, nil, nil, huser.NewUserHandler(l.svcCtx).MarkAllNotificationsRead)
	if err != nil {
		return err
	}
	return nil
}
