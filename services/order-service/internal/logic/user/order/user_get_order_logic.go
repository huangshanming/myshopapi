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

type UserGetOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUserGetOrderLogic(svcCtx *svc.ServiceContext) *UserGetOrderLogic {
	return &UserGetOrderLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *UserGetOrderLogic) UserGetOrder(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/orders/:id", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, huser.NewOrderHandler(l.svcCtx).Detail)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
