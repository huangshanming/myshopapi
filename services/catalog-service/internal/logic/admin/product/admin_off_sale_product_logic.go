package product

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	padmin "mymall/services/catalog-service/internal/product/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type AdminOffSaleProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminOffSaleProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminOffSaleProductLogic {
	return &AdminOffSaleProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminOffSaleProductLogic) AdminOffSaleProduct(w http.ResponseWriter, r *http.Request) {
	padmin.NewPlatformProductHandler(l.svcCtx).OffSale(w, r)
}
