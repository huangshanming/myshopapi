package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	padmin "mymall/services/catalog-service/internal/product/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type AdminUpdateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateCategoryLogic {
	return &AdminUpdateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateCategoryLogic) AdminUpdateCategory(w http.ResponseWriter, r *http.Request) {
	padmin.NewCatalogHandler(l.svcCtx).AdminUpdateCategory(w, r)
}
