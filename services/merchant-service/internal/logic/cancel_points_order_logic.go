package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelPointsOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelPointsOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelPointsOrderLogic {
	return &CancelPointsOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelPointsOrderLogic) CancelPointsOrder(w http.ResponseWriter, r *http.Request) {
	hadmin.NewPointsOrderHandler(l.svcCtx).Cancel(w, r)
}
