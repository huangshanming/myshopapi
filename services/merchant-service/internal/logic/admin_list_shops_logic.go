package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListShopsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListShopsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListShopsLogic {
	return &AdminListShopsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListShopsLogic) AdminListShops(w http.ResponseWriter, r *http.Request) {
	hadmin.NewShopHandler(l.svcCtx).AdminListShops(w, r)
}
