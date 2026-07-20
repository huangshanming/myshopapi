package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateShopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateShopLogic {
	return &AdminUpdateShopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateShopLogic) AdminUpdateShop(w http.ResponseWriter, r *http.Request) {
	hadmin.NewShopHandler(l.svcCtx).AdminUpdateShop(w, r)
}
