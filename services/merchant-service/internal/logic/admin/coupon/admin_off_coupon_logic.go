package coupon

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminOffCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminOffCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminOffCouponLogic {
	return &AdminOffCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminOffCouponLogic) AdminOffCoupon(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	id := req.Id
	if err := biz.NewMerchantLogic(l.svcCtx).OffCoupon(id, 0, true); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{}, nil
}
