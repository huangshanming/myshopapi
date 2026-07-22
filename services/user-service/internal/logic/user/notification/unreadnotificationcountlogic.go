// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notification

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnreadNotificationCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnreadNotificationCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnreadNotificationCountLogic {
	return &UnreadNotificationCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnreadNotificationCountLogic) UnreadNotificationCount() (resp *types.CountResp, err error) {
	// todo: add your logic here and delete this line

	return
}
