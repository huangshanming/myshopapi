package logic

import (
	"net/http"

	"context"

	hinternal "mymall/services/user-service/internal/httpapi/internalapi"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FreezeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFreezeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FreezeLogic {
	return &FreezeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FreezeLogic) Freeze(w http.ResponseWriter, r *http.Request) {
	hinternal.NewWalletHandler(l.svcCtx).Freeze(w, r)
}
