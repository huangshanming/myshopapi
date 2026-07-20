package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCouponRedeemsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCouponRedeemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCouponRedeemsLogic {
	return &AdminCouponRedeemsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCouponRedeemsLogic) AdminCouponRedeems(w http.ResponseWriter, r *http.Request) {
	hadmin.NewCouponHandler(l.svcCtx).AdminCouponRedeems(w, r)
}
