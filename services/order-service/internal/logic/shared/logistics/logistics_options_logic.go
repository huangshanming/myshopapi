package logistics

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

type LogisticsOptionsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewLogisticsOptionsLogic(svcCtx *svc.ServiceContext) *LogisticsOptionsLogic {
	return &LogisticsOptionsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *LogisticsOptionsLogic) LogisticsOptions(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/logistics/options", nil, nil, nil, hadmin.NewLogisticsHandler(l.svcCtx).Options)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
