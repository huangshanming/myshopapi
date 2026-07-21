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

type MerchantOffCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantOffCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantOffCouponLogic {
	return &MerchantOffCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantOffCouponLogic) MerchantOffCoupon(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	id := req.Id
	shopID := middleware.GetShopID(ctx)
	if err := biz.NewMerchantLogic(l.svcCtx).OffCoupon(id, shopID, false); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
