package admin

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"
)

type AdminUpdateTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateTaskLogic {
	return &AdminUpdateTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateTaskLogic) AdminUpdateTask(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewTaskHandler(l.svcCtx).AdminUpdate
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "marketing:task:edit"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
