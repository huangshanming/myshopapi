package logic

import (
	"net/http"

	"context"

	hinternal "mymall/services/merchant-service/internal/httpapi/internalapi"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SeckillRestoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSeckillRestoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillRestoreLogic {
	return &SeckillRestoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SeckillRestoreLogic) SeckillRestore(w http.ResponseWriter, r *http.Request) {
	hinternal.NewSeckillHandler(l.svcCtx).SeckillRestore(w, r)
}
