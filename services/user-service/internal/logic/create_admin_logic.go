package logic

import (
	"net/http"

	"context"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminLogic {
	return &CreateAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAdminLogic) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewAdminHandler(l.svcCtx).CreateAdmin
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "system:admin:add"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
