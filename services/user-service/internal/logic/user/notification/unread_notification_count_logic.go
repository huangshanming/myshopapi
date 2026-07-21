package notification

import (
	"context"
	"encoding/json"
	"mymall/pkg/appinput"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnreadNotificationCountLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUnreadNotificationCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnreadNotificationCountLogic {
	return &UnreadNotificationCountLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UnreadNotificationCountLogic) UnreadNotificationCount(ctx context.Context) (resp *types.CountResp, err error) {
	data, err := huser.NewUserHandler(l.svcCtx).UnreadNotificationCount(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(data)
	var out types.CountResp
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
