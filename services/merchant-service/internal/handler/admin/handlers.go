package admin

import (
	"context"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

type CouponHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.MerchantLogic
}

func NewCouponHandler(svcCtx *svc.ServiceContext) *CouponHandler {
	return &CouponHandler{
		svcCtx: svcCtx,
		logic:  logic.NewMerchantLogic(context.Background(), svcCtx),
	}
}

type HomepageSlotHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.MerchantLogic
}

func NewHomepageSlotHandler(svcCtx *svc.ServiceContext) *HomepageSlotHandler {
	return &HomepageSlotHandler{
		svcCtx: svcCtx,
		logic:  logic.NewMerchantLogic(context.Background(), svcCtx),
	}
}

type HomepageThemeHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.MerchantLogic
}

func NewHomepageThemeHandler(svcCtx *svc.ServiceContext) *HomepageThemeHandler {
	return &HomepageThemeHandler{
		svcCtx: svcCtx,
		logic:  logic.NewMerchantLogic(context.Background(), svcCtx),
	}
}

type PointsOrderHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.PointsOrderLogic
}

func NewPointsOrderHandler(svcCtx *svc.ServiceContext) *PointsOrderHandler {
	return &PointsOrderHandler{
		svcCtx: svcCtx,
		logic:  logic.NewPointsOrderLogic(context.Background(), svcCtx),
	}
}

type PointsProductHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.PointsProductLogic
}

func NewPointsProductHandler(svcCtx *svc.ServiceContext) *PointsProductHandler {
	return &PointsProductHandler{
		svcCtx: svcCtx,
		logic:  logic.NewPointsProductLogic(context.Background(), svcCtx),
	}
}

type SeckillHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.MerchantLogic
}

func NewSeckillHandler(svcCtx *svc.ServiceContext) *SeckillHandler {
	return &SeckillHandler{
		svcCtx: svcCtx,
		logic:  logic.NewMerchantLogic(context.Background(), svcCtx),
	}
}

type ShopHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.MerchantLogic
}

func NewShopHandler(svcCtx *svc.ServiceContext) *ShopHandler {
	return &ShopHandler{
		svcCtx: svcCtx,
		logic:  logic.NewMerchantLogic(context.Background(), svcCtx),
	}
}

type WalletHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.MerchantLogic
}

func NewWalletHandler(svcCtx *svc.ServiceContext) *WalletHandler {
	return &WalletHandler{
		svcCtx: svcCtx,
		logic:  logic.NewMerchantLogic(context.Background(), svcCtx),
	}
}
