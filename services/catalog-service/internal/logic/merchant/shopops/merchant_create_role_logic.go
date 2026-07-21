package shopops

import (
	"context"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/shopops/model"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantCreateRoleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantCreateRoleLogic {
	return &MerchantCreateRoleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantCreateRoleLogic) MerchantCreateRole(ctx context.Context, req *types.ShopRoleReq) (resp *types.AnyResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, uid, ok := shopUser(ctx)
	if !ok || !l.svcCtx.ShopRBAC.IsOwner(ctx, shopID, uid) {
		return nil, xerr.New(http.StatusForbidden, "仅店主可操作")
	}
	role := &model.ShopRole{ShopID: shopID, Code: req.Code, Name: req.Name, Remark: req.Remark, Status: 1}
	if role.Code == "" {
		role.Code = "custom"
	}
	if err := l.svcCtx.ShopRBAC.SaveRole(ctx, role, req.MenuIDs); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: role}, nil
}
