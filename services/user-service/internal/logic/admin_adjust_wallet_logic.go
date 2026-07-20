package logic

import (
	"net/http"

	"context"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminAdjustWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminAdjustWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminAdjustWalletLogic {
	return &AdminAdjustWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminAdjustWalletLogic) AdminAdjustWallet(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewWalletHandler(l.svcCtx).AdminAdjustWallet
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "system:user:wallet"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
