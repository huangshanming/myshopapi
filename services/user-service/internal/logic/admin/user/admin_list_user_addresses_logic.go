package user

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	pkgmw "mymall/pkg/middleware"
	hadmin "mymall/services/user-service/internal/httpapi/admin"
	"mymall/services/user-service/internal/svc"
)

type AdminListUserAddressesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListUserAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListUserAddressesLogic {
	return &AdminListUserAddressesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListUserAddressesLogic) AdminListUserAddresses(w http.ResponseWriter, r *http.Request) {
	h := hadmin.NewAddressHandler(l.svcCtx).AdminList
	admin := hadmin.NewAdminHandler(l.svcCtx)
	if code := "system:user:list"; code != "" {
		h = pkgmw.RequirePermission(admin, code)(h)
	}
	h(w, r)
}
