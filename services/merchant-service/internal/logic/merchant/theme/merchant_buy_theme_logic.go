package theme

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantBuyThemeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantBuyThemeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantBuyThemeLogic {
	return &MerchantBuyThemeLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantBuyThemeLogic) MerchantBuyTheme(ctx context.Context, req *types.ThemeBuyReq) (resp *types.AnyResp, err error) {
	shopID := middleware.GetShopID(ctx)
	userID, _ := middleware.GetUserID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺")
	}
	o, err := biz.NewMerchantLogic(l.svcCtx).BuyTheme(shopID, userID, *req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: o}, nil
}
