package theme

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"
)

type MerchantListThemeOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantListThemeOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListThemeOrdersLogic {
	return &MerchantListThemeOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantListThemeOrdersLogic) MerchantListThemeOrders(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewHomepageThemeHandler(l.svcCtx).MerchantListThemeOrders(w, r)
}
