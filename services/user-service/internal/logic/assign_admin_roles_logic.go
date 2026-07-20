package logic

import (
	"net/http"

	"context"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignAdminRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssignAdminRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignAdminRolesLogic {
	return &AssignAdminRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignAdminRolesLogic) AssignAdminRoles(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewAdminHandler(l.svcCtx).AssignAdminRoles
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "system:admin:assign"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
