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

type MerchantGrantCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantGrantCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantGrantCouponLogic {
	return &MerchantGrantCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantGrantCouponLogic) MerchantGrantCoupon(ctx context.Context, req *types.GrantCouponReq) (resp *types.CouponGrantResp, err error) {
	shopID := middleware.GetShopID(ctx)
	userID, _ := middleware.GetUserID(ctx)
	g, err := biz.NewMerchantLogic(l.svcCtx).GrantCoupon(userID, req.CouponID, req.UserIDs, shopID, false)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.CouponGrantResp{Data: g}, nil
}
