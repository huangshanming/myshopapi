package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	notifyhandler "mymall/services/catalog-service/internal/notify/handler"
	"mymall/services/catalog-service/internal/svc"
)

type MarkAllReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarkAllReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkAllReadLogic {
	return &MarkAllReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkAllReadLogic) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	notifyhandler.NewNotificationHandler(l.svcCtx).MarkAllRead(w, r)
}
