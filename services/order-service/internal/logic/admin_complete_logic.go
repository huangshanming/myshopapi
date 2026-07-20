package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/order-service/internal/httpapi/admin"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCompleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCompleteLogic {
	return &AdminCompleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCompleteLogic) AdminComplete(w http.ResponseWriter, r *http.Request) {
	hadmin.NewOrderHandler(l.svcCtx).AdminComplete(w, r)
}
