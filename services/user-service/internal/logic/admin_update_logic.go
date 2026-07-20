package logic

import (
	"net/http"

	"context"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateLogic {
	return &AdminUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateLogic) AdminUpdate(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewTaskHandler(l.svcCtx).AdminUpdate
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "marketing:task:edit"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
