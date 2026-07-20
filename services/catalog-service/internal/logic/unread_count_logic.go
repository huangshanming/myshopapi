package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	notifyhandler "mymall/services/catalog-service/internal/notify/handler"
	"mymall/services/catalog-service/internal/svc"
)

type UnreadCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnreadCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnreadCountLogic {
	return &UnreadCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnreadCountLogic) UnreadCount(w http.ResponseWriter, r *http.Request) {
	notifyhandler.NewNotificationHandler(l.svcCtx).UnreadCount(w, r)
}
