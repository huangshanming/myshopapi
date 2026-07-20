package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListCouponsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListCouponsLogic {
	return &AdminListCouponsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListCouponsLogic) AdminListCoupons(w http.ResponseWriter, r *http.Request) {
	hadmin.NewCouponHandler(l.svcCtx).AdminListCoupons(w, r)
}
