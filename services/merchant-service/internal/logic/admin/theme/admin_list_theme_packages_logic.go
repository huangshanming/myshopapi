package theme

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"
)

type AdminListThemePackagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListThemePackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListThemePackagesLogic {
	return &AdminListThemePackagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListThemePackagesLogic) AdminListThemePackages(w http.ResponseWriter, r *http.Request) {
	hadmin.NewHomepageThemeHandler(l.svcCtx).AdminListThemePackages(w, r)
}
