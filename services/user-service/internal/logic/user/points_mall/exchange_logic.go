package points_mall

import (
	"context"
	"mymall/pkg/httpinvoke"
	huser "mymall/services/user-service/internal/app/user"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewExchangeLogic(svcCtx *svc.ServiceContext) *ExchangeLogic {
	return &ExchangeLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *ExchangeLogic) Exchange(ctx context.Context, req *types.ExchangeReq) (resp *types.AnyResp, err error) {
	raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/user/points-mall/exchange", nil, nil, req, huser.NewPointsOrderHandler(l.svcCtx).Exchange)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
