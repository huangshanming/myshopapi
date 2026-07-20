package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type ScheduleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScheduleLogic {
	return &ScheduleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScheduleLogic) Schedule(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).Schedule(w, r)
}
