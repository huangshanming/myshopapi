package notification

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
	"net/url"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListNotificationRecipientsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListNotificationRecipientsLogic(svcCtx *svc.ServiceContext) *AdminListNotificationRecipientsLogic {
	return &AdminListNotificationRecipientsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminListNotificationRecipientsLogic) AdminListNotificationRecipients(ctx context.Context, req *types.NotificationRecipientsReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/notifications/sends/{Id}/recipients", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}, nil, hadmin.NewAdminHandler(l.svcCtx).AdminListNotificationRecipients)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
