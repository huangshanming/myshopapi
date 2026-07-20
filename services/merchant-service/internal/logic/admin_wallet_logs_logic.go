package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
	hadmin.NewWalletHandler(l.svcCtx).AdminWalletLogs(w, r)
}
