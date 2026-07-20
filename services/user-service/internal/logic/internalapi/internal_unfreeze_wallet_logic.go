package internalapi

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hinternal "mymall/services/user-service/internal/httpapi/internalapi"
	"mymall/services/user-service/internal/svc"
)

type InternalUnfreezeWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalUnfreezeWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalUnfreezeWalletLogic {
	return &InternalUnfreezeWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InternalUnfreezeWalletLogic) InternalUnfreezeWallet(w http.ResponseWriter, r *http.Request) {
	hinternal.NewWalletHandler(l.svcCtx).Unfreeze(w, r)
}
