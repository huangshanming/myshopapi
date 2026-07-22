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

type ClaimCouponLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewClaimCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimCouponLogic {
	return &ClaimCouponLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ClaimCouponLogic) ClaimCoupon(ctx context.Context, req *types.ClaimCouponBodyReq) (resp *types.UserCouponResp, err error) {
	userID, _ := middleware.GetUserID(ctx)
	if userID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "请先登录")
	}
	id := req.Id
	uc, err := biz.NewMerchantLogic(l.svcCtx).ClaimCoupon(userID, id, req.Source)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.UserCouponResp{Data: uc}, nil
}
