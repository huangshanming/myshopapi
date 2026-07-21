package coupon

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
	hadmin "mymall/services/merchant-service/internal/app/admin"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCreateCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCreateCouponLogic(svcCtx *svc.ServiceContext) *AdminCreateCouponLogic {
	return &AdminCreateCouponLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateCouponLogic) AdminCreateCoupon(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/coupons", nil, nil, req, hadmin.NewCouponHandler(l.svcCtx).AdminCreateCoupon)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
