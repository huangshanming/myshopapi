package staff

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"
)

type GetAdminRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdminRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminRolesLogic {
	return &GetAdminRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdminRolesLogic) GetAdminRoles(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewAdminHandler(l.svcCtx).AssignAdminRoles
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "system:admin:assign"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
