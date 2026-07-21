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

type MerchantListSlotOrdersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListSlotOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListSlotOrdersLogic {
	return &MerchantListSlotOrdersLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantListSlotOrdersLogic) MerchantListSlotOrders(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺")
	}
	p, ps := req.Page, req.PageSize
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListSlotOrders(shopID, "" /* was query:slot_type */, "" /* was query:status */, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil

}
