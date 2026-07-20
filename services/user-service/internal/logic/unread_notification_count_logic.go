package logic

import (
	"net/http"

	"context"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"

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

func (l *UnreadNotificationCountLogic) UnreadNotificationCount(w http.ResponseWriter, r *http.Request) {
	huser.NewUserHandler(l.svcCtx).UnreadNotificationCount(w, r)
}
