package logic

import (
	"net/http"

	"context"

	hpublic "mymall/services/merchant-service/internal/httpapi/public"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicSeckillListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicSeckillListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicSeckillListLogic {
	return &PublicSeckillListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicSeckillListLogic) PublicSeckillList(w http.ResponseWriter, r *http.Request) {
	hpublic.NewSeckillHandler(l.svcCtx).PublicSeckillList(w, r)
}
