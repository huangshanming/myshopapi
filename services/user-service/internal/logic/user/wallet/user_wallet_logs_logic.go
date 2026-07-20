package wallet

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"
)

type UserWalletLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserWalletLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserWalletLogsLogic {
	return &UserWalletLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserWalletLogsLogic) UserWalletLogs(w http.ResponseWriter, r *http.Request) {
	huser.NewWalletHandler(l.svcCtx).UserWalletLogs(w, r)
}
