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

type AdminOffCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminOffCouponLogic(svcCtx *svc.ServiceContext) *AdminOffCouponLogic {
	return &AdminOffCouponLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminOffCouponLogic) AdminOffCoupon(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/coupons/:id/off", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hadmin.NewCouponHandler(l.svcCtx).AdminOffCoupon)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
