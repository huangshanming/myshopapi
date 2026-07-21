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

type MerchantShipLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantShipLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantShipLogic {
	return &MerchantShipLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *MerchantShipLogic) MerchantShip(ctx context.Context, req *types.ShipBodyReq) (*types.EmptyResp, error) {
	shopID := middleware.GetShopID(ctx)
	if err := biz.NewOrderLogic(l.svcCtx).Ship(ctx, req.Id, shopID, req.ShipCompany, req.ShipNo); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
