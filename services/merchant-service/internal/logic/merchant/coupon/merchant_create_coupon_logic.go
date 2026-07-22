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

type MerchantCreateCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCreateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCreateCouponLogic {
	return &MerchantCreateCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCreateCouponLogic) MerchantCreateCoupon(ctx context.Context, req *types.CouponSaveReq) (resp *types.CouponResp, err error) {
	shopID := middleware.GetShopID(ctx)
	userID, _ := middleware.GetUserID(ctx)
	c, err := biz.NewMerchantLogic(l.svcCtx).MerchantCreateCoupon(shopID, userID, *req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.CouponResp{Data: c}, nil
}
