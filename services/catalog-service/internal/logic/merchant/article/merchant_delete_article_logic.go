package article

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"
	"strconv"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantDeleteArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantDeleteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantDeleteArticleLogic {
	return &MerchantDeleteArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantDeleteArticleLogic) MerchantDeleteArticle(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{PathVars: map[string]string{"id": fmt.Sprintf("%d", req.Id)}}

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

	if err := requirePerm(ctx, "article:delete"); err != nil {
		return nil, err
	}
	shopID, _, _ := shopUser(ctx)
	id, _ := strconv.ParseUint(in.Path("id"), 10, 64)
	if err := clogic.NewArticleLogic(l.svcCtx).MerchantDelete(ctx, shopID, id); err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: &types.AnyResp{}}, nil
}
