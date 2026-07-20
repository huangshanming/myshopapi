package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	notifyhandler "mymall/services/catalog-service/internal/notify/handler"
	"mymall/services/catalog-service/internal/svc"
)

type List5Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewList5Logic(ctx context.Context, svcCtx *svc.ServiceContext) *List5Logic {
	return &List5Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *List5Logic) List5(w http.ResponseWriter, r *http.Request) {
	notifyhandler.NewNotificationHandler(l.svcCtx).List(w, r)
}
