package logic

import (
	"net/http"

	"context"

	hadmin "mymall/services/merchant-service/internal/httpapi/admin"
	"mymall/services/merchant-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGrantThemeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGrantThemeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGrantThemeLogic {
	return &AdminGrantThemeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGrantThemeLogic) AdminGrantTheme(w http.ResponseWriter, r *http.Request) {
	hadmin.NewHomepageThemeHandler(l.svcCtx).AdminGrantTheme(w, r)
}
