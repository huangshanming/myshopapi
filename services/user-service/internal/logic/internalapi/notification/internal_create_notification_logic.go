package notification

import (
	"context"
	"mymall/pkg/appinput"
	hinternal "mymall/services/user-service/internal/app/internalapi"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *InternalCreateNotificationLogic) InternalCreateNotification(ctx context.Context, req *types.NotifyCreateReq) (resp *types.AnyResp, err error) {
	data, err := hinternal.NewNotificationHandler(l.svcCtx).InternalCreateNotification(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
