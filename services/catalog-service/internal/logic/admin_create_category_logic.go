package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	padmin "mymall/services/catalog-service/internal/product/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type AdminCreateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateCategoryLogic {
	return &AdminCreateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateCategoryLogic) AdminCreateCategory(w http.ResponseWriter, r *http.Request) {
	padmin.NewCatalogHandler(l.svcCtx).AdminCreateCategory(w, r)
}
