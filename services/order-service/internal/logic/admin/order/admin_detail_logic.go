package order

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDetailLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDetailLogic {
	return &AdminDetailLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminDetailLogic) AdminDetail(ctx context.Context, req *types.IdPathReq) (*types.AnyResp, error) {
	ol := biz.NewOrderLogic(l.svcCtx)
	order, err := ol.GetOrderAdmin(ctx, req.Id)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, "订单不存在")
	}
	as, _ := ol.ListAfterSalesByOrder(ctx, req.Id)
	return &types.AnyResp{Data: map[string]interface{}{"order": order, "after_sales": as}}, nil
}
