package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	ppublic "mymall/services/catalog-service/internal/product/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type GetCategoryDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryDetailLogic {
	return &GetCategoryDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryDetailLogic) GetCategoryDetail(w http.ResponseWriter, r *http.Request) {
	ppublic.NewCatalogHandler(l.svcCtx).GetCategoryDetail(w, r)
}
