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

type PublicCouponCenterLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPublicCouponCenterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicCouponCenterLogic {
	return &PublicCouponCenterLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PublicCouponCenterLogic) PublicCouponCenter(ctx context.Context, req *types.ShopIdQueryReq) (resp *types.ListResp, err error) {
	userID, _ := middleware.GetUserID(ctx)
	list, err := biz.NewMerchantLogic(l.svcCtx).ListCenter(userID, req.ShopId)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.ListResp{List: list}, nil
}
