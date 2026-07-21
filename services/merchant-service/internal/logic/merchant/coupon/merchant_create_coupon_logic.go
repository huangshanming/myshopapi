package coupon

import (
	"context"
	"mymall/pkg/appinput"
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

func (l *MerchantCreateCouponLogic) MerchantCreateCoupon(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	shopID := middleware.GetShopID(ctx)
	userID, _ := middleware.GetUserID(ctx)
	var body biz.CouponSaveReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	c, err := biz.NewMerchantLogic(l.svcCtx).MerchantCreateCoupon(shopID, userID, body)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: c}, nil
}
