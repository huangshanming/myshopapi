package logic

import (
	"net/http"

	"context"

	huser "mymall/services/merchant-service/internal/httpapi/user"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailUserPointsOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailUserPointsOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailUserPointsOrderLogic {
	return &DetailUserPointsOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailUserPointsOrderLogic) DetailUserPointsOrder(w http.ResponseWriter, r *http.Request) {
	huser.NewPointsOrderHandler(l.svcCtx).Detail(w, r)
}
