package theme

import (
	"context"
	"mymall/pkg/appinput"
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

func (l *MerchantBuyThemeLogic) MerchantBuyTheme(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	shopID := middleware.GetShopID(ctx)
	userID, _ := middleware.GetUserID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺")
	}
	var body biz.ThemeBuyReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	o, err := biz.NewMerchantLogic(l.svcCtx).BuyTheme(shopID, userID, body)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: o}, nil
}
