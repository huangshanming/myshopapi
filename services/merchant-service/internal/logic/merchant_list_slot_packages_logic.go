package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListSlotPackagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantListSlotPackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListSlotPackagesLogic {
	return &MerchantListSlotPackagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantListSlotPackagesLogic) MerchantListSlotPackages(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewHomepageSlotHandler(l.svcCtx).MerchantListSlotPackages(w, r)
}
