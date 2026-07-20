package health

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"mymall/services/catalog-service/internal/svc"
)

type ReadyzLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadyzLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadyzLogic {
	return &ReadyzLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadyzLogic) Readyz(w http.ResponseWriter, r *http.Request) {
	l.svcCtx.Health.ReadyHandler()(w, r)
}
