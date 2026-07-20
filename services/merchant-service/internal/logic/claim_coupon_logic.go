package logic

import (
	"net/http"

	"context"

	huser "mymall/services/merchant-service/internal/httpapi/user"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClaimCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClaimCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimCouponLogic {
	return &ClaimCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClaimCouponLogic) ClaimCoupon(w http.ResponseWriter, r *http.Request) {
	huser.NewCouponHandler(l.svcCtx).ClaimCoupon(w, r)
}
