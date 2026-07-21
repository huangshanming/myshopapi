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

type MerchantListThemeOrdersLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantListThemeOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantListThemeOrdersLogic {
	return &MerchantListThemeOrdersLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantListThemeOrdersLogic) MerchantListThemeOrders(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	shopID := middleware.GetShopID(ctx)
	if shopID == 0 {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺")
	}
	page, pageSize := int(req.Page), int(req.PageSize)
	list, total, err := biz.NewMerchantLogic(l.svcCtx).ListThemeOrders(shopID, 0, page, pageSize)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list, Total: total}, nil

}
