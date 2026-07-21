package shopops

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/shopops/repository"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantAuthMeLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantAuthMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantAuthMeLogic {
	return &MerchantAuthMeLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantAuthMeLogic) MerchantAuthMe(ctx context.Context) (resp *types.AnyResp, err error) {

	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, uid, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	_ = l.svcCtx.ShopRBAC.EnsureOwnerRole(ctx, shopID, uid)
	perms, _ := l.svcCtx.ShopRBAC.ListPerms(ctx, shopID, uid)
	menus, _ := l.svcCtx.ShopRBAC.ListMenusForUser(ctx, shopID, uid)
	return &types.AnyResp{Data: map[string]interface{}{
		"perms": perms, "menus": menus, "menu_tree": repository.BuildShopMenuTree(menus),
		"is_owner": l.svcCtx.ShopRBAC.IsOwner(ctx, shopID, uid),
	}}, nil

}
