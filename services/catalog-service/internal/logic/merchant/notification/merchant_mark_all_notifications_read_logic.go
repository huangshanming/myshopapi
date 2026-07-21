package notification

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hhandler "mymall/services/catalog-service/internal/notify/handler"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantMarkAllNotificationsReadLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantMarkAllNotificationsReadLogic(svcCtx *svc.ServiceContext) *MerchantMarkAllNotificationsReadLogic {
	return &MerchantMarkAllNotificationsReadLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantMarkAllNotificationsReadLogic) MerchantMarkAllNotificationsRead(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/notifications/read-all", nil, nil, req, hhandler.NewNotificationHandler(l.svcCtx).MarkAllRead)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
