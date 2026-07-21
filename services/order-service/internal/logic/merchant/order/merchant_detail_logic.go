package order

import (
	"context"
	"net/http"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantDetailLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantDetailLogic {
	return &MerchantDetailLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *MerchantDetailLogic) MerchantDetail(ctx context.Context, req *types.IdPathReq) (*types.AnyResp, error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少 shop_id")
	}
	ol := biz.NewOrderLogic(l.svcCtx)
	order, err := ol.GetOrderByShop(ctx, shopID, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "订单不存在")
	}
	as, _ := ol.ListAfterSalesByOrder(ctx, req.Id)
	return &types.AnyResp{Data: map[string]interface{}{"order": order, "after_sales": as}}, nil
}
