package shopops

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	shopopshandler "mymall/services/catalog-service/internal/shopops/handler"
	"mymall/services/catalog-service/internal/svc"
)

type RoleMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleMenusLogic {
	return &RoleMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleMenusLogic) RoleMenus(w http.ResponseWriter, r *http.Request) {
	shopopshandler.NewShopOpsHandler(l.svcCtx).RoleMenus(w, r)
}
