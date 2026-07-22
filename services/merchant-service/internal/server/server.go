package server

import (
	"context"
	"fmt"

	merchantv1 "mymall/api/gen/merchant/v1"
	"mymall/pkg/grpc/interceptor"
	"mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MerchantServer struct {
	merchantv1.UnimplementedMerchantServiceServer
	logic *biz.MerchantLogic
}

func NewMerchantServer(logic *biz.MerchantLogic) *MerchantServer {
	return &MerchantServer{logic: logic}
}

func (s *MerchantServer) ConsumeSeckill(ctx context.Context, req *merchantv1.ConsumeSeckillRequest) (*merchantv1.ConsumeSeckillResponse, error) {
	out, err := s.logic.ConsumeSeckill(req.GetEntryId(), req.GetProductId(), int(req.GetQuantity()))
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	price, _ := out["seckill_price"].(float64)
	return &merchantv1.ConsumeSeckillResponse{
		EntryId:      req.GetEntryId(),
		ProductId:    req.GetProductId(),
		SeckillPrice: price,
		Quantity:     req.GetQuantity(),
	}, nil
}

func (s *MerchantServer) RestoreSeckill(ctx context.Context, req *merchantv1.RestoreSeckillRequest) (*merchantv1.EmptyResponse, error) {
	if err := s.logic.RestoreSeckill(req.GetEntryId(), int(req.GetQuantity())); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &merchantv1.EmptyResponse{}, nil
}

func (s *MerchantServer) MatchCoupons(ctx context.Context, req *merchantv1.MatchCouponsRequest) (*merchantv1.MatchCouponsResponse, error) {
	items := make([]types.MatchItem, 0, len(req.GetItems()))
	for _, it := range req.GetItems() {
		items = append(items, types.MatchItem{
			ProductID:      it.GetProductId(),
			CategoryID:     it.GetCategoryId(),
			Amount:         it.GetAmount(),
			SeckillEntryID: it.GetSeckillEntryId(),
		})
	}
	mr, err := s.logic.MatchCoupons(types.MatchCouponsReq{
		UserID:       req.GetUserId(),
		ShopID:       req.GetShopId(),
		UserCouponID: req.GetUserCouponId(),
		Items:        items,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	resp := &merchantv1.MatchCouponsResponse{
		GoodsAmount:      mr.GoodsAmount,
		DiscountAmount:   mr.DiscountAmount,
		PayAmount:        mr.PayAmount,
		BestUserCouponId: mr.BestUserCouponID,
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
		UserCouponId:    v.UserCouponID,
		CouponId:        v.CouponID,
		Name:            v.Name,
		CouponType:      v.CouponType,
		DiscountAmount:  v.DiscountAmount,
		ThresholdAmount: v.ThresholdAmount,
		ValidEnd:        v.ValidEnd,
		Usable:          v.Usable,
		Reason:          v.Reason,
		Best:            v.Best,
	}
}

func (s *MerchantServer) LockCoupon(ctx context.Context, req *merchantv1.LockCouponRequest) (*merchantv1.EmptyResponse, error) {
	if err := s.logic.LockCoupon(req.GetUserCouponId(), req.GetUserId(), req.GetOrderId(), req.GetDiscountAmount()); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &merchantv1.EmptyResponse{}, nil
}

func (s *MerchantServer) UnlockCoupon(ctx context.Context, req *merchantv1.UnlockCouponRequest) (*merchantv1.EmptyResponse, error) {
	if err := s.logic.UnlockCoupon(req.GetUserCouponId(), req.GetOrderId()); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &merchantv1.EmptyResponse{}, nil
}

func (s *MerchantServer) RedeemCoupon(ctx context.Context, req *merchantv1.RedeemCouponRequest) (*merchantv1.EmptyResponse, error) {
	if err := s.logic.RedeemCoupon(req.GetUserCouponId(), req.GetOrderId(), req.GetDiscountAmount()); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &merchantv1.EmptyResponse{}, nil
}

func (s *MerchantServer) ReturnCoupon(ctx context.Context, req *merchantv1.ReturnCouponRequest) (*merchantv1.EmptyResponse, error) {
	if err := s.logic.ReturnCoupon(req.GetUserCouponId(), req.GetOrderId()); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &merchantv1.EmptyResponse{}, nil
}

func (s *MerchantServer) OrderGiftCoupons(ctx context.Context, req *merchantv1.OrderGiftCouponsRequest) (*merchantv1.OrderGiftCouponsResponse, error) {
	n, err := s.logic.OrderGiftCoupons(req.GetUserId(), req.GetShopId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	return &merchantv1.OrderGiftCouponsResponse{Granted: int32(n)}, nil
}

func StartZRPC(port int, svcCtx *svc.ServiceContext, logger *zap.Logger) *zrpc.RpcServer {
	c := zrpc.RpcServerConf{
		ListenOn: fmt.Sprintf("0.0.0.0:%d", port),
	}
	c.Mode = service.DevMode
	c.Log.Mode = "console"
	c.Log.Encoding = "plain"

	logic := biz.NewMerchantLogic(svcCtx)
	s := zrpc.MustNewServer(c, func(grpcServer *grpc.Server) {
		merchantv1.RegisterMerchantServiceServer(grpcServer, NewMerchantServer(logic))
	})
	s.AddUnaryInterceptors(interceptor.Logging(logger))
	return s
}
