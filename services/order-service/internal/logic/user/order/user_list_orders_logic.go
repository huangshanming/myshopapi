package order

import (
	"context"
	"net/http"

	"mymall/pkg/middleware"
	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListOrdersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListOrdersLogic {
	return &UserListOrdersLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserListOrdersLogic) UserListOrders(ctx context.Context, req *types.PageReq) (*types.PageListResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	page, pageSize, _ := pagination.Normalize(&pagination.PageReq{Page: req.Page, PageSize: req.PageSize})
	// status comes from query; PageReq may not include it — keep empty for now
	orders, total, err := biz.NewOrderLogic(l.svcCtx).ListOrders(ctx, userID, page, pageSize, "")
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: orders}, nil
}
