package internalapi

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hinternal "mymall/services/merchant-service/internal/httpapi/internalapi"
	"mymall/services/merchant-service/internal/svc"
)

type SeckillConsumeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSeckillConsumeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillConsumeLogic {
	return &SeckillConsumeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SeckillConsumeLogic) SeckillConsume(w http.ResponseWriter, r *http.Request) {
	hinternal.NewSeckillHandler(l.svcCtx).SeckillConsume(w, r)
}
