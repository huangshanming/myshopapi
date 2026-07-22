package order

import (
	"context"
	"net/http"

	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/order-service/internal/biz"
	"mymall/services/order-service/internal/svc"
	"mymall/services/order-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponPreviewLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCouponPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponPreviewLogic {
	return &CouponPreviewLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *CouponPreviewLogic) CouponPreview(ctx context.Context, req *types.CouponPreviewReq) (*types.CouponPreviewResp, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	data, err := biz.NewOrderLogic(l.svcCtx).CouponPreview(ctx, userID, req.Items, req.UserCouponID)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.CouponPreviewResp{
		GoodsAmount:      data.GoodsAmount,
		DiscountAmount:   data.DiscountAmount,
		PayAmount:        data.PayAmount,
		BestUserCouponID: data.BestUserCouponID,
		Available:        data.Available,
		Unavailable:      data.Unavailable,
	}, nil
}
