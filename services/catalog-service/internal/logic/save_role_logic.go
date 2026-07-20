package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	shopopshandler "mymall/services/catalog-service/internal/shopops/handler"
	"mymall/services/catalog-service/internal/svc"
)

type SaveRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveRoleLogic {
	return &SaveRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveRoleLogic) SaveRole(w http.ResponseWriter, r *http.Request) {
	shopopshandler.NewShopOpsHandler(l.svcCtx).SaveRole(w, r)
}
