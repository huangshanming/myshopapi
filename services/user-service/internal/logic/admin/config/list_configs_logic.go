package config

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"
)

type ListConfigsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListConfigsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListConfigsLogic {
	return &ListConfigsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListConfigsLogic) ListConfigs(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewAdminHandler(l.svcCtx).SaveConfigs
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "system:config:edit"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
