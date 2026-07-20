package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListSeckillEntriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListSeckillEntriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListSeckillEntriesLogic {
	return &AdminListSeckillEntriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListSeckillEntriesLogic) AdminListSeckillEntries(w http.ResponseWriter, r *http.Request) {
	hadmin.NewSeckillHandler(l.svcCtx).AdminListSeckillEntries(w, r)
}
