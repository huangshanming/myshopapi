// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notification

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalCreateNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalCreateNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalCreateNotificationLogic {
	return &InternalCreateNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InternalCreateNotificationLogic) InternalCreateNotification(req *types.NotifyCreateReq) (resp *types.NotificationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
