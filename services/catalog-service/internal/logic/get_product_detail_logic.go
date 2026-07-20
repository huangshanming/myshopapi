package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	ppublic "mymall/services/catalog-service/internal/product/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type GetProductDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductDetailLogic {
	return &GetProductDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductDetailLogic) GetProductDetail(w http.ResponseWriter, r *http.Request) {
	ppublic.NewCatalogHandler(l.svcCtx).GetProductDetail(w, r)
}
