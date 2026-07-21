package notification

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNotificationsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotificationsLogic {
	return &ListNotificationsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *ListNotificationsLogic) ListNotifications(ctx context.Context, req *types.PageReq) (*types.PageListResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	page, pageSize, _ := pagination.Normalize(&pagination.PageReq{Page: req.Page, PageSize: req.PageSize})
	list, total, err := biz.NewUserLogic(l.svcCtx).ListMyNotifications(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil
}
