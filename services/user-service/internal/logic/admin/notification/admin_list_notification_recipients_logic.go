package notification

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
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
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	page, pageSize, _ := pagination.Normalize(&pagination.PageReq{Page: req.Page, PageSize: req.PageSize})
	userLogic := biz.NewUserLogic(l.svcCtx)
	batch, err := userLogic.GetSendBatch(ctx, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "发送记录不存在")
	}
	list, total, err := userLogic.ListBatchRecipients(ctx, req.Id, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: map[string]interface{}{
		"batch": batch,
		"list":  list,
		"total": total,
	}}, nil
}
