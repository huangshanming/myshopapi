package notification

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
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

func NewAdminListNotificationRecipientsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListNotificationRecipientsLogic {
	return &AdminListNotificationRecipientsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListNotificationRecipientsLogic) AdminListNotificationRecipients(ctx context.Context, req *types.NotificationRecipientsReq) (resp *types.AnyResp, err error) {
	data, err := hadmin.NewAdminHandler(l.svcCtx).AdminListNotificationRecipients(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%v", req.Id)}, Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
