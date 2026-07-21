package shopops

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListStaffLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListStaffLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListStaffLogic {
	return &ListStaffLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListStaffLogic) ListStaff(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {

	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	list, err := l.svcCtx.ShopRBAC.ListStaff(ctx, shopID)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list}, nil
}
