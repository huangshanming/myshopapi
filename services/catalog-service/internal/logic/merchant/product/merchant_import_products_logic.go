package product

import (
	"context"
	"io"
	"mymall/pkg/middleware"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"net/http"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantImportProductsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMerchantImportProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantImportProductsLogic {
	return &MerchantImportProductsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *MerchantImportProductsLogic) MerchantImportProducts(ctx context.Context, r *http.Request) (resp *types.AnyResp, err error) {
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
	_ = r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, "缺少文件")
	}
	defer file.Close()
	data, _ := io.ReadAll(file)
	res, err := plogic.NewProductAdminLogic(l.svcCtx).ImportCSV(shopID, uid, string(data))
	if err != nil {
		return nil, xerr.New(http.StatusBadRequest, err.Error())
	}
	return &types.AnyResp{Data: res}, nil
}
