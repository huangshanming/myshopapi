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

type MerchantUpdateCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantUpdateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantUpdateCouponLogic {
	return &MerchantUpdateCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantUpdateCouponLogic) MerchantUpdateCoupon(ctx context.Context, req *types.CouponUpdateBodyReq) (resp *types.AnyResp, err error) {
	id := req.Id
	shopID := middleware.GetShopID(ctx)
	if err := biz.NewMerchantLogic(l.svcCtx).UpdateCoupon(id, shopID, false, req.ToCouponSaveReq()); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
