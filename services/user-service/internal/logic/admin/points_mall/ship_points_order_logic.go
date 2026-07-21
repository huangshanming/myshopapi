package points_mall

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type ShipPointsOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewShipPointsOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShipPointsOrderLogic {
	return &ShipPointsOrderLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ShipPointsOrderLogic) ShipPointsOrder(ctx context.Context, req *types.ShipReq) (resp *types.AnyResp, err error) {
	if req.Id == 0 {
		return nil, xerr.New(http.StatusBadRequest, "订单ID无效")
	}
	o, err := biz.NewPointsOrderLogic(l.svcCtx).AdminShip(ctx, req.Id, biz.ShipReq{
		ShipCompany: req.ExpressCompany,
		ShipNo:      req.ExpressNo,
	})
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: o}, nil
}
