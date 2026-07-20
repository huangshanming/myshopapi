package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateSeckillRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateSeckillRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateSeckillRuleLogic {
	return &AdminUpdateSeckillRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateSeckillRuleLogic) AdminUpdateSeckillRule(w http.ResponseWriter, r *http.Request) {
	hadmin.NewSeckillHandler(l.svcCtx).AdminUpdateSeckillRule(w, r)
}
