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

type AdminCreateCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateCouponLogic {
	return &AdminCreateCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateCouponLogic) AdminCreateCoupon(ctx context.Context, req *types.CouponSaveReq) (resp *types.CouponResp, err error) {
	adminID, _ := middleware.GetUserID(ctx)
	c, err := biz.NewMerchantLogic(l.svcCtx).AdminCreateCoupon(adminID, *req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.CouponResp{Data: c}, nil
}
