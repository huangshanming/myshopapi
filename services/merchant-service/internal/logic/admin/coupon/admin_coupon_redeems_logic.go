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

type AdminCouponRedeemsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminCouponRedeemsLogic(svcCtx *svc.ServiceContext) *AdminCouponRedeemsLogic {
	return &AdminCouponRedeemsLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminCouponRedeemsLogic) AdminCouponRedeems(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "GET", "/api/v1/admin/coupons/:id/redeems", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, nil, hadmin.NewCouponHandler(l.svcCtx).AdminCouponRedeems)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
