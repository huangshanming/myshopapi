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

type InternalReturnCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalReturnCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalReturnCouponLogic {
	return &InternalReturnCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalReturnCouponLogic) InternalReturnCoupon(ctx context.Context, req *types.ReturnCouponReq) (resp *types.AnyResp, err error) {
	if err := biz.NewMerchantLogic(l.svcCtx).ReturnCoupon(req.UserCouponID, req.OrderID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
