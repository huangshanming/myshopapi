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

type AdminGrantCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminGrantCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGrantCouponLogic {
	return &AdminGrantCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminGrantCouponLogic) AdminGrantCoupon(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	adminID, _ := middleware.GetUserID(ctx)
	var body struct {
		CouponID uint64   `json:"coupon_id"`
		UserIDs  []uint64 `json:"user_ids"`
	}
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	g, err := biz.NewMerchantLogic(l.svcCtx).GrantCoupon(adminID, body.CouponID, body.UserIDs, 0, true)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: g}, nil
}
