package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	ppublic "mymall/services/catalog-service/internal/product/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type GetSalesRankLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSalesRankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSalesRankLogic {
	return &GetSalesRankLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSalesRankLogic) GetSalesRank(w http.ResponseWriter, r *http.Request) {
	ppublic.NewCatalogHandler(l.svcCtx).GetSalesRank(w, r)
}
