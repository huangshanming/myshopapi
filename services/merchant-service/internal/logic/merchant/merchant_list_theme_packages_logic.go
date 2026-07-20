package merchant

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hmerchant "mymall/services/merchant-service/internal/httpapi/merchant"
	"mymall/services/merchant-service/internal/svc"
)

type MerchantListThemePackagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantListThemePackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListThemePackagesLogic {
	return &MerchantListThemePackagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantListThemePackagesLogic) MerchantListThemePackages(w http.ResponseWriter, r *http.Request) {
	hmerchant.NewHomepageThemeHandler(l.svcCtx).MerchantListThemePackages(w, r)
}
