package logistics

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/order-service/internal/httpapi/admin"
	"mymall/services/order-service/internal/svc"
)

type LogisticsOptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogisticsOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogisticsOptionsLogic {
	return &LogisticsOptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogisticsOptionsLogic) LogisticsOptions(w http.ResponseWriter, r *http.Request) {
	hadmin.NewLogisticsHandler(l.svcCtx).Options(w, r)
}
