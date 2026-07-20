package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemarkPointsOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemarkPointsOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemarkPointsOrderLogic {
	return &RemarkPointsOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemarkPointsOrderLogic) RemarkPointsOrder(w http.ResponseWriter, r *http.Request) {
	hadmin.NewPointsOrderHandler(l.svcCtx).Remark(w, r)
}
