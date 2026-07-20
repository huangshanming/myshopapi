package order

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"
)

type CouponPreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCouponPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponPreviewLogic {
	return &CouponPreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CouponPreviewLogic) CouponPreview(w http.ResponseWriter, r *http.Request) {
	huser.NewOrderHandler(l.svcCtx).CouponPreview(w, r)
}
