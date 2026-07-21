package product

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	ptypes "mymall/services/catalog-service/internal/product/types"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantBatchProductsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantBatchProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantBatchProductsLogic {
	return &MerchantBatchProductsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantBatchProductsLogic) MerchantBatchProducts(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	requirePerm := func(ctx context.Context, code string) error {
		shopID, uid, ok := shopUser(ctx)
		if !ok {
			return xerr.New(http.StatusForbidden, "缺少店铺上下文")
		}
		if middleware.GetUserRole(ctx) == jwt.RoleMerchantOwner {
			_ = l.svcCtx.ShopRBAC.EnsureOwnerRole(ctx, shopID, uid)
		}
		if !l.svcCtx.ShopRBAC.HasPerm(ctx, shopID, uid, code) {
			return xerr.New(http.StatusForbidden, "无权限: "+code)
		}
		return nil
	}

	if err := requirePerm(ctx, "product:batch"); err != nil {
		return nil, err
	}
	shopID, uid, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	var body ptypes.BatchProductReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	job, err := plogic.NewProductAdminLogic(l.svcCtx).Batch(ctx, shopID, uid, body)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: job}, nil
}
