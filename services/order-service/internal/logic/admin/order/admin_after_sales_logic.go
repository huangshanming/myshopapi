package order

import (
	"context"
	"net/http"

	"mymall/pkg/pagination"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/repository"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminAfterSalesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminAfterSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminAfterSalesLogic {
	return &AdminAfterSalesLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminAfterSalesLogic) AdminAfterSales(ctx context.Context, req *types.PageReq) (*types.PageListResp, error) {
	page, pageSize, _ := pagination.Normalize(&pagination.PageReq{Page: req.Page, PageSize: req.PageSize})
	list, total, err := biz.NewOrderLogic(l.svcCtx).ListAfterSales(ctx, repository.AfterSaleListFilter{
		Page: page, PageSize: pageSize,
	})
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil
}
