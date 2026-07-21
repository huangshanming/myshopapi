package coupon

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalRedeemCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalRedeemCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalRedeemCouponLogic {
	return &InternalRedeemCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalRedeemCouponLogic) InternalRedeemCoupon(ctx context.Context, req *types.RedeemCouponReq) (resp *types.AnyResp, err error) {
	if err := biz.NewMerchantLogic(l.svcCtx).RedeemCoupon(req.UserCouponID, req.OrderID, req.DiscountAmount); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
