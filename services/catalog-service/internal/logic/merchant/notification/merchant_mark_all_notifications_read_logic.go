package notification

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hhandler "mymall/services/catalog-service/internal/notify/handler"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantMarkAllNotificationsReadLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantMarkAllNotificationsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantMarkAllNotificationsReadLogic {
	return &MerchantMarkAllNotificationsReadLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantMarkAllNotificationsReadLogic) MerchantMarkAllNotificationsRead(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hhandler.NewNotificationHandler(l.svcCtx).MarkAllRead(ctx, appinput.CallInput{Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
