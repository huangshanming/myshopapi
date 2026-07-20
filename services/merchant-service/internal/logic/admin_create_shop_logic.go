package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCreateShopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCreateShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateShopLogic {
	return &AdminCreateShopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateShopLogic) AdminCreateShop(w http.ResponseWriter, r *http.Request) {
	hadmin.NewShopHandler(l.svcCtx).AdminCreateShop(w, r)
}
