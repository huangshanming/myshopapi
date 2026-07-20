package handler

import (
	"context"

	"mymall/services/merchant-service/internal/logic"
	"mymall/services/merchant-service/internal/svc"
)

type merchantDeps struct {
	svcCtx *svc.ServiceContext
	logic  *logic.MerchantLogic
}

func newMerchantDeps(svcCtx *svc.ServiceContext) merchantDeps {
	return merchantDeps{
		svcCtx: svcCtx,
		logic:  logic.NewMerchantLogic(context.Background(), svcCtx),
	}
}

type CouponAdminHandler struct {
	merchantDeps
}

func NewCouponAdminHandler(svcCtx *svc.ServiceContext) *CouponAdminHandler {
	return &CouponAdminHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type CouponInternalHandler struct {
	merchantDeps
}

func NewCouponInternalHandler(svcCtx *svc.ServiceContext) *CouponInternalHandler {
	return &CouponInternalHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type CouponMerchantHandler struct {
	merchantDeps
}

func NewCouponMerchantHandler(svcCtx *svc.ServiceContext) *CouponMerchantHandler {
	return &CouponMerchantHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type CouponPublicHandler struct {
	merchantDeps
}

func NewCouponPublicHandler(svcCtx *svc.ServiceContext) *CouponPublicHandler {
	return &CouponPublicHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type CouponUserHandler struct {
	merchantDeps
}

func NewCouponUserHandler(svcCtx *svc.ServiceContext) *CouponUserHandler {
	return &CouponUserHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type HomepageSlotAdminHandler struct {
	merchantDeps
}

func NewHomepageSlotAdminHandler(svcCtx *svc.ServiceContext) *HomepageSlotAdminHandler {
	return &HomepageSlotAdminHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type HomepageSlotMerchantHandler struct {
	merchantDeps
}

func NewHomepageSlotMerchantHandler(svcCtx *svc.ServiceContext) *HomepageSlotMerchantHandler {
	return &HomepageSlotMerchantHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type HomepageSlotPublicHandler struct {
	merchantDeps
}

func NewHomepageSlotPublicHandler(svcCtx *svc.ServiceContext) *HomepageSlotPublicHandler {
	return &HomepageSlotPublicHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type HomepageThemeAdminHandler struct {
	merchantDeps
}

func NewHomepageThemeAdminHandler(svcCtx *svc.ServiceContext) *HomepageThemeAdminHandler {
	return &HomepageThemeAdminHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type HomepageThemeMerchantHandler struct {
	merchantDeps
}

func NewHomepageThemeMerchantHandler(svcCtx *svc.ServiceContext) *HomepageThemeMerchantHandler {
	return &HomepageThemeMerchantHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type HomepageThemePublicHandler struct {
	merchantDeps
}

func NewHomepageThemePublicHandler(svcCtx *svc.ServiceContext) *HomepageThemePublicHandler {
	return &HomepageThemePublicHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type SeckillAdminHandler struct {
	merchantDeps
}

func NewSeckillAdminHandler(svcCtx *svc.ServiceContext) *SeckillAdminHandler {
	return &SeckillAdminHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type SeckillInternalHandler struct {
	merchantDeps
}

func NewSeckillInternalHandler(svcCtx *svc.ServiceContext) *SeckillInternalHandler {
	return &SeckillInternalHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type SeckillMerchantHandler struct {
	merchantDeps
}

func NewSeckillMerchantHandler(svcCtx *svc.ServiceContext) *SeckillMerchantHandler {
	return &SeckillMerchantHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type SeckillPublicHandler struct {
	merchantDeps
}

func NewSeckillPublicHandler(svcCtx *svc.ServiceContext) *SeckillPublicHandler {
	return &SeckillPublicHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type ShopAdminHandler struct {
	merchantDeps
}

func NewShopAdminHandler(svcCtx *svc.ServiceContext) *ShopAdminHandler {
	return &ShopAdminHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type ShopMerchantHandler struct {
	merchantDeps
}

func NewShopMerchantHandler(svcCtx *svc.ServiceContext) *ShopMerchantHandler {
	return &ShopMerchantHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type ShopPublicHandler struct {
	merchantDeps
}

func NewShopPublicHandler(svcCtx *svc.ServiceContext) *ShopPublicHandler {
	return &ShopPublicHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type WalletAdminHandler struct {
	merchantDeps
}

func NewWalletAdminHandler(svcCtx *svc.ServiceContext) *WalletAdminHandler {
	return &WalletAdminHandler{merchantDeps: newMerchantDeps(svcCtx)}
}

type WalletMerchantHandler struct {
	merchantDeps
}

func NewWalletMerchantHandler(svcCtx *svc.ServiceContext) *WalletMerchantHandler {
	return &WalletMerchantHandler{merchantDeps: newMerchantDeps(svcCtx)}
}
