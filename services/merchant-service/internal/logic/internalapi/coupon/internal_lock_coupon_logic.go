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

func (l *InternalLockCouponLogic) InternalLockCoupon(ctx context.Context, req *types.LockCouponReq) (resp *types.AnyResp, err error) {
	if err := biz.NewMerchantLogic(l.svcCtx).LockCoupon(req.UserCouponID, req.UserID, req.OrderID, req.DiscountAmount); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
