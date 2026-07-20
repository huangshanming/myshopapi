package logic

import (
	"net/http"

	"context"

	hpublic "mymall/services/user-service/internal/httpapi/public"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(w http.ResponseWriter, r *http.Request) {
	hpublic.NewRegionHandler(l.svcCtx).List(w, r)
}
