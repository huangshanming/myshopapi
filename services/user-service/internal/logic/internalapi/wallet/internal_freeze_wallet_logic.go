package wallet

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hinternal "mymall/services/user-service/internal/httpapi/internalapi"
	"mymall/services/user-service/internal/svc"
)

type InternalFreezeWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalFreezeWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalFreezeWalletLogic {
	return &InternalFreezeWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InternalFreezeWalletLogic) InternalFreezeWallet(w http.ResponseWriter, r *http.Request) {
	hinternal.NewWalletHandler(l.svcCtx).Freeze(w, r)
}
