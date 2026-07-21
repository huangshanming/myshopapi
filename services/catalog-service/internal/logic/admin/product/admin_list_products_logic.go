package product

import (
	"context"
	"mymall/pkg/xerr"
	plogic "mymall/services/catalog-service/internal/product/logic"
	"mymall/services/catalog-service/internal/product/repository"
	"net/http"
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
	return &AdminListProductsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *AdminListProductsLogic) AdminListProducts(ctx context.Context, req *types.AdminProductListReq) (resp *types.PageListResp, err error) {
	f := repository.ProductListFilter{
		ShopID: req.ShopId, Name: req.Name, ProductNo: req.ProductNo,
		CategoryID: req.CategoryId, Status: req.Status, ProductType: req.ProductType,
		Page: req.Page, PageSize: req.PageSize, OrderBy: req.OrderBy,
		PlatformScope: true,
	}
	if req.CreatedFrom != "" {
		if t, err := time.ParseInLocation("2006-01-02", req.CreatedFrom, time.Local); err == nil {
			f.CreatedFrom = &t
		}
	}
	if req.CreatedTo != "" {
		if t, err := time.ParseInLocation("2006-01-02", req.CreatedTo, time.Local); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			f.CreatedTo = &end
		}
	}
	if req.PublishFrom != "" {
		if t, err := time.ParseInLocation("2006-01-02", req.PublishFrom, time.Local); err == nil {
			f.PublishFrom = &t
		}
	}
	if req.PublishTo != "" {
		if t, err := time.ParseInLocation("2006-01-02", req.PublishTo, time.Local); err == nil {
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
