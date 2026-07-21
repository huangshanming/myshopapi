package product

import (
	"context"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/product/repository"
	"net/http"
	"strconv"
	"time"

	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListProductsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminListProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListProductsLogic {
	return &AdminListProductsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdminListProductsLogic) AdminListProducts(ctx context.Context, req *types.PageReq) (resp *types.PageListResp, err error) {
	page, pageSize := req.Page, req.PageSize
	shopID, _ := strconv.ParseUint("" /* was query:shop_id */, 10, 64)
	catID, _ := strconv.ParseUint("" /* was query:category_id */, 10, 64)
	f := repository.ProductListFilter{
		ShopID: shopID, Name: "" /* was query:name */, ProductNo: "" /* was query:product_no */,
		CategoryID: catID, Status: "" /* was query:status */, ProductType: "" /* was query:product_type */,
		Page: page, PageSize: pageSize, OrderBy: "" /* was query:order_by */,
		PlatformScope: true,
	}
	if s := "" /* was query:created_from */; s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			f.CreatedFrom = &t
		}
	}
	if s := "" /* was query:created_to */; s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			f.CreatedTo = &end
		}
	}
	if s := "" /* was query:publish_from */; s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			f.PublishFrom = &t
		}
	}
	if s := "" /* was query:publish_to */; s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			f.PublishTo = &end
		}
	}
	data, err := plogic.NewPlatformProductLogic(l.svcCtx).List(ctx, f)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: data}, nil
}
