package notification

import (
	"context"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminSendNotificationLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminSendNotificationLogic(svcCtx *svc.ServiceContext) *AdminSendNotificationLogic {
	return &AdminSendNotificationLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminSendNotificationLogic) AdminSendNotification(ctx context.Context, req *types.AdminSendReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/notifications/send", nil, nil, req, hadmin.NewAdminHandler(l.svcCtx).AdminSendNotification)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
