package order

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	huser "mymall/services/order-service/internal/app/user"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAfterSaleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateAfterSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAfterSaleLogic {
	return &CreateAfterSaleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreateAfterSaleLogic) CreateAfterSale(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := huser.NewOrderHandler(l.svcCtx).CreateAfterSale(ctx, appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
