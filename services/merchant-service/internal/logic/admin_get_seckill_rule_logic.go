package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetSeckillRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetSeckillRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetSeckillRuleLogic {
	return &AdminGetSeckillRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetSeckillRuleLogic) AdminGetSeckillRule(w http.ResponseWriter, r *http.Request) {
	hadmin.NewSeckillHandler(l.svcCtx).AdminGetSeckillRule(w, r)
}
