package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	ppublic "mymall/services/catalog-service/internal/product/httpapi/public"
	"mymall/services/catalog-service/internal/svc"
)

type GetCategoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryListLogic {
	return &GetCategoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryListLogic) GetCategoryList(w http.ResponseWriter, r *http.Request) {
	ppublic.NewCatalogHandler(l.svcCtx).GetCategoryList(w, r)
}
