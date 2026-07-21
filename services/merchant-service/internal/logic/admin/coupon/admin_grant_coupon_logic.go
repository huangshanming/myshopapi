package coupon

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGrantCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminGrantCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGrantCouponLogic {
	return &AdminGrantCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminGrantCouponLogic) AdminGrantCoupon(ctx context.Context, req *types.GrantCouponReq) (resp *types.AnyResp, err error) {
	adminID, _ := middleware.GetUserID(ctx)
	g, err := biz.NewMerchantLogic(l.svcCtx).GrantCoupon(adminID, req.CouponID, req.UserIDs, 0, true)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: g}, nil
}
