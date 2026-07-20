package seckill

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"
)

type AdminListSeckillSessionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListSeckillSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListSeckillSessionsLogic {
	return &AdminListSeckillSessionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListSeckillSessionsLogic) AdminListSeckillSessions(w http.ResponseWriter, r *http.Request) {
	hadmin.NewSeckillHandler(l.svcCtx).AdminListSeckillSessions(w, r)
}
