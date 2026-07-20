package product

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	padmin "mymall/services/catalog-service/internal/product/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type AdminDeleteProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteProductLogic {
	return &AdminDeleteProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteProductLogic) AdminDeleteProduct(w http.ResponseWriter, r *http.Request) {
	padmin.NewPlatformProductHandler(l.svcCtx).Delete(w, r)
}
