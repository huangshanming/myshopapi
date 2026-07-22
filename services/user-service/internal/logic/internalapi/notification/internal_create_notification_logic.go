package notification

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type InternalCreateNotificationLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalCreateNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalCreateNotificationLogic {
	return &InternalCreateNotificationLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalCreateNotificationLogic) InternalCreateNotification(ctx context.Context, req *types.NotifyCreateReq) (resp *types.NotificationResp, err error) {
	n, err := biz.NewUserLogic(l.svcCtx).CreateNotification(ctx, biz.NotifyCreateReq{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
		MsgType: req.Type,
	})
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.NotificationResp{Data: n}, nil
}
