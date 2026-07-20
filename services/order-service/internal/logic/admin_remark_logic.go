package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/order-service/internal/httpapi/admin"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminRemarkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminRemarkLogic {
	return &AdminRemarkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminRemarkLogic) AdminRemark(w http.ResponseWriter, r *http.Request) {
	hadmin.NewOrderHandler(l.svcCtx).AdminRemark(w, r)
}
