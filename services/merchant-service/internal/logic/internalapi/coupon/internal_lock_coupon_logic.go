package coupon

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InternalLockCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalLockCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalLockCouponLogic {
	return &InternalLockCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalLockCouponLogic) InternalLockCoupon(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	var body struct {
		UserCouponID   uint64  `json:"user_coupon_id"`
		UserID         uint64  `json:"user_id"`
		OrderID        uint64  `json:"order_id"`
		DiscountAmount float64 `json:"discount_amount"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := biz.NewMerchantLogic(l.svcCtx).LockCoupon(body.UserCouponID, body.UserID, body.OrderID, body.DiscountAmount); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
