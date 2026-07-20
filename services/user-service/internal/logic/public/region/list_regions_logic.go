package region

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hpublic "mymall/services/user-service/internal/httpapi/public"
	"mymall/services/user-service/internal/svc"
)

type ListRegionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRegionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRegionsLogic {
	return &ListRegionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRegionsLogic) ListRegions(w http.ResponseWriter, r *http.Request) {
	hpublic.NewRegionHandler(l.svcCtx).List(w, r)
}
