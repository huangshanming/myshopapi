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

type InternalUnlockCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalUnlockCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalUnlockCouponLogic {
	return &InternalUnlockCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalUnlockCouponLogic) InternalUnlockCoupon(ctx context.Context, req *types.UnlockCouponReq) (resp *types.AnyResp, err error) {
	if err := biz.NewMerchantLogic(l.svcCtx).UnlockCoupon(req.UserCouponID, req.OrderID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
