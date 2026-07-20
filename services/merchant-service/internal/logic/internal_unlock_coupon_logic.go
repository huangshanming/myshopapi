package logic

import (
	"net/http"

	"context"

	hinternal "mymall/services/merchant-service/internal/httpapi/internalapi"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalUnlockCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalUnlockCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalUnlockCouponLogic {
	return &InternalUnlockCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InternalUnlockCouponLogic) InternalUnlockCoupon(w http.ResponseWriter, r *http.Request) {
	hinternal.NewCouponHandler(l.svcCtx).InternalUnlockCoupon(w, r)
}
