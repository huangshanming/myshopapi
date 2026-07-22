package homepage

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

func (l *MerchantBuySlotLogic) MerchantBuySlot(ctx context.Context, req *types.BuySlotReq) (resp *types.HomepageOrderResp, err error) {
	shopID := middleware.GetShopID(ctx)
	userID, ok := middleware.GetUserID(ctx)
	if !ok || shopID == 0 {
		return nil, xerr.New(http.StatusUnauthorized, "未授权")
	}
	order, err := biz.NewMerchantLogic(l.svcCtx).BuySlot(shopID, userID, *req)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.HomepageOrderResp{Data: order}, nil
}
