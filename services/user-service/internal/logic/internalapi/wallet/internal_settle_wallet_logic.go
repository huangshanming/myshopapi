package wallet

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hinternal "mymall/services/user-service/internal/httpapi/internalapi"
	"mymall/services/user-service/internal/svc"
)

type InternalSettleWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalSettleWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalSettleWalletLogic {
	return &InternalSettleWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InternalSettleWalletLogic) InternalSettleWallet(w http.ResponseWriter, r *http.Request) {
	hinternal.NewWalletHandler(l.svcCtx).Settle(w, r)
}
