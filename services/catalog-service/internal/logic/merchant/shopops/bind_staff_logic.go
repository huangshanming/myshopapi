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

type BindStaffLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewBindStaffLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindStaffLogic {
	return &BindStaffLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *BindStaffLogic) BindStaff(ctx context.Context, req *types.ShopStaffReq) (resp *types.BindStaffResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	shopID, uid, ok := shopUser(ctx)
	if !ok || !l.svcCtx.ShopRBAC.IsOwner(ctx, shopID, uid) {
		return nil, xerr.New(http.StatusForbidden, "仅店主可操作")
	}
	if req.Mobile == "" {
		return nil, xerr.New(http.StatusBadRequest, "请填写手机号")
	}
	if req.RoleID == 0 {
		return nil, xerr.New(http.StatusBadRequest, "请选择角色")
	}
	mode := req.Mode
	if mode == "" {
		mode = "bind"
	}
	var userID uint64
	switch mode {
	case "create":
		userID, err = l.svcCtx.ShopRBAC.CreateStaffUser(ctx, req.Mobile, req.Password, req.Nickname)
		if err != nil {
			return nil, xerr.New(http.StatusBadRequest, err.Error())
		}
	default:
		userID, err = l.svcCtx.ShopRBAC.FindUserIDByMobile(ctx, req.Mobile)
		if err != nil {
			return nil, xerr.New(http.StatusBadRequest, err.Error())
		}
	}
	if err = l.svcCtx.ShopRBAC.BindStaff(ctx, shopID, userID, req.RoleID); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	msg := "已绑定"
	if mode == "create" {
		msg = "已创建账号并绑定店铺"
	}
	return &types.BindStaffResp{UserId: userID, Msg: msg}, nil

}
