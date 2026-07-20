package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateThemePackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateThemePackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateThemePackageLogic {
	return &AdminUpdateThemePackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateThemePackageLogic) AdminUpdateThemePackage(w http.ResponseWriter, r *http.Request) {
	hadmin.NewHomepageThemeHandler(l.svcCtx).AdminUpdateThemePackage(w, r)
}
