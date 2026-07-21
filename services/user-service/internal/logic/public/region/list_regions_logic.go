package region

import (
	"context"
	"net/http"
	"github.com/zeromicro/go-zero/core/logx"
	"mymall/pkg/xerr"
	"mymall/services/user-service/internal/biz"
	"mymall/services/user-service/internal/svc"
	"mymall/services/user-service/internal/types"
)

type ListRegionsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListRegionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRegionsLogic {
	return &ListRegionsLogic{Logger: logx.WithContext(ctx), svcCtx: svcCtx}
}

func (l *ListRegionsLogic) ListRegions(ctx context.Context, req *types.RegionListReq) (*types.PageListResp, error) {
	list, err := biz.NewRegionLogic(l.svcCtx).ListChildren(ctx, req.ParentCode)
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return &types.PageListResp{List: list}, nil
}
