package logistics

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

type LogisticsOptionsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewLogisticsOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogisticsOptionsLogic {
	return &LogisticsOptionsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *LogisticsOptionsLogic) LogisticsOptions(ctx context.Context) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	data, err := hadmin.NewLogisticsHandler(l.svcCtx).Options(ctx, appinput.CallInput{})
	if err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
