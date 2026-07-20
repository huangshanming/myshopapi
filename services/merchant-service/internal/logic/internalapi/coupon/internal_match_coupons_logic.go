package coupon

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hinternal "mymall/services/merchant-service/internal/httpapi/internalapi"
	"mymall/services/merchant-service/internal/svc"
)

type InternalMatchCouponsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalMatchCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalMatchCouponsLogic {
	return &InternalMatchCouponsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InternalMatchCouponsLogic) InternalMatchCoupons(w http.ResponseWriter, r *http.Request) {
	hinternal.NewCouponHandler(l.svcCtx).InternalMatchCoupons(w, r)
}
