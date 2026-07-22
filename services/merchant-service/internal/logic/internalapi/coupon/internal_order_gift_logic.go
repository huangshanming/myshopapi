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

type InternalOrderGiftLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewInternalOrderGiftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalOrderGiftLogic {
	return &InternalOrderGiftLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *InternalOrderGiftLogic) InternalOrderGift(ctx context.Context, req *types.OrderGiftCouponReq) (resp *types.GrantedCountResp, err error) {
	n, err := biz.NewMerchantLogic(l.svcCtx).OrderGiftCoupons(req.UserID, req.ShopID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.GrantedCountResp{Granted: int64(n)}, nil
}
