package notification

import (
	"context"
	"mymall/pkg/httpinvoke"
	hinternal "mymall/services/user-service/internal/app/internalapi"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalCreateNotificationLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalCreateNotificationLogic(svcCtx *svc.ServiceContext) *InternalCreateNotificationLogic {
	return &InternalCreateNotificationLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *InternalCreateNotificationLogic) InternalCreateNotification(ctx context.Context, req *types.NotifyCreateReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/internal/notifications", nil, nil, req, hinternal.NewNotificationHandler(l.svcCtx).InternalCreateNotification)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
