package notification

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkAllNotificationsReadLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMarkAllNotificationsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkAllNotificationsReadLogic {
	return &MarkAllNotificationsReadLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *MarkAllNotificationsReadLogic) MarkAllNotificationsRead(ctx context.Context) error {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return xerr.New(http.StatusUnauthorized, "未授权")
	}
	if err := biz.NewUserLogic(l.svcCtx).MarkAllRead(ctx, userID); err != nil {
		return xerr.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}
