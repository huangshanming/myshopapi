package logic

import (
	"net/http"

	"context"

	huser "mymall/services/user-service/internal/httpapi/user"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type List2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewList2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *List2Logic {
	return &List2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *List2Logic) List2(w http.ResponseWriter, r *http.Request) {
	huser.NewAddressHandler(l.svcCtx).Create(w, r)
}
