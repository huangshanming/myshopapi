package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminRejectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminRejectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminRejectLogic {
	return &AdminRejectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminRejectLogic) AdminReject(w http.ResponseWriter, r *http.Request) {
	hadmin.NewShopHandler(l.svcCtx).AdminReject(w, r)
}
