package logic

import (
	"net/http"

	"context"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminList2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminList2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminList2Logic {
	return &AdminList2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminList2Logic) AdminList2(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewTaskHandler(l.svcCtx).AdminList
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "marketing:task:list"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
