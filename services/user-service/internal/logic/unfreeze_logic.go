package logic

import (
	"net/http"

	"context"

	hinternal "mymall/services/user-service/internal/httpapi/internalapi"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnfreezeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnfreezeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfreezeLogic {
	return &UnfreezeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnfreezeLogic) Unfreeze(w http.ResponseWriter, r *http.Request) {
	hinternal.NewWalletHandler(l.svcCtx).Unfreeze(w, r)
}
