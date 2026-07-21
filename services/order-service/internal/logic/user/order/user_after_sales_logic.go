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

type UserAfterSalesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserAfterSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAfterSalesLogic {
	return &UserAfterSalesLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *UserAfterSalesLogic) UserAfterSales(ctx context.Context, req *types.PageReq) (*types.PageListResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	page, pageSize, _ := pagination.Normalize(&pagination.PageReq{Page: req.Page, PageSize: req.PageSize})
	list, total, err := biz.NewOrderLogic(l.svcCtx).ListUserAfterSales(ctx, userID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil
}
