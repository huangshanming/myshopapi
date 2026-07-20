package logic

import (
	"net/http"

	"context"

	hpublic "mymall/services/merchant-service/internal/httpapi/public"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicSeckillCurrentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicSeckillCurrentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicSeckillCurrentLogic {
	return &PublicSeckillCurrentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicSeckillCurrentLogic) PublicSeckillCurrent(w http.ResponseWriter, r *http.Request) {
	hpublic.NewSeckillHandler(l.svcCtx).PublicSeckillCurrent(w, r)
}
