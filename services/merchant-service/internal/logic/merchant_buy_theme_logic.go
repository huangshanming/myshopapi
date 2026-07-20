package logic

import (
	"net/http"

	"context"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantBuyThemeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantBuyThemeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantBuyThemeLogic {
	return &MerchantBuyThemeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantBuyThemeLogic) MerchantBuyTheme(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewHomepageThemeHandler(l.svcCtx).MerchantBuyTheme(w, r)
}
