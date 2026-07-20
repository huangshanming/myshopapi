package notification

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	notifyhandler "mymall/services/catalog-service/internal/notify/handler"
	"mymall/services/catalog-service/internal/svc"
)

type MerchantListNotificationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantListNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListNotificationsLogic {
	return &MerchantListNotificationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantListNotificationsLogic) MerchantListNotifications(w http.ResponseWriter, r *http.Request) {
	notifyhandler.NewNotificationHandler(l.svcCtx).List(w, r)
}
