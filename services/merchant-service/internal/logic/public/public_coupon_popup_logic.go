package public

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hpublic "mymall/services/merchant-service/internal/httpapi/public"
	"mymall/services/merchant-service/internal/svc"
)

type PublicCouponPopupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicCouponPopupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicCouponPopupLogic {
	return &PublicCouponPopupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicCouponPopupLogic) PublicCouponPopup(w http.ResponseWriter, r *http.Request) {
	hpublic.NewCouponHandler(l.svcCtx).PublicCouponPopup(w, r)
}
