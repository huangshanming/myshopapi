package points_mall

import (
	"context"
	"fmt"
	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/user-service/internal/app/admin"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShipPointsOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewShipPointsOrderLogic(svcCtx *svc.ServiceContext) *ShipPointsOrderLogic {
	return &ShipPointsOrderLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *ShipPointsOrderLogic) ShipPointsOrder(ctx context.Context, req *types.ShipReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/points-orders/{Id}/ship", map[string]string{"id": fmt.Sprintf("%v", req.Id)}, nil, req, hadmin.NewPointsOrderHandler(l.svcCtx).Ship)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
