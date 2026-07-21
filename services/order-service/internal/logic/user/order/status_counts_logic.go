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

type StatusCountsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewStatusCountsLogic(svcCtx *svc.ServiceContext) *StatusCountsLogic {
	return &StatusCountsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *StatusCountsLogic) StatusCounts(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/orders/status-counts", nil, nil, nil, huser.NewOrderHandler(l.svcCtx).StatusCounts)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
