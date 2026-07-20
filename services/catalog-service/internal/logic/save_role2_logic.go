package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	shopopshandler "mymall/services/catalog-service/internal/shopops/handler"
	"mymall/services/catalog-service/internal/svc"
)

type SaveRole2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveRole2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveRole2Logic {
	return &SaveRole2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveRole2Logic) SaveRole2(w http.ResponseWriter, r *http.Request) {
	shopopshandler.NewShopOpsHandler(l.svcCtx).SaveRole(w, r)
}
