package logic

import (
	"net/http"

	"context"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"

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

func (l *AdminListNotificationRecipientsLogic) AdminListNotificationRecipients(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewAdminHandler(l.svcCtx).AdminListNotificationRecipients
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "business:message:send"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
