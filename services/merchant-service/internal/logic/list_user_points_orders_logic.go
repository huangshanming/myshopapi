package logic

import (
	"net/http"

	"context"

	huser "mymall/services/merchant-service/internal/httpapi/user"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserPointsOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUserPointsOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserPointsOrdersLogic {
	return &ListUserPointsOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUserPointsOrdersLogic) ListUserPointsOrders(w http.ResponseWriter, r *http.Request) {
	huser.NewPointsOrderHandler(l.svcCtx).List(w, r)
}
