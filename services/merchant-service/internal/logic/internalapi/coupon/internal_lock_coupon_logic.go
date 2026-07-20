package coupon

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hinternal "mymall/services/merchant-service/internal/httpapi/internalapi"
	"mymall/services/merchant-service/internal/svc"
)

type InternalLockCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalLockCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalLockCouponLogic {
	return &InternalLockCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InternalLockCouponLogic) InternalLockCoupon(w http.ResponseWriter, r *http.Request) {
	hinternal.NewCouponHandler(l.svcCtx).InternalLockCoupon(w, r)
}
