package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGrantCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGrantCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGrantCouponLogic {
	return &AdminGrantCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGrantCouponLogic) AdminGrantCoupon(w http.ResponseWriter, r *http.Request) {
	hadmin.NewCouponHandler(l.svcCtx).AdminGrantCoupon(w, r)
}
