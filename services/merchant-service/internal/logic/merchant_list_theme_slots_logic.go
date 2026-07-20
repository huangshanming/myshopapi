package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantListThemeSlotsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantListThemeSlotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListThemeSlotsLogic {
	return &MerchantListThemeSlotsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantListThemeSlotsLogic) MerchantListThemeSlots(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewHomepageThemeHandler(l.svcCtx).MerchantListThemeSlots(w, r)
}
