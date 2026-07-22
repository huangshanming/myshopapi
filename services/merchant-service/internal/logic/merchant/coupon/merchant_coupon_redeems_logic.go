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

type MerchantCouponRedeemsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCouponRedeemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCouponRedeemsLogic {
	return &MerchantCouponRedeemsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCouponRedeemsLogic) MerchantCouponRedeems(ctx context.Context, req *types.IdPageReq) (resp *types.PageListResp, err error) {
	id := req.Id
	page, pageSize := req.Page, req.PageSize
	list, total, err := biz.NewMerchantLogic(l.svcCtx).CouponRedeems(id, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: total}, nil

}
