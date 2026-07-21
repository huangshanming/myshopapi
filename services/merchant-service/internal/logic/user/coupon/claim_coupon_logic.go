package coupon

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"
	huser "mymall/services/merchant-service/internal/app/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClaimCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewClaimCouponLogic(svcCtx *svc.ServiceContext) *ClaimCouponLogic {
	return &ClaimCouponLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *ClaimCouponLogic) ClaimCoupon(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
raw, err := httpinvoke.Run(ctx, "POST", "/api/v1/coupons/:id/claim", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, huser.NewCouponHandler(l.svcCtx).ClaimCoupon)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
