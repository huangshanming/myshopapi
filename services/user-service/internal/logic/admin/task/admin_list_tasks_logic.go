package task

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"
)

type AdminListTasksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListTasksLogic {
	return &AdminListTasksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListTasksLogic) AdminListTasks(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewTaskHandler(l.svcCtx).AdminList
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "marketing:task:list"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
