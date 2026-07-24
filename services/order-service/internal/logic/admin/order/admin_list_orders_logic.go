package order

import (
	"context"
	"net/http"

	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListOrdersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListOrdersLogic {
	return &AdminListOrdersLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminListOrdersLogic) AdminListOrders(ctx context.Context, req *types.PageReq) (*types.PageListResp, error) {
	page, pageSize, _ := pagination.Normalize(&pagination.PageReq{Page: req.Page, PageSize: req.PageSize})
	orders, total, err := biz.NewOrderLogic(l.svcCtx).ListAll(ctx, 0, page, pageSize, req.Status, "")
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: orders}, nil
}
