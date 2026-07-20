package logic

import (
	"net/http"

	"context"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
