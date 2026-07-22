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

type AdminUpdateCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateCouponLogic {
	return &AdminUpdateCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateCouponLogic) AdminUpdateCoupon(ctx context.Context, req *types.CouponUpdateBodyReq) (resp *types.EmptyResp, err error) {
	id := req.Id
	if err := biz.NewMerchantLogic(l.svcCtx).UpdateCoupon(id, 0, true, req.ToCouponSaveReq()); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.EmptyResp{}, nil
}
