package logic

import (
	"net/http"

	"context"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewAdminHandler(l.svcCtx).DeleteRole
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "system:role:delete"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
