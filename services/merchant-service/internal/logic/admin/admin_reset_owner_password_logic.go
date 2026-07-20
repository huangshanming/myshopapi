package admin

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"
)

type AdminResetOwnerPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminResetOwnerPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminResetOwnerPasswordLogic {
	return &AdminResetOwnerPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminResetOwnerPasswordLogic) AdminResetOwnerPassword(w http.ResponseWriter, r *http.Request) {
	hadmin.NewShopHandler(l.svcCtx).AdminResetOwnerPassword(w, r)
}
