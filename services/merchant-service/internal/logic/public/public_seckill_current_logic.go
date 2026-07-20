package public

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hpublic "mymall/services/merchant-service/internal/httpapi/public"
	"mymall/services/merchant-service/internal/svc"
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
