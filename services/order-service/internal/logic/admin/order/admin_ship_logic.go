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

type AdminShipLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminShipLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminShipLogic {
	return &AdminShipLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminShipLogic) AdminShip(ctx context.Context, req *types.ShipBodyReq) (*types.EmptyResp, error) {
	if err := biz.NewOrderLogic(l.svcCtx).Ship(ctx, req.Id, 0, req.ShipCompany, req.ShipNo); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
