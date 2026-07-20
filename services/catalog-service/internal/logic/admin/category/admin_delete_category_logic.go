package category

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	padmin "mymall/services/catalog-service/internal/product/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type AdminDeleteCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteCategoryLogic {
	return &AdminDeleteCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteCategoryLogic) AdminDeleteCategory(w http.ResponseWriter, r *http.Request) {
	padmin.NewCatalogHandler(l.svcCtx).AdminDeleteCategory(w, r)
}
