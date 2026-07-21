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

type AdminCreateCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateCouponLogic {
	return &AdminCreateCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateCouponLogic) AdminCreateCoupon(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	adminID, _ := middleware.GetUserID(ctx)
	var body biz.CouponSaveReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	c, err := biz.NewMerchantLogic(l.svcCtx).AdminCreateCoupon(adminID, body)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: c}, nil
}
