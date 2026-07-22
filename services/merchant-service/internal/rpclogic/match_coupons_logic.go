package rpclogic

import (
	"context"

	merchantv1 "mymall/api/gen/merchant/v1"
	"mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MatchCouponsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMatchCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MatchCouponsLogic {
	return &MatchCouponsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *MatchCouponsLogic) MatchCoupons(in *merchantv1.MatchCouponsRequest) (*merchantv1.MatchCouponsResponse, error) {
	items := make([]types.MatchItem, 0, len(in.GetItems()))
	for _, it := range in.GetItems() {
		items = append(items, types.MatchItem{
			ProductID: it.GetProductId(), CategoryID: it.GetCategoryId(),
			Amount: it.GetAmount(), SeckillEntryID: it.GetSeckillEntryId(),
		})
	}
	mr, err := biz.NewMerchantLogic(l.svcCtx).MatchCoupons(types.MatchCouponsReq{
		UserID: in.GetUserId(), ShopID: in.GetShopId(), UserCouponID: in.GetUserCouponId(), Items: items,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	resp := &merchantv1.MatchCouponsResponse{
		GoodsAmount: mr.GoodsAmount, DiscountAmount: mr.DiscountAmount,
		PayAmount: mr.PayAmount, BestUserCouponId: mr.BestUserCouponID,
	}
	for _, v := range mr.Available {
		resp.Available = append(resp.Available, toCouponView(v))
	}
	for _, v := range mr.Unavailable {
		resp.Unavailable = append(resp.Unavailable, toCouponView(v))
	}
	return resp, nil
}

func toCouponView(v biz.MatchCouponView) *merchantv1.MatchCouponView {
	return &merchantv1.MatchCouponView{
		UserCouponId: v.UserCouponID, CouponId: v.CouponID, Name: v.Name, CouponType: v.CouponType,
		DiscountAmount: v.DiscountAmount, ThresholdAmount: v.ThresholdAmount, ValidEnd: v.ValidEnd,
		Usable: v.Usable, Reason: v.Reason, Best: v.Best,
	}
}
