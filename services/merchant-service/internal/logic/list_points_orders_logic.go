package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPointsOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPointsOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPointsOrdersLogic {
	return &ListPointsOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPointsOrdersLogic) ListPointsOrders(w http.ResponseWriter, r *http.Request) {
	hadmin.NewPointsOrderHandler(l.svcCtx).List(w, r)
}
