package product

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	hadmin "mymall/services/catalog-service/internal/product/app/admin"
	"mymall/services/catalog-service/internal/svc"
	"mymall/services/catalog-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminOffSaleProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdminOffSaleProductLogic(svcCtx *svc.ServiceContext) *AdminOffSaleProductLogic {
	return &AdminOffSaleProductLogic{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}
}

func (l *AdminOffSaleProductLogic) AdminOffSaleProduct(ctx context.Context, req *types.IdPathReq) (resp *types.AnyResp, err error) {
	_ = fmt.Sprintf
	_ = url.Values{}
	raw, err := httpinvoke.Run(ctx, "PUT", "/api/v1/admin/products/:id/off_sale", map[string]string{"id": fmt.Sprintf("%d", req.Id)}, nil, req, hadmin.NewPlatformProductHandler(l.svcCtx).OffSale)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := httpinvoke.Decode(raw, &data); err != nil {
		return nil, err
	}
	return &types.AnyResp{Data: data}, nil
}
