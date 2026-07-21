package public

import (

	"mymall/services/merchant-service/internal/biz"
	"mymall/services/merchant-service/internal/svc"
)

type CouponHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.MerchantLogic
}

func NewCouponHandler(svcCtx *svc.ServiceContext) *CouponHandler {
	return &CouponHandler{
		svcCtx: svcCtx,
		logic:  biz.NewMerchantLogic(svcCtx),
	}
}

type HomepageSlotHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.MerchantLogic
}

func NewHomepageSlotHandler(svcCtx *svc.ServiceContext) *HomepageSlotHandler {
	return &HomepageSlotHandler{
		svcCtx: svcCtx,
		logic:  biz.NewMerchantLogic(svcCtx),
	}
}

type HomepageThemeHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.MerchantLogic
}

func NewHomepageThemeHandler(svcCtx *svc.ServiceContext) *HomepageThemeHandler {
	return &HomepageThemeHandler{
		svcCtx: svcCtx,
		logic:  biz.NewMerchantLogic(svcCtx),
	}
}

type SeckillHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.MerchantLogic
}

func NewSeckillHandler(svcCtx *svc.ServiceContext) *SeckillHandler {
	return &SeckillHandler{
		svcCtx: svcCtx,
		logic:  biz.NewMerchantLogic(svcCtx),
	}
}

type ShopHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.MerchantLogic
}

func NewShopHandler(svcCtx *svc.ServiceContext) *ShopHandler {
	return &ShopHandler{
		svcCtx: svcCtx,
		logic:  biz.NewMerchantLogic(svcCtx),
	}
}
