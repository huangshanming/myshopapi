package internalapi

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hinternal "mymall/services/user-service/internal/httpapi/internalapi"
	"mymall/services/user-service/internal/svc"
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

func (l *InternalCreateNotificationLogic) InternalCreateNotification(w http.ResponseWriter, r *http.Request) {
	hinternal.NewNotificationHandler(l.svcCtx).InternalCreateNotification(w, r)
}
