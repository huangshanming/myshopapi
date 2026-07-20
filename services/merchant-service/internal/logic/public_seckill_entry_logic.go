package logic

import (
	"net/http"

	"context"

	hpublic "mymall/services/merchant-service/internal/httpapi/public"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
