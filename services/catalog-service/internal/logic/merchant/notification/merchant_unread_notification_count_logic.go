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

type MerchantUnreadNotificationCountLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantUnreadNotificationCountLogic(svcCtx *svc.ServiceContext) *MerchantUnreadNotificationCountLogic {
	return &MerchantUnreadNotificationCountLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantUnreadNotificationCountLogic) MerchantUnreadNotificationCount(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/merchant/notifications/unread-count", nil, nil, nil, hhandler.NewNotificationHandler(l.svcCtx).UnreadCount)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
