package order

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/order-service/internal/app/admin"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminShipLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminShipLogic(svcCtx *svc.ServiceContext) *AdminShipLogic {
	return &AdminShipLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminShipLogic) AdminShip(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/orders/:id/ship", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hadmin.NewOrderHandler(l.svcCtx).AdminShip)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
