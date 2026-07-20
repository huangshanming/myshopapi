package logic

import (
	"net/http"

	"context"

	huser "mymall/services/merchant-service/internal/httpapi/user"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyCouponsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyCouponsLogic {
	return &ListMyCouponsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyCouponsLogic) ListMyCoupons(w http.ResponseWriter, r *http.Request) {
	huser.NewCouponHandler(l.svcCtx).ListMyCoupons(w, r)
}
