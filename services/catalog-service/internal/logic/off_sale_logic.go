package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	padmin "mymall/services/catalog-service/internal/product/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type OffSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOffSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OffSaleLogic {
	return &OffSaleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OffSaleLogic) OffSale(w http.ResponseWriter, r *http.Request) {
	padmin.NewPlatformProductHandler(l.svcCtx).OffSale(w, r)
}
