package notification

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnreadNotificationCountLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUnreadNotificationCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnreadNotificationCountLogic {
	return &UnreadNotificationCountLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UnreadNotificationCountLogic) UnreadNotificationCount(ctx context.Context) (*types.CountResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	n, err := biz.NewUserLogic(l.svcCtx).UnreadCount(ctx, userID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.CountResp{Count: n}, nil
}
