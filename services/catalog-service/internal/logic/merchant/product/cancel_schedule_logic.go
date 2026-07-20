package product

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pmerchant "mymall/services/catalog-service/internal/product/httpapi/merchant"
	"mymall/services/catalog-service/internal/svc"
)

type CancelScheduleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelScheduleLogic {
	return &CancelScheduleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelScheduleLogic) CancelSchedule(w http.ResponseWriter, r *http.Request) {
	pmerchant.NewProductHandler(l.svcCtx).CancelSchedule(w, r)
}
