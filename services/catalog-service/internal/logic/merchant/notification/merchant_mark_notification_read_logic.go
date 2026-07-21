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

type MerchantMarkNotificationReadLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantMarkNotificationReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantMarkNotificationReadLogic {
	return &MerchantMarkNotificationReadLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantMarkNotificationReadLogic) MerchantMarkNotificationRead(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hhandler.NewNotificationHandler(l.svcCtx).MarkRead(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
