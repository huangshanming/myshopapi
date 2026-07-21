package homepage

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

type MerchantBuySlotLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantBuySlotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantBuySlotLogic {
	return &MerchantBuySlotLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantBuySlotLogic) MerchantBuySlot(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	shopID := middleware.GetShopID(ctx)
	userID, ok := middleware.GetUserID(ctx)
	if !ok || shopID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	var body biz.BuySlotReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	order, err := biz.NewMerchantLogic(l.svcCtx).BuySlot(shopID, userID, body)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: order}, nil
}
