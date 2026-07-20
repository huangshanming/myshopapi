package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	padmin "mymall/services/catalog-service/internal/product/httpapi/admin"
	"mymall/services/catalog-service/internal/svc"
)

type Delete2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelete2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Delete2Logic {
	return &Delete2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Delete2Logic) Delete2(w http.ResponseWriter, r *http.Request) {
	padmin.NewPlatformProductHandler(l.svcCtx).Delete(w, r)
}
