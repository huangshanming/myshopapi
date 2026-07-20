package notification

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	notifyhandler "mymall/services/catalog-service/internal/notify/handler"
	"mymall/services/catalog-service/internal/svc"
)

type MerchantMarkAllNotificationsReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantMarkAllNotificationsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantMarkAllNotificationsReadLogic {
	return &MerchantMarkAllNotificationsReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantMarkAllNotificationsReadLogic) MerchantMarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	notifyhandler.NewNotificationHandler(l.svcCtx).MarkAllRead(w, r)
}
