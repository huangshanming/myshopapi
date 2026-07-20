package logic

import (
	"net/http"

	"context"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignRoleMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssignRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRoleMenusLogic {
	return &AssignRoleMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignRoleMenusLogic) AssignRoleMenus(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewAdminHandler(l.svcCtx).AssignRoleMenus
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "system:role:assign"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
