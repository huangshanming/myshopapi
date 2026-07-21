package product

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantExportProductsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantExportProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantExportProductsLogic {
	return &MerchantExportProductsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantExportProductsLogic) MerchantExportProducts(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {

	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	url, err := plogic.NewProductAdminLogic(l.svcCtx).ExportCSV(ctx, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: map[string]string{"url": url}}, nil

}
