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

type AdminListNotificationSendsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListNotificationSendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListNotificationSendsLogic {
	return &AdminListNotificationSendsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListNotificationSendsLogic) AdminListNotificationSends(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	page, pageSize, _ := pagination.Normalize(&pagination.PageReq{Page: req.Page, PageSize: req.PageSize})
	list, total, err := biz.NewUserLogic(l.svcCtx).ListSendBatches(ctx, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: total}, nil
}
