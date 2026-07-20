package logic

import (
	"net/http"

	"context"

	hinternal "mymall/services/merchant-service/internal/httpapi/internalapi"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalRedeemCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalRedeemCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalRedeemCouponLogic {
	return &InternalRedeemCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InternalRedeemCouponLogic) InternalRedeemCoupon(w http.ResponseWriter, r *http.Request) {
	hinternal.NewCouponHandler(l.svcCtx).InternalRedeemCoupon(w, r)
}
