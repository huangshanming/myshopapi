package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateCouponLogic {
	return &AdminUpdateCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateCouponLogic) AdminUpdateCoupon(w http.ResponseWriter, r *http.Request) {
	hadmin.NewCouponHandler(l.svcCtx).AdminUpdateCoupon(w, r)
}
