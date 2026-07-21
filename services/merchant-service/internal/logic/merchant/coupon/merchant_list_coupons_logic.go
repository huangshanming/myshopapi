package coupon

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListCouponsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListCouponsLogic {
	return &MerchantListCouponsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantListCouponsLogic) MerchantListCoupons(ctx context.Context, req *types.CouponListReq) (resp *types.PageListResp, err error) {
	shopID := middleware.GetShopID(ctx)
	page, pageSize := int(req.Page), int(req.PageSize)
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListCoupons("shop", shopID, req.Status, req.Keyword, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: total}, nil

}
