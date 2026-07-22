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

type MerchantGetProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantGetProductLogic {
	return &MerchantGetProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantGetProductLogic) MerchantGetProduct(ctx context.Context, req *types.IdPathReq) (resp *types.ProductResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id := req.Id
	data, err := plogic.NewProductAdminLogic(l.svcCtx).Detail(ctx, id, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusNotFound, err.Error())
	}
	return &types.ProductResp{Data: data}, nil
}
