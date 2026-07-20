package theme

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"
)

type AdminCreateThemePackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCreateThemePackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateThemePackageLogic {
	return &AdminCreateThemePackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateThemePackageLogic) AdminCreateThemePackage(w http.ResponseWriter, r *http.Request) {
	hadmin.NewHomepageThemeHandler(l.svcCtx).AdminCreateThemePackage(w, r)
}
