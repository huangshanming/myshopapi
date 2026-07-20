package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	shopopshandler "mymall/services/catalog-service/internal/shopops/handler"
	"mymall/services/catalog-service/internal/svc"
)

type ListRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRolesLogic {
	return &ListRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRolesLogic) ListRoles(w http.ResponseWriter, r *http.Request) {
	shopopshandler.NewShopOpsHandler(l.svcCtx).ListRoles(w, r)
}
