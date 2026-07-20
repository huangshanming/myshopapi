package logic

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	puser "mymall/services/catalog-service/internal/product/httpapi/user"
	"mymall/services/catalog-service/internal/svc"
)

type Status2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatus2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Status2Logic {
	return &Status2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Status2Logic) Status2(w http.ResponseWriter, r *http.Request) {
	puser.NewFavoriteHandler(l.svcCtx).Status(w, r)
}
