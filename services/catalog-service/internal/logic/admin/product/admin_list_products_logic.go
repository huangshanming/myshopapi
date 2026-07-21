package product

import (
	"context"
	"fmt"
	"mymall/pkg/appinput"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/product/repository"
	"net/http"
	"net/url"
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
	in := appinput.CallInput{Query: url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}}

	page, pageSize := in.Page()
	shopID, _ := strconv.ParseUint(in.QueryGet("shop_id"), 10, 64)
	catID, _ := strconv.ParseUint(in.QueryGet("category_id"), 10, 64)
	f := repository.ProductListFilter{
		ShopID: shopID, Name: in.QueryGet("name"), ProductNo: in.QueryGet("product_no"),
		CategoryID: catID, Status: in.QueryGet("status"), ProductType: in.QueryGet("product_type"),
		Page: page, PageSize: pageSize, OrderBy: in.QueryGet("order_by"),
		PlatformScope: true,
	}
	if s := in.QueryGet("created_from"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			f.CreatedFrom = &t
		}
	}
	if s := in.QueryGet("created_to"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			f.CreatedTo = &end
		}
	}
	if s := in.QueryGet("publish_from"); s != "" {
		if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			f.PublishFrom = &t
		}
	}
	if s := in.QueryGet("publish_to"); s != "" {
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
