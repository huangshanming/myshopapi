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

type ListMenusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenusLogic {
	return &ListMenusLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListMenusLogic) ListMenus(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {

	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, uid, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	_ = l.svcCtx.ShopRBAC.EnsureShopMenus(ctx)
	_ = l.svcCtx.ShopRBAC.EnsureOwnerRole(ctx, shopID, uid)
	menus, err := l.svcCtx.ShopRBAC.MenuTree(ctx)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: repository.BuildShopMenuTree(menus)}, nil

}
