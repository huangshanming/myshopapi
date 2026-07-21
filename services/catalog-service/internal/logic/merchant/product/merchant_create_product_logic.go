package product

import (
	"context"
	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantCreateProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCreateProductLogic {
	return &MerchantCreateProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCreateProductLogic) MerchantCreateProduct(ctx context.Context, req *types.MerchantProductSaveReq) (resp *types.AnyResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, uid, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	if middleware.GetUserRole(ctx) == jwt.RoleMerchantOwner {
		_ = l.svcCtx.ShopRBAC.EnsureOwnerRole(ctx, shopID, uid)
	}
	if !l.svcCtx.ShopRBAC.HasPerm(ctx, shopID, uid, "product:add") && !l.svcCtx.ShopRBAC.HasPerm(ctx, shopID, uid, "product:edit") {
		return nil, xerr.New(http.StatusForbidden, "无权限: product:add")
	}
	p, err := plogic.NewProductAdminLogic(l.svcCtx).Save(ctx, shopID, uid, 0, req.ToProduct())
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: p}, nil
}
