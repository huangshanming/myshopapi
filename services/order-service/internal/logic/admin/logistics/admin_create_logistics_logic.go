package logistics

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/order-service/internal/httpapi/admin"
	"mymall/services/order-service/internal/svc"
)

type AdminCreateLogisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCreateLogisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateLogisticsLogic {
	return &AdminCreateLogisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateLogisticsLogic) AdminCreateLogistics(w http.ResponseWriter, r *http.Request) {
	hadmin.NewLogisticsHandler(l.svcCtx).Create(w, r)
}
