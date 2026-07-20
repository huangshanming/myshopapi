package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompletePointsOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompletePointsOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompletePointsOrderLogic {
	return &CompletePointsOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompletePointsOrderLogic) CompletePointsOrder(w http.ResponseWriter, r *http.Request) {
	hadmin.NewPointsOrderHandler(l.svcCtx).Complete(w, r)
}
