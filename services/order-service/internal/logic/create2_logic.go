package logic

import (
	"net/http"

	"context"

	huser "mymall/services/order-service/internal/httpapi/user"
	"mymall/services/order-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type Create2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreate2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Create2Logic {
	return &Create2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Create2Logic) Create2(w http.ResponseWriter, r *http.Request) {
	huser.NewReviewHandler(l.svcCtx).Create(w, r)
}
