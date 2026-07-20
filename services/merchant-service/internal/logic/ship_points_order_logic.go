package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShipPointsOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShipPointsOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShipPointsOrderLogic {
	return &ShipPointsOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShipPointsOrderLogic) ShipPointsOrder(w http.ResponseWriter, r *http.Request) {
	hadmin.NewPointsOrderHandler(l.svcCtx).Ship(w, r)
}
