package region

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	hpublic "mymall/services/user-service/internal/httpapi/public"
	"mymall/services/user-service/internal/svc"
)

type RegionTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegionTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegionTreeLogic {
	return &RegionTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegionTreeLogic) RegionTree(w http.ResponseWriter, r *http.Request) {
	hpublic.NewRegionHandler(l.svcCtx).Tree(w, r)
}
