package logic

import (
	"net/http"

	"context"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatusCountsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatusCountsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatusCountsLogic {
	return &StatusCountsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatusCountsLogic) StatusCounts(w http.ResponseWriter, r *http.Request) {
	huser.NewOrderHandler(l.svcCtx).StatusCounts(w, r)
}
