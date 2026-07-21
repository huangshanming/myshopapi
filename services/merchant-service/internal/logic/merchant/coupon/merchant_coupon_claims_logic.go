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

type MerchantCouponClaimsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCouponClaimsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCouponClaimsLogic {
	return &MerchantCouponClaimsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCouponClaimsLogic) MerchantCouponClaims(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	id := req.Id
	page, pageSize := 1, 10
	list, total, err := biz.NewMerchantLogic(l.svcCtx).CouponClaims(id, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.AnyResp{Data: map[string]interface{}{"list": list, "total": total}}, nil

}
