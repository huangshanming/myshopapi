package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/order-service/internal/httpapi/admin"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type Create3Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreate3Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Create3Logic {
	return &Create3Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Create3Logic) Create3(w http.ResponseWriter, r *http.Request) {
	hadmin.NewLogisticsHandler(l.svcCtx).Create(w, r)
}
