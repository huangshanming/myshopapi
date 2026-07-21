package shopops

import (
	"context"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	"mymall/services/catalog-service/internal/shopops/model"
	sotypes "mymall/services/catalog-service/internal/shopops/types"
	"net/http"
	"strconv"

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

func (l *MerchantCreateRoleLogic) MerchantCreateRole(ctx context.Context, req *types.JSONBody) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Body: req}

	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, uid, ok := shopUser(ctx)
	if !ok || !l.svcCtx.ShopRBAC.IsOwner(ctx, shopID, uid) {
		return nil, xerr.New(http.StatusForbidden, "仅店主可操作")
	}
	var body sotypes.ShopRoleReq
	_ = appinput.BindBody(in, &body)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	role := &model.ShopRole{ID: id, ShopID: shopID, Code: body.Code, Name: body.Name, Remark: body.Remark, Status: 1}
	if role.Code == "" {
		role.Code = "custom"
	}
	if err := l.svcCtx.ShopRBAC.SaveRole(ctx, role, body.MenuIDs); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: role}, nil
}
