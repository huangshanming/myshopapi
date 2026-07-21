package article

import (
	"context"
	"io"
	"mymall/pkg/jwt"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	clogic "mymall/services/catalog-service/internal/content/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantUploadArticleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantUploadArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantUploadArticleLogic {
	return &MerchantUploadArticleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantUploadArticleLogic) MerchantUploadArticle(ctx context.Context, r *http.Request) (resp *types.AnyResp, err error) {
	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	if r == nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少上传请求")
	}

	shopID, uid, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	if middleware.GetUserRole(ctx) == jwt.RoleMerchantOwner {
		_ = l.svcCtx.ShopRBAC.EnsureOwnerRole(ctx, shopID, uid)
	}
	if !l.svcCtx.ShopRBAC.HasPerm(ctx, shopID, uid, "article:edit") &&
		!l.svcCtx.ShopRBAC.HasPerm(ctx, shopID, uid, "article:add") {
		return nil, xerr.New(http.StatusForbidden, "无权限: article:edit")
	}
	file, hdr, err := r.FormFile("file")
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少文件")
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "读取失败")
	}
	url, err := clogic.NewArticleLogic(l.svcCtx).SaveUpload(shopID, hdr.Filename, data)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: map[string]string{"url": url}}, nil
}
