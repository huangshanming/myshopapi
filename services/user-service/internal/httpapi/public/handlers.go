package public

import (
	"context"

	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
)

type RegionHandler struct {
	svcCtx *svc.ServiceContext
	logic  *biz.RegionLogic
}

func NewRegionHandler(svcCtx *svc.ServiceContext) *RegionHandler {
	return &RegionHandler{
		svcCtx: svcCtx,
		logic:  biz.NewRegionLogic(context.Background(), svcCtx),
	}
}
