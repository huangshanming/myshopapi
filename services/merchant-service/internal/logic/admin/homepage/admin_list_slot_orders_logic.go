package homepage

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListSlotOrdersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListSlotOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListSlotOrdersLogic {
	return &AdminListSlotOrdersLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListSlotOrdersLogic) AdminListSlotOrders(ctx context.Context, req *types.SlotOrderListReq) (resp *types.PageListResp, err error) {
	p, ps := req.Page, req.PageSize
	shopID := req.ShopId
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListSlotOrders(shopID, req.SlotType, req.Status, p, ps)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{Total: total, List: list}, nil

}
