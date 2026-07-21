package notification

import (
	"context"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnreadNotificationCountLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUnreadNotificationCountLogic(svcCtx *svc.ServiceContext) *UnreadNotificationCountLogic {
	return &UnreadNotificationCountLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UnreadNotificationCountLogic) UnreadNotificationCount(ctx context.Context) (resp *types.CountResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/user/notifications/unread-count", nil, nil, nil, huser.NewUserHandler(l.svcCtx).UnreadNotificationCount)
	if err != nil {
		return nil, err
	}
	var out types.CountResp
	if err := httpinvoke.Decode(raw, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
