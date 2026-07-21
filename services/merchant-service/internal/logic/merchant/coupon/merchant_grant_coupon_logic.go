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

func (l *MerchantGrantCouponLogic) MerchantGrantCoupon(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	shopID := middleware.GetShopID(ctx)
	userID, _ := middleware.GetUserID(ctx)
	var body struct {
		CouponID uint64   `json:"coupon_id"`
		UserIDs  []uint64 `json:"user_ids"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	g, err := biz.NewMerchantLogic(l.svcCtx).GrantCoupon(userID, body.CouponID, body.UserIDs, shopID, false)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: g}, nil
}
