package product

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	ptypes "mymall/services/catalog-service/internal/product/types"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantSetProductStatusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantSetProductStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantSetProductStatusLogic {
	return &MerchantSetProductStatusLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantSetProductStatusLogic) MerchantSetProductStatus(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}, Body: req}

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

	if err := requirePerm(ctx, "product:status"); err != nil {
		return nil, err
	}
	shopID, uid, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	var body ptypes.SetStatusReq
	if err := appinput.BindBody(in, &body); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "参数错误")
	}
	if err := plogic.NewProductAdminLogic(l.svcCtx).SetStatus(ctx, shopID, uid, id, body.Status); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
