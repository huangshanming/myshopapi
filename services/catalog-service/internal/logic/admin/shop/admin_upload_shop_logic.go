package shop

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	padmin "mymall/services/catalog-service/internal/product/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type AdminUploadShopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUploadShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUploadShopLogic {
	return &AdminUploadShopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUploadShopLogic) AdminUploadShop(w http.ResponseWriter, r *http.Request) {
	padmin.NewShopUploadHandler().Upload(w, r)
}
