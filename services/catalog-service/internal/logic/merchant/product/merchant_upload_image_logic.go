package product

import (
	"context"
	"io"
	"mymall/pkg/appinput"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantUploadImageLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantUploadImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantUploadImageLogic {
	return &MerchantUploadImageLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantUploadImageLogic) MerchantUploadImage(ctx context.Context, r *http.Request) (resp *types.AnyResp, err error) {
	in := appinput.CallInput{Request: r}

	shopUser := func(ctx context.Context) (shopID, userID uint64, ok bool) {
		shopID = middleware.GetShopID(ctx)
		userID, _ = middleware.GetUserID(ctx)
		return shopID, userID, shopID > 0 && userID > 0
	}

	if in.Request == nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少上传请求")
	}

	shopID, _, ok := shopUser(ctx)
	if !ok {
		return nil, xerr.New(http.StatusForbidden, "缺少店铺上下文")
	}
	if err := in.Request.ParseMultipartForm(6 << 20); err != nil {
		return nil, xerr.New(http.StatusBadRequest, "上传失败")
	}
	file, hdr, err := in.Request.FormFile("file")
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少文件")
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "读取失败")
	}
	url, err := plogic.NewProductAdminLogic(l.svcCtx).SaveUpload(shopID, hdr.Filename, data)
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: map[string]string{"url": url}}, nil
}
