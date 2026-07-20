package logic

import (
	"net/http"

	"context"

	hinternal "mymall/services/user-service/internal/httpapi/internalapi"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SettleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSettleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SettleLogic {
	return &SettleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SettleLogic) Settle(w http.ResponseWriter, r *http.Request) {
	hinternal.NewWalletHandler(l.svcCtx).Settle(w, r)
}
