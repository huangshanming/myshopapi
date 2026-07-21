package order

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"net/url"

	hadmin "mymall/services/order-service/internal/app/admin"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminAfterSalesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminAfterSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminAfterSalesLogic {
	return &AdminAfterSalesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminAfterSalesLogic) AdminAfterSales(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hadmin.NewOrderHandler(l.svcCtx).AdminAfterSales(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
