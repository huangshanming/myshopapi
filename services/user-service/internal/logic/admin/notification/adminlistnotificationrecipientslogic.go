// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notification

import (
	"context"

	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListNotificationRecipientsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListNotificationRecipientsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListNotificationRecipientsLogic {
	return &AdminListNotificationRecipientsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListNotificationRecipientsLogic) AdminListNotificationRecipients(req *types.NotificationRecipientsReq) (resp *types.NotificationRecipientsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
