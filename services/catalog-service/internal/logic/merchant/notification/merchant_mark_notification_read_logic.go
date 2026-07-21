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

type MerchantMarkNotificationReadLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantMarkNotificationReadLogic(svcCtx *svc.ServiceContext) *MerchantMarkNotificationReadLogic {
	return &MerchantMarkNotificationReadLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *MerchantMarkNotificationReadLogic) MerchantMarkNotificationRead(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/merchant/notifications/:id/read", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hhandler.NewNotificationHandler(l.svcCtx).MarkRead)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
