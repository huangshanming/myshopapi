package banner

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	cadmin "mymall/services/catalog-service/internal/content/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type AdminListBannersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListBannersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListBannersLogic {
	return &AdminListBannersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListBannersLogic) AdminListBanners(w http.ResponseWriter, r *http.Request) {
	cadmin.NewArticleHandler(l.svcCtx).ListBanners(w, r)
}
