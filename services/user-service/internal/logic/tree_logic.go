package logic

import (
	"net/http"

	"context"

	hpublic "mymall/services/user-service/internal/httpapi/public"
	"mymall/services/user-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type TreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TreeLogic {
	return &TreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TreeLogic) Tree(w http.ResponseWriter, r *http.Request) {
	hpublic.NewRegionHandler(l.svcCtx).Tree(w, r)
}
