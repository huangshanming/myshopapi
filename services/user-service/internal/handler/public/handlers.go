package public

import (
	"context"

	"mymall/services/user-service/internal/logic"
	"mymall/services/user-service/internal/svc"
)

type RegionHandler struct {
	svcCtx *svc.ServiceContext
	logic  *logic.RegionLogic
}

func NewRegionHandler(svcCtx *svc.ServiceContext) *RegionHandler {
	return &RegionHandler{
		svcCtx: svcCtx,
		logic:  logic.NewRegionLogic(context.Background(), svcCtx),
	}
}
