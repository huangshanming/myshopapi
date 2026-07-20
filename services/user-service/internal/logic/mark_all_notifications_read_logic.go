package logic

import (
	"net/http"

	"context"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkAllNotificationsReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarkAllNotificationsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkAllNotificationsReadLogic {
	return &MarkAllNotificationsReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkAllNotificationsReadLogic) MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	huser.NewUserHandler(l.svcCtx).MarkAllNotificationsRead(w, r)
}
