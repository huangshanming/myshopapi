package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/order-service/internal/httpapi/admin"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminList3Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminList3Logic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminList3Logic {
	return &AdminList3Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminList3Logic) AdminList3(w http.ResponseWriter, r *http.Request) {
	hadmin.NewLogisticsHandler(l.svcCtx).Create(w, r)
}
