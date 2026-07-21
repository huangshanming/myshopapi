package logistics

import (
	"context"
	"net/http"

	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogisticsOptionsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewLogisticsOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogisticsOptionsLogic {
	return &LogisticsOptionsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *LogisticsOptionsLogic) LogisticsOptions(ctx context.Context) (*types.AnyResp, error) {
	list, err := biz.NewLogisticsLogic(l.svcCtx).Options(ctx, "")
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: list}, nil
}
