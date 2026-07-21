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

type MerchantListOrdersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListOrdersLogic {
	return &MerchantListOrdersLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *MerchantListOrdersLogic) MerchantListOrders(ctx context.Context, req *types.PageReq) (*types.PageListResp, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少 shop_id")
	}
	page, pageSize, _ := pagination.Normalize(&pagination.PageReq{Page: req.Page, PageSize: req.PageSize})
	orders, total, err := biz.NewOrderLogic(l.svcCtx).ListByShop(ctx, shopID, page, pageSize, "", "")
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: orders}, nil
}
