package order

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	huser "mymall/services/order-service/internal/app/user"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAfterSaleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateAfterSaleLogic(svcCtx *svc.ServiceContext) *CreateAfterSaleLogic {
	return &CreateAfterSaleLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *CreateAfterSaleLogic) CreateAfterSale(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/orders/:id/after-sales", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, huser.NewOrderHandler(l.svcCtx).CreateAfterSale)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
