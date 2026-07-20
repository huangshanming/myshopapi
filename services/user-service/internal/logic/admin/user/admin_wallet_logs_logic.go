package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"
)

type AdminWalletLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminWalletLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminWalletLogsLogic {
	return &AdminWalletLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminWalletLogsLogic) AdminWalletLogs(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewWalletHandler(l.svcCtx).AdminWalletLogs
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "system:user:wallet"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
