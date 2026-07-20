package seckill

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hpublic "mymall/services/merchant-service/internal/httpapi/public"
	"mymall/services/merchant-service/internal/svc"
)

type PublicSeckillEntryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicSeckillEntryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicSeckillEntryLogic {
	return &PublicSeckillEntryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicSeckillEntryLogic) PublicSeckillEntry(w http.ResponseWriter, r *http.Request) {
	hpublic.NewSeckillHandler(l.svcCtx).PublicSeckillEntry(w, r)
}
