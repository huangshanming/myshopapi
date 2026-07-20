package points_mall

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"
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
