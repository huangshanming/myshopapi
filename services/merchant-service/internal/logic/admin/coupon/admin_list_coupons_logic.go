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

type AdminListCouponsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListCouponsLogic {
	return &AdminListCouponsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListCouponsLogic) AdminListCoupons(ctx context.Context, req *types.CouponListReq) (resp *types.PageListResp, err error) {
	page, pageSize := int(req.Page), int(req.PageSize)
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListCoupons("platform", 0, req.Status, req.Keyword, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: total}, nil

}
