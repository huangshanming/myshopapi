package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCopyCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCopyCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCopyCouponLogic {
	return &AdminCopyCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCopyCouponLogic) AdminCopyCoupon(w http.ResponseWriter, r *http.Request) {
	hadmin.NewCouponHandler(l.svcCtx).AdminCopyCoupon(w, r)
}
