package theme

import (
	"context"
	"mymall/pkg/xerr"
	"mymall/services/merchant-service/internal/biz"
	"net/http"

	"mymall/services/merchant-service/internal/svc"
	"mymall/services/merchant-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListThemeOrdersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListThemeOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListThemeOrdersLogic {
	return &AdminListThemeOrdersLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListThemeOrdersLogic) AdminListThemeOrders(ctx context.Context, req *types.ThemeOrderListReq) (resp *types.PageListResp, err error) {
	shopID := req.ShopId
	slotID := req.ThemeSlotId
	page, pageSize := int(req.Page), int(req.PageSize)
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListThemeOrders(shopID, slotID, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: total}, nil

}
