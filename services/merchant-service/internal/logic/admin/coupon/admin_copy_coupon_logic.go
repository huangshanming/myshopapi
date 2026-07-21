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

type AdminCopyCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCopyCouponLogic(svcCtx *svc.ServiceContext) *AdminCopyCouponLogic {
	return &AdminCopyCouponLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminCopyCouponLogic) AdminCopyCoupon(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/admin/coupons/:id/copy", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hadmin.NewCouponHandler(l.svcCtx).AdminCopyCoupon)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
