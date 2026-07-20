package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/order-service/internal/httpapi/admin"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminList2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminList2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminList2Logic {
	return &AdminList2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminList2Logic) AdminList2(w http.ResponseWriter, r *http.Request) {
	hadmin.NewOrderHandler(l.svcCtx).AdminList(w, r)
}
