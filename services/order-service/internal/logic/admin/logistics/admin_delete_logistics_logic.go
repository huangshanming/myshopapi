package logistics

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/order-service/internal/httpapi/admin"
	"mymall/services/order-service/internal/svc"
)

type AdminDeleteLogisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteLogisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteLogisticsLogic {
	return &AdminDeleteLogisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteLogisticsLogic) AdminDeleteLogistics(w http.ResponseWriter, r *http.Request) {
	hadmin.NewLogisticsHandler(l.svcCtx).Delete(w, r)
}
