package logic

import (
	"net/http"

	"context"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminSendNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminSendNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminSendNotificationLogic {
	return &AdminSendNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminSendNotificationLogic) AdminSendNotification(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewAdminHandler(l.svcCtx).AdminSendNotification
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "business:message:send"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
