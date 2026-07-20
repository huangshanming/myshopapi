package coupon

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hpublic "mymall/services/merchant-service/internal/httpapi/public"
	"mymall/services/merchant-service/internal/svc"
)

type PublicCouponCenterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicCouponCenterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicCouponCenterLogic {
	return &PublicCouponCenterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicCouponCenterLogic) PublicCouponCenter(w http.ResponseWriter, r *http.Request) {
	hpublic.NewCouponHandler(l.svcCtx).PublicCouponCenter(w, r)
}
